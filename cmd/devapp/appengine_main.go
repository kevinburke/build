// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Build for running on app engine flex, following the instructions set here:
// https://godoc.org/google.golang.org/appengine#Main

// +build appenginevm

package main

import (
	// this registers HTTP handlers
	_ "golang.org/x/build/devapp"
	"google.golang.org/appengine"
)

func main() {
	appengine.Main()
}
