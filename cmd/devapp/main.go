// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Devapp generates the dashboard that powers dev.golang.org.
//
// Usage:
//
//	devapp --port=8081
//
// By default devapp listens on port 8081.
//
// Github issues and Gerrit CL's are stored in memory in the running process.
// To trigger an initial download, visit http://localhost:8081/update or
// http://localhost:8081/update/stats in your browser.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/build/devapp"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func init() {
	flag.Usage = func() {
		os.Stderr.WriteString(`usage: devapp [-port=port]

Devapp generates the dashboard that powers dev.golang.org.

`)
		flag.PrintDefaults()
	}
}

func main() {
	templateDir := flag.String("templates", "", "load templates/JS/CSS from disk in this directory (for local development)")
	port := flag.Uint("port", 8081, "Port to listen on")
	flag.Parse()
	var fs vfs.FileSystem
	if *templateDir != "" {
		fs = vfs.OS(*templateDir)
	} else {
		fs = mapfs.New(devapp.Files)
	}
	mux := devapp.NewServeMux(fs)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "Listening on port %d\n", *port)
	log.Fatal(http.Serve(ln, mux))
}
