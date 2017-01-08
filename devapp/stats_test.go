// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package devapp

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReleaseRender(t *testing.T) {
	req, _ := http.NewRequest("GET", "/stats/release", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	// for now, check that this doesn't panic, soon test that it renders
	// correctly
}
