// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/env"
	"github.com/mariomac/goblog/src/feed"
	"github.com/mariomac/goblog/src/visual"
)

// Env var names
const envGoblogRoot = "GOBLOG_ROOT"
const envGoblogPort = "GOBLOG_PORT"
const envGoblogDomain = "GOBLOG_DOMAIN"

// Template names
const templateIndex = "index"
const templateEntry = "entry"

// Directory names
const dirTemplate = "template/"
const dirEntry = "entries/"
const dirStatic = "static/"

// Path names
const pathStatic = "/static/"
const pathEntry = "/entry/"
const pathIndex = "/"
const pathAtom = "/atom.xml"

var entries blog.Content
var templates = visual.Templates{}

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
	var blogDomain = env.GetDef(envGoblogDomain, osHostname)
	var blogRoot = env.GetDef(envGoblogRoot, "./sample")
	var blogPort = env.GetDef(envGoblogPort, "8080")

	log.Printf("Environment: { %s=\"%s\", %s=\"%s\", %s=\"%s\",",
		envGoblogDomain, blogDomain,
		envGoblogPort, blogPort,
		envGoblogRoot, blogRoot)

	// Load blog entries
	entries = blog.Content{}
	entries.Load(blogRoot + "/" + dirEntry)

	// Create Atom XML feed
	atomxml := bytes.NewBufferString(
		feed.BuildAtomFeed(entries.GetEntries(), blogDomain, pathEntry)).Bytes()

	// Load templates
	templates.Load(blogRoot+"/"+dirTemplate, entries.GetEntries)

	http.HandleFunc(pathIndex, makeIndexHandler(templateIndex))
	http.HandleFunc(pathAtom, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/atom+xml")
		writer.Write(atomxml)
	})
	http.HandleFunc(pathEntry, makePageHandler(pathEntry, templateEntry, viewHandler))
	http.Handle(pathStatic, http.StripPrefix(pathStatic,
		http.FileServer(http.Dir(blogRoot+"/"+dirStatic))))

	log.Printf("GoBlog is listening at port %s", blogPort)
	http.ListenAndServe(":"+blogPort, nil)
}
