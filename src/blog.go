// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"./bentry"
	"./btemplate"
	"./env"
	"./feed"
	"log"
	"net/http"
	"regexp"
	"os"
	"bytes"
)

const ENV_GOBLOG_ROOT = "GOBLOG_ROOT"
const ENV_GOBLOG_PORT = "GOBLOG_PORT"
const ENV_GOBLOG_DOMAIN = "GOBLOG_DOMAIN"

const TMPL_INDEX = "index"
const TMPL_ENTRY = "entry"
const TMPL_DIR = "template/"
const ENTRY_DIR = "entries/"
const STATIC_DIR = "static/"

const PATH_STATIC = "/static/"
const PATH_ENTRY = "/entry/"
const PATH_INDEX = "/"
const PATH_ATOM = "/atom.xml"

var entries bentry.BlogContent
var templates = btemplate.Templates{}

func viewHandler(w http.ResponseWriter, r *http.Request, fileName string, template string) {
	log.Printf("viewHandler(_, _, %s, %s)", fileName, template)
	entry, found := entries.Get(fileName)
	if !found {
		http.Error(w, "Entry not found "+fileName, http.StatusNotFound) // Todo: redirect or template 404
		return
	}
	templates.Render(w, template, entry)
}

func makeIndexHandler(template string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates.Render(w, template, nil)
	}
}

var validPagePath = regexp.MustCompile("^([_\\-a-zA-Z0-9]+)\\.md$")

func makePageHandler(rootPath string, template string,
	fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPagePath.FindStringSubmatch(r.URL.Path[len(rootPath):])
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[0], template)
	}
}

func main() {
	log.Print("Starting GoBlog...")

	osHostname, _ := os.Hostname()
	var BLOG_DOMAIN = env.GetDef(ENV_GOBLOG_DOMAIN, osHostname)
	var BLOG_ROOT = env.GetDef(ENV_GOBLOG_ROOT, "../sample")
	var BLOG_PORT = env.GetDef(ENV_GOBLOG_PORT, "8080")

	log.Printf("Environment: { %s=\"%s\", %s=\"%s\", %s=\"%s\",",
		ENV_GOBLOG_DOMAIN, BLOG_DOMAIN,
		ENV_GOBLOG_PORT, BLOG_PORT,
		ENV_GOBLOG_ROOT, BLOG_ROOT)

	// Load blog entries
	entries = bentry.BlogContent{}
	entries.Load(BLOG_ROOT + "/" + ENTRY_DIR)

	// Create Atom XML feed
	atomxml := bytes.NewBufferString(
		feed.BuildAtomFeed(entries.GetEntries(), BLOG_DOMAIN, PATH_ENTRY)).Bytes()

	// Load templates
	templates.Load(BLOG_ROOT+"/"+TMPL_DIR, entries.GetEntries)

	http.HandleFunc(PATH_INDEX, makeIndexHandler(TMPL_INDEX))
	http.HandleFunc(PATH_ATOM, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/atom+xml")
		writer.Write(atomxml)
	})
	http.HandleFunc(PATH_ENTRY, makePageHandler(PATH_ENTRY, TMPL_ENTRY, viewHandler))
	http.Handle(PATH_STATIC, http.StripPrefix(PATH_STATIC,
		http.FileServer(http.Dir(BLOG_ROOT+"/"+STATIC_DIR))))

	log.Printf("GoBlog is listening at port %s", BLOG_PORT)
	http.ListenAndServe(":" + BLOG_PORT, nil)
}
