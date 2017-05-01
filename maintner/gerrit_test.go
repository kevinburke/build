// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maintner

import (
	"testing"
	"time"
)

var statusTests = []struct {
	msg  string
	want string
}{
	{`    Create change

    Uploaded patch set 1.

    Patch-set: 1 (draft)
    Change-id: I38a08cacc17bcd9587475495111fe98f10d6875c
    Subject: test: test
    Branch: refs/heads/master
    Status: draft
    Topic:
    Commit: fee468c613a70d89f60fb5d683b0f796aabecaac`, "draft"},
	{`   Update patch set 1

    Change has been successfully cherry-picked as 117ac82c422a11e4dd5f4c14b50bafc1df840481 by Brad Fitzpatrick

    Patch-set: 1
    Status: merged
    Submission-id: 16401-1446004855021-a20b3823`, "merged"},
	{`    Create patch set 8

    Uploaded patch set 8: Patch Set 7 was rebased.

    Patch-set: 8
    Subject: devapp: initial support for App Engine Flex
    Commit: 17839a9f284b473986f235ad2757a2b445d05068
    Tag: autogenerated:gerrit:newPatchSet
    Groups: 17839a9f284b473986f235ad2757a2b445d05068`, ""},
}

func TestGetGerritStatus(t *testing.T) {
	for _, tt := range statusTests {
		gc := &GitCommit{Msg: tt.msg}
		got := getGerritStatus(gc)
		if got != tt.want {
			t.Errorf("getGerritStatus msg:\n%s\ngot %s, want %s", tt.msg, got, tt.want)
		}
	}
}

var messageTests = []struct {
	in  string
	out string
}{
	{`Update patch set 1

Patch Set 1: Code-Review+2

Just to confirm, "go test" will consider an empty test file to be passing?

Patch-set: 1
Reviewer: Quentin Smith <13020@62eb7196-b449-3ce5-99f1-c037f21e1705>
Label: Code-Review=+2
`, "Patch Set 1: Code-Review+2\n\nJust to confirm, \"go test\" will consider an empty test file to be passing?"},
}

func TestGetGerritMessage(t *testing.T) {
	var c Corpus
	c.EnableLeaderMode(new(dummyMutationLogger), "/fake/dir")
	c.TrackGerrit("go.googlesource.com/build")
	gp := c.gerrit.projects["go.googlesource.com/build"]
	for _, tt := range messageTests {
		gc := &GitCommit{
			Msg:        tt.in,
			CommitTime: time.Now().UTC(),
		}
		msg := gp.getGerritMessage(gc)
		if msg.Date.IsZero() {
			t.Errorf("getGerritMessage: expected Date to be non-zero, got zero")
		}
		if msg.Message != tt.out {
			t.Errorf("getGerritMessage: want %q, got %q", tt.out, msg.Message)
		}
	}
}
