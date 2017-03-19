// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Logic to interact with a Gerrit server. Gerrit has an entire Git-based
// protocol for fetching metadata about CL's, reviewers, patch comments, which
// is used here - we don't use the x/build/gerrit client, which hits the API.
// TODO: write about Gerrit's Git API.

package maintner

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/build/maintner/maintpb"
)

// Gerrit holds information about a number of Gerrit projects.
type Gerrit struct {
	c       *Corpus
	dataDir string // the root Corpus data directory
	// keys are like "https://go-review.googlesource.com/build"
	projects map[string]*GerritProject
}

// c.mu must be held
func (g *Gerrit) getOrCreateProject(gerritURL string) *GerritProject {
	proj, ok := g.projects[gerritURL]
	if ok {
		return proj
	}
	proj = &GerritProject{
		gerrit: g,
		url:    gerritURL,
		gitDir: filepath.Join(g.dataDir, url.PathEscape(gerritURL)),
		cls:    map[int32]*gerritCL{},
	}
	g.projects[gerritURL] = proj
	return proj
}

// GerritProject represents a single Gerrit project.
type GerritProject struct {
	gerrit *Gerrit
	url    string // "https://go-review.googlesource.com"
	// TODO: Many different Git remotes can share the same Gerrit instance, e.g.
	// the Go Gerrit instance supports build, gddo, go. For the moment these are
	// all treated separately, since the remotes are separate.
	gitDir string
	cls    map[int32]*gerritCL
}

type gerritCL struct {
	Hash       gitHash
	Number     int32
	Author     *gitPerson
	AuthorTime time.Time
	Status     string // "merged", "abandoned", "" ("open" is implicit)
	// TODO...
}

// c.mu must be held
func (c *Corpus) initGerrit() {
	if c.gerrit != nil {
		return
	}
	c.gerrit = &Gerrit{
		c:        c,
		dataDir:  c.dataDir,
		projects: map[string]*GerritProject{},
	}
}

type watchedGerritRepo struct {
	project *GerritProject
}

// AddGerrit adds the Gerrit project with the given URL to the corpus.
func (c *Corpus) AddGerrit(gerritURL string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.initGerrit()
	project := c.gerrit.getOrCreateProject(gerritURL)
	if project == nil {
		panic("gerrit project not created")
	}
	c.watchedGerritRepos = append(c.watchedGerritRepos, watchedGerritRepo{
		project: project,
	})
}

func (c *Corpus) processGerritMutation(gm *maintpb.GerritMutation) {
	// TODO
}

// map of CL numbers to hashes, see the regex below
type remoteResponse map[int32]gitHash

// sample row:
// fd1e71f1594ce64941a85428ddef2fbb0ad1023e	refs/changes/99/30599/3
//
// The "99" in the middle covers all CL's that end in "99", so
// refs/changes/99/99/1, refs/changes/99/199/meta.
//
// The last value is an integer representing a patch set (1, 2, 3), or "meta" (a
// special commit holding the comments for a commit)
var remoteRegex = regexp.MustCompile(`^(.+)\s+refs/changes/[0-9a-f]{2}/([0-9]+)/(.+)$`)

func (gp *GerritProject) sync(ctx context.Context) error {
	c := gp.gerrit.c
	if err := gp.init(ctx); err != nil {
		return err
	}
	// TODO abstract the cmd running boilerplate
	fetchCtx, cancel := context.WithTimeout(ctx, time.Minute)
	cmd := exec.CommandContext(fetchCtx, "git", "fetch", "origin")
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Dir = gp.gitDir
	err := cmd.Run()
	cancel()
	if err != nil {
		os.Stderr.Write(buf.Bytes())
		return err
	}
	cmd = exec.CommandContext(ctx, "git", "ls-remote")
	buf.Reset()
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Dir = gp.gitDir
	if err := cmd.Run(); err != nil {
		os.Stderr.Write(buf.Bytes())
		return err
	}
	rr := make(remoteResponse)
	bs := bufio.NewScanner(buf)
	for bs.Scan() {
		m := remoteRegex.FindStringSubmatch(bs.Text())
		if m == nil {
			continue
		}
		if m[3] == "meta" {
			i, err := strconv.ParseInt(m[2], 10, 32)
			if err != nil {
				return fmt.Errorf("maintner: error parsing CL as a number: %v", err)
			}
			rr[int32(i)] = gitHashFromHexStr(m[1])
		}
	}
	if err := bs.Err(); err != nil {
		return err
	}
	var toFetch []gitHash
	for cl, hash := range rr {
		c.mu.RLock()
		val, ok := gp.cls[cl]
		c.mu.RUnlock()
		if !ok || !val.Hash.Equal(hash) {
			toFetch = append(toFetch, hash)
		}
	}
	if err := gp.fetchHashes(ctx, toFetch); err != nil {
		return err
	}
	for cl, hash := range rr {
		c.mu.RLock()
		val, ok := gp.cls[cl]
		c.mu.RUnlock()
		if !ok || !val.Hash.Equal(hash) {
			// TODO: parallelize updates if this gets slow, we can probably do
			// lots of filesystem reads without penalty
			if err := gp.updateCL(ctx, cl, hash); err != nil {
				return err
			}
		}
	}
	return nil
}

