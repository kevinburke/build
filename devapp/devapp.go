// Copyright 2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package devapp implements a simple App Engine app for generating
// and serving Go project release dashboards using the godash
// command/library.
package devapp

// Run this to copy all the static files on disk into the binary

//go:generate go run makestatic.go

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/build/gerrit"
	"golang.org/x/build/godash"
	"golang.org/x/net/context"
	"golang.org/x/tools/godoc/vfs"
)

const entityPrefix = "DevApp"

var gerritTransport http.RoundTripper
var releaseTpl *template.Template
var dashBytes []byte

func NewServeMux(fs vfs.FileSystem) *http.ServeMux {
	dashBytes = mustOpen(fs, "templates/dash.html")
	tplBytes := mustOpen(fs, "templates/release.html")
	releaseTpl = template.Must(template.New("release").Parse(string(tplBytes)))
	mux := http.NewServeMux()
	for _, page := range []string{"release", "cl"} {
		page := page
		mux.Handle("/"+page, hstsHandler(func(w http.ResponseWriter, r *http.Request) { servePage(w, r, page) }))
	}
	mux.Handle("/dash", hstsHandler(func(w http.ResponseWriter, r *http.Request) {
		showDash(w, r, fs)
	}))
	mux.Handle("/update", ctxHandler(update))
	// Defined in stats.go
	mux.HandleFunc("/stats/raw", rawHandler)
	mux.HandleFunc("/stats/svg", svgHandler)
	mux.Handle("/stats/release", ctxHandler(release))
	mux.Handle("/stats/release/data.js", ctxHandler(releaseData))
	mux.Handle("/update/stats", ctxHandler(updateStats))
	registerPerBuildHandlers(mux, fs)
	return mux
}

// mustOpen returns the contents of the file with the given name from fs, or
// panics.
func mustOpen(fs vfs.FileSystem, filename string) []byte {
	f, err := fs.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	bits, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return bits
}

// hstsHandler wraps an http.HandlerFunc such that it sets the HSTS header.
func hstsHandler(fn http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; preload")
		fn(w, r)
	})
}

func ctxHandler(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := getContext(r)
		if err := fn(ctx, w, r); err != nil {
			log.Criticalf(ctx, "handler failed: %v", err)
			http.Error(w, err.Error(), 500)
		}
	})
}

func logFn(ctx context.Context, w io.Writer) func(string, ...interface{}) {
	logger := stdlog.New(w, "", stdlog.Lmicroseconds)
	return func(format string, args ...interface{}) {
		logger.Printf(format, args...)
		log.Infof(ctx, format, args...)
	}
}

type Page struct {
	// Content is the complete HTML of the page.
	Content []byte
}

func servePage(w http.ResponseWriter, r *http.Request, page string) {
	ctx := getContext(r)
	entity, err := getPage(ctx, page)
	if err != nil {
		http.Error(w, "page not found", 404)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(entity.Content)
}

type countTransport struct {
	http.RoundTripper
	count int64
}

func (ct *countTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt64(&ct.count, 1)
	return ct.RoundTripper.RoundTrip(req)
}

func (ct *countTransport) Count() int64 {
	return atomic.LoadInt64(&ct.count)
}

func update(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	token, err := getToken(ctx)
	if err != nil {
		return err
	}
	gzdata, _ := getCache(ctx, "gzdata")
	ct := &countTransport{newTransport(ctx), 0}
	gh := godash.NewGitHubClient("golang/go", token, ct)
	defer func() {
		log.Infof(ctx, "Sent %d requests to GitHub", ct.Count())
	}()
	ger := gerrit.NewClient("https://go-review.googlesource.com", gerrit.NoAuth)
	// Without a deadline, urlfetch will use a 5s timeout which is too slow for Gerrit.
	gerctx, cancel := context.WithTimeout(ctx, 9*time.Minute)
	defer cancel()
	ger.HTTPClient = &http.Client{Transport: newTransport(gerctx)}

	data, err := parseData(gzdata)
	if err != nil {
		return err
	}

	if err := data.Reviewers.LoadGithub(ctx, gh); err != nil {
		log.Criticalf(ctx, "failed to load reviewers: %v", err)
		return err
	}
	l := logFn(ctx, w)
	if err := data.FetchData(gerctx, gh, ger, l, 7, false, false); err != nil {
		log.Criticalf(ctx, "failed to fetch data: %v", err)
		return err
	}

	for _, cls := range []bool{false, true} {
		var output bytes.Buffer
		kind := "release"
		if cls {
			kind = "CL"
		}
		fmt.Fprintf(&output, "Go %s dashboard\n", kind)
		fmt.Fprintf(&output, "%v\n\n", time.Now().UTC().Format(time.UnixDate))
		fmt.Fprintf(&output, "HOWTO\n\n")
		if cls {
			data.PrintCLs(&output)
		} else {
			fmt.Fprintf(&output, fmt.Sprintf(`<a href="/stats/release?cycle=%d">Go 1.%d Issue Stats Dashboard</a>`, data.GoReleaseCycle, data.GoReleaseCycle))
			data.PrintIssues(&output)
		}
		var html bytes.Buffer
		godash.PrintHTML(&html, output.String())

		if err := writePage(ctx, strings.ToLower(kind), html.Bytes()); err != nil {
			return err
		}
	}
	return writeCache(ctx, "gzdata", &data)
}
