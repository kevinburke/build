// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file gets built on App Engine Flex.

// +build appenginevm

package devapp

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/user"
)

var lg *logging.Logger

func init() {
	onAppengine = !appengine.IsDevAppServer()
	log = &appengineLogger{}
	client, err := logging.NewClient(context.Background(), "devapp")
	if err == nil {
		lg = client.Logger("log")
	}

	http.HandleFunc("/setToken", setTokenHandler)
}

type appengineLogger struct{}

func (a *appengineLogger) Infof(_ context.Context, format string, args ...interface{}) {
	if lg == nil {
		return
	}
	lg.Log(logging.Entry{
		Severity: logging.Info,
		Payload:  fmt.Sprintf(format, args...),
	})
}

func (a *appengineLogger) Errorf(_ context.Context, format string, args ...interface{}) {
	if lg == nil {
		return
	}
	lg.Log(logging.Entry{
		Severity: logging.Error,
		Payload:  fmt.Sprintf(format, args...),
	})
}

func (a *appengineLogger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	if lg == nil {
		return
	}
	lg.LogSync(ctx, logging.Entry{
		Severity: logging.Critical,
		Payload:  fmt.Sprintf(format, args...),
	})
}

func newTransport(ctx context.Context) http.RoundTripper {
	return &urlfetch.Transport{Context: ctx}
}

func currentUserEmail(ctx context.Context) string {
	u := user.Current(ctx)
	if u == nil {
		return ""
	}
	return u.Email
}

// loginURL returns a URL that, when visited, prompts the user to sign in,
// then redirects the user to the URL specified by dest.
func loginURL(ctx context.Context, path string) (string, error) {
	return user.LoginURL(ctx, path)
}

func logoutURL(ctx context.Context, path string) (string, error) {
	return user.LogoutURL(ctx, path)
}

func getCache(ctx context.Context, name string) (*Cache, error) {
	var cache Cache
	if err := datastore.Get(ctx, datastore.NewKey(ctx, entityPrefix+"Cache", name, 0, nil), &cache); err != nil {
		return &cache, err
	}
	return &cache, nil
}

func getCaches(ctx context.Context, names ...string) map[string]*Cache {
	out := make(map[string]*Cache)
	var keys []*datastore.Key
	var ptrs []*Cache
	for _, name := range names {
		keys = append(keys, datastore.NewKey(ctx, entityPrefix+"Cache", name, 0, nil))
		out[name] = &Cache{}
		ptrs = append(ptrs, out[name])
	}
	datastore.GetMulti(ctx, keys, ptrs) // Ignore errors since they might not exist.
	return out
}

func getPage(ctx context.Context, page string) (*Page, error) {
	entity := new(Page)
	err := datastore.Get(ctx, datastore.NewKey(ctx, entityPrefix+"Page", page, 0, nil), entity)
	return entity, err
}

func writePage(ctx context.Context, page string, content []byte) error {
	entity := &Page{
		Content: content,
	}
	_, err := datastore.Put(ctx, datastore.NewKey(ctx, entityPrefix+"Page", page, 0, nil), entity)
	return err
}

func putCache(ctx context.Context, name string, c *Cache) error {
	_, err := datastore.Put(ctx, datastore.NewKey(ctx, entityPrefix+"Cache", name, 0, nil), c)
	return err
}

func getToken(ctx context.Context) (string, error) {
	cache, err := getCache(ctx, "github-token")
	if err != nil {
		return "", err
	}
	return string(cache.Value), nil
}

// Store a token in the database
func setTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", 405)
		return
	}
	ctx := r.Context()
	r.ParseForm()
	if value := r.Form.Get("value"); value != "" {
		var token Cache
		token.Value = []byte(value)
		if err := putCache(ctx, "github-token", &token); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func getContext(r *http.Request) context.Context {
	return r.Context()
}
