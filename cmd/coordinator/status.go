// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"time"
)

func handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	round := func(t time.Duration) time.Duration {
		return t / time.Second * time.Second
	}
	df := diskFree()

	statusMu.Lock()
	data := statusData{
		Total:        len(status),
		Uptime:       round(time.Now().Sub(processStartTime)),
		Recent:       append([]*buildStatus{}, statusDone...),
		DiskFree:     df,
		Version:      Version,
		NumFD:        fdCount(),
		NumGoroutine: runtime.NumGoroutine(),
	}
	for _, st := range status {
		if atomic.LoadInt32(&st.hasBuildlet) != 0 {
			data.ActiveBuilds++
			data.Active = append(data.Active, st)
		} else {
			data.Pending = append(data.Pending, st)
		}
	}
	// TODO: make this prettier.
	var buf bytes.Buffer
	for _, key := range tryList {
		if ts := tries[key]; ts != nil {
			state := ts.state()
			fmt.Fprintf(&buf, "Change-ID: %v Commit: %v (<a href='/try?commit=%v'>status</a>)\n",
				key.ChangeTriple(), key.Commit, key.Commit[:8])
			fmt.Fprintf(&buf, "   Remain: %d, fails: %v\n", state.remain, state.failed)
			for _, bs := range ts.builds {
				fmt.Fprintf(&buf, "  %s: running=%v\n", bs.name, bs.isRunning())
			}
		}
	}
	statusMu.Unlock()

	data.RemoteBuildlets = template.HTML(remoteBuildletStatus())

	sort.Sort(byAge(data.Active))
	sort.Sort(byAge(data.Pending))
	sort.Sort(sort.Reverse(byAge(data.Recent)))
	if errTryDeps != nil {
		data.TrybotsErr = errTryDeps.Error()
	} else {
		if buf.Len() == 0 {
			data.Trybots = template.HTML("<i>(none)</i>")
		} else {
			data.Trybots = template.HTML("<pre>" + buf.String() + "</pre>")
		}
	}

	buf.Reset()
	gcePool.WriteHTMLStatus(&buf)
	data.GCEPoolStatus = template.HTML(buf.String())
	buf.Reset()

	kubePool.WriteHTMLStatus(&buf)
	data.KubePoolStatus = template.HTML(buf.String())
	buf.Reset()

	reversePool.WriteHTMLStatus(&buf)
	data.ReversePoolStatus = template.HTML(buf.String())

	buf.Reset()
	if err := statusTmpl.Execute(&buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func fdCount() int {
	f, err := os.Open("/proc/self/fd")
	if err != nil {
		return -1
	}
	defer f.Close()
	n := 0
	for {
		names, err := f.Readdirnames(1000)
		n += len(names)
		if err == io.EOF {
			return n
		}
		if err != nil {
			return -1
		}
	}
}

func diskFree() string {
	out, _ := exec.Command("df", "-h").Output()
	return string(out)
}

// statusData is the data that fills out statusTmpl.
type statusData struct {
	Total             int // number of total builds (including those waiting for a buildlet)
	ActiveBuilds      int // number of running builds (subset of Total with a buildlet)
	NumFD             int
	NumGoroutine      int
	Uptime            time.Duration
	Active            []*buildStatus // have a buildlet
	Pending           []*buildStatus // waiting on a buildlet
	Recent            []*buildStatus
	TrybotsErr        string
	Trybots           template.HTML
	GCEPoolStatus     template.HTML // TODO: embed template
	KubePoolStatus    template.HTML // TODO: embed template
	ReversePoolStatus template.HTML // TODO: embed template
	RemoteBuildlets   template.HTML
	DiskFree          string
	Version           string
}

var statusTmpl = template.Must(template.New("status").Parse(`
<!DOCTYPE html>
<html>
<head><link rel="stylesheet" href="/style.css"/><title>Go Farmer</title></head>
<body>
<header id="topbar">
	<div class="header-container">
		<div class="top-heading">
			<h1>Go Build Coordinator</h1>
		</div>
		<ul id="menu" class="clearfix">
			<li class="button">
				<a href="https://build.golang.org">Dashboard</a>
			</li>
			<li class="button">
				<a href="/builders">Builders</a>
			</li>
		</nav>
	</div>
</header>

<div class="container">
<h2>Running</h2>
<p>{{printf "%d" .Total}} total builds; {{printf "%d" .ActiveBuilds}} active. Uptime {{printf "%s" .Uptime}}. Version {{.Version}}.

<h2 id=trybots>Active Trybot Runs <a href='#trybots'>¶</a></h2>
{{- if .TrybotsErr}}
<b>trybots disabled:</b>: {{.TrybotsErr}}
{{else}}
{{.Trybots}}
{{end}}

<h2 id=remote>Remote buildlets <a href='#remote'>¶</a></h3>
{{.RemoteBuildlets}}

<h2 id=pools>Buildlet pools <a href='#pools'>¶</a></h2>
<ul>
<li>{{.GCEPoolStatus}}</li>
<li>{{.KubePoolStatus}}</li>
<li>{{.ReversePoolStatus}}</li>
</ul>

<h2 id=active>Active builds <a href='#active'>¶</a></h2>
<ul>
{{range .Active}}
<li><pre>{{.HTMLStatusLine}}</pre></li>
{{end}}
</ul>

<h2 id=pending>Pending builds <a href='#pending'>¶</a></h2>
<ul>
{{range .Pending}}
<li><pre>{{.HTMLStatusLine}}</pre></li>
{{end}}
</ul>

<h2 id=completed>Recently completed <a href='#completed'>¶</a></h2>
<ul>
{{range .Recent}}
<li><span>{{.HTMLStatusLine_done}}</span></li>
{{end}}
</ul>

<h2 id=disk>Disk Space <a href='#disk'>¶</a></h2>
<pre>{{.DiskFree}}</pre>

<h2 id=disk>File Descriptors <a href='#fd'>¶</a></h2>
<p>{{.NumFD}}</p>

<h2 id=disk>Goroutines <a href='#goroutines'>¶</a></h2>
<p>{{.NumGoroutine}} <a href='/debug/goroutines'>goroutines</a></p>

</div>
</body>
</html>
`))

func handleStyleCSS(w http.ResponseWriter, r *http.Request) {
	src := strings.NewReader(styleCSS)
	http.ServeContent(w, r, "style.css", processStartTime, src)
}

const styleCSS = `
body {
	font-family: sans-serif;
	color: #222;
	margin: 0;
}

pre {
	font-family: Menlo, Consolas, monospace;
	font-size: 9pt;
}

h2 {
	font-size: 20px;
	color: #480048;
	background-color: #ddd3ee;
	line-height: 1.25;
	font-weight: normal;
	padding: 8px;
	margin: 20px 0 20px;
}

h2 a {
	text-decoration: none;
	color: #480048;
}

h2 a:hover {
	cursor: pointer;
	text-decoration: underline;
}

.header-container {
	max-width: 950px;
	margin: 0 auto;
}

#topbar {
	background-color: #ddd3ee;
	height: 64px;
}
.top-heading {
	float: left;
}
#topbar h1 {
	font-size: 20px;
	font-weight: normal;
	margin: 0;
	color: #222;
	line-height: 1.25;
	padding: 21px 0;
}
#menu {
	float: left;
	padding: 10px;
	margin-left: 100px;
}

.button {
	float: left;
	list-style-type: none;
}

#menu a {
	padding: 10px;
	margin: 0;
	margin-right: 5px;
	color: white;
	background: #480048;
	text-decoration: none;
	font-size: 16px;
	border: 1px solid #480048;
	border-radius: 5px;
}

.container {
	max-width: 950px;
	margin: 0 auto;
}

table {
	border-collapse: collapse;
	font-size: 9pt;
}

table td, table th, table td, table th {
	text-align: left;
	vertical-align: top;
	padding: 2px 6px;
}

table thead tr {
	background: #fff !important;
}
`