var (
	statusSpace = []byte("Status: ")
)

// mustTimestamp turns a time.Time into a *timestamp.Timestamp or panics if in
// is invalid.
func mustTimestamp(in time.Time) *timestamp.Timestamp {
	tp, err := ptypes.TimestampProto(in)
	if err != nil {
		panic(err)
	}
	return tp
}

// newMutationFromCL generates a GerritCLMutation using the smallest possible
// diff between a (the state we have in memory) and b (the current Gerrit
// state).
//
// If newMutationFromCL returns nil, the provided gerrit CL is no newer than
// the data we have in the corpus. 'a' may be nil.
func (gp *GerritProject) newMutationFromCL(a, b *gerritCL) *maintpb.Mutation {
	if b == nil {
		panic("newMutationFromCL: provided nil gerritCL")
	}
	if a == nil {
		var sha1 string
		switch b.Hash.(type) {
		case gitSHA1:
			sha1 = b.Hash.String()
		default:
			panic(fmt.Sprintf("unsupported git hash type %T", b.Hash))
		}
		return &maintpb.Mutation{
			Gerrit: &maintpb.GerritMutation{
				Url:        gp.url,
				Sha1:       sha1,
				Number:     b.Number,
				Status:     b.Status,
				Author:     b.Author.str,
				AuthorTime: mustTimestamp(b.AuthorTime),
			},
		}
	}
	// TODO: update the existing proto
	return nil
}

// updateCL updates the local CL.
func (gp *GerritProject) updateCL(ctx context.Context, clNum int32, hash gitHash) error {
	cmd := exec.CommandContext(ctx, "git", "cat-file", "-p", hash.String())
	cmd.Dir = gp.gitDir
	buf, errBuf := new(bytes.Buffer), new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = errBuf
	if err := cmd.Run(); err != nil {
		return err
	}
	catFile := buf.Bytes()
	cl := &gerritCL{
		Number: clNum,
		Hash:   hash,
	}
	c := gp.gerrit.c
	c.mu.Lock()
	err := foreachLine(catFile, func(ln []byte) error {
		if bytes.HasPrefix(ln, authorSpace) {
			p, t, err := c.parsePerson(ln[len(authorSpace):])
			if err != nil {
				return fmt.Errorf("unrecognized author line %q: %v", ln, err)
			}
			cl.Author = p
			cl.AuthorTime = t
		}
		if bytes.HasPrefix(ln, statusSpace) {
			cl.Status = string(ln[len(statusSpace):])
		}
		return nil
	})
	if err != nil {
		c.mu.Unlock()
		return err
	}
	proto := gp.newMutationFromCL(gp.cls[clNum], cl)
	gp.cls[clNum] = cl
	c.mu.Unlock()
	c.processMutation(proto)
	return nil
}

func (gp *GerritProject) fetchHashes(ctx context.Context, hashes []gitHash) error {
	for i := 0; i < len(hashes); i = i + 500 {
		var slice []gitHash
		if i+500 < len(hashes) {
			slice = hashes[i : i+500]
		} else {
			slice = hashes[i:len(hashes)]
		}
		args := []string{"fetch", "--quiet", "origin"}
		for _, hash := range hashes {
			args = append(args, hash.String())
		}
		cmd := exec.CommandContext(ctx, "git", args...)
		buf := new(bytes.Buffer)
		cmd.Dir = gp.gitDir
		cmd.Stdout = buf
		cmd.Stderr = buf
		if err := cmd.Run(); err != nil {
			log.Println("error fetching", len(slice), "hashes from git remote", gp.url)
			os.Stderr.Write(buf.Bytes())
			return err
		}
	}
	return nil
}

func (gp *GerritProject) init(ctx context.Context) error {
	if err := os.MkdirAll(gp.gitDir, 0755); err != nil {
		return err
	}
	// try to short circuit a git init error, since the init error matching is
	// brittle
	if _, err := exec.LookPath("git"); err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(gp.gitDir, ".git", "config")); err == nil {
		remoteBytes, err := exec.CommandContext(ctx, "git", "remote", "-v").Output()
		if err != nil {
			return err
		}
		if !strings.Contains(string(remoteBytes), "origin") && !strings.Contains(string(remoteBytes), gp.url) {
			return fmt.Errorf("didn't find origin & gp.url in remote output %s", string(remoteBytes))
		}
	} else {
		cmd := exec.CommandContext(ctx, "git", "init")
		buf := new(bytes.Buffer)
		cmd.Stdout = buf
		cmd.Stderr = buf
		cmd.Dir = gp.gitDir
		if err := cmd.Run(); err != nil {
			log.Printf(`Error running "git init": %s`, buf.String())
			return err
		}
		buf.Reset()
		cmd = exec.CommandContext(ctx, "git", "remote", "add", "origin", gp.url)
		cmd.Stdout = buf
		cmd.Stderr = buf
		cmd.Dir = gp.gitDir
		if err := cmd.Run(); err != nil {
			log.Printf(`Error running "git remote add origin": %s`, buf.String())
			return err
		}
	}
	return nil
}
