// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"./bentry"
	"./btemplate"
	"./env"
	"log"
	"net/http"
	"regexp"
)

const ENV_GOBLOG_ROOT = "GOBLOG_ROOT"

const TMPL_INDEX = "index"
const TMPL_ENTRY = "entry"
const TMPL_DIR = "template/"
const ENTRY_DIR = "entries/"
const ENTRY_EXT = ".md"
const STATIC_DIR = "static/"

const PATH_STATIC = "/static/"
const PATH_ENTRY = "/entry/"
const PATH_INDEX = "/"

var BLOG_ROOT = env.GetDef(ENV_GOBLOG_ROOT, "../sample")

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
	entries = bentry.BlogContent{}
	entries.Load(BLOG_ROOT + "/" + ENTRY_DIR)

	templates.Load(BLOG_ROOT+"/"+TMPL_DIR, entries.GetEntries)
	http.HandleFunc(PATH_INDEX, makeIndexHandler(TMPL_INDEX))
	http.HandleFunc(PATH_ENTRY, makePageHandler(PATH_ENTRY, TMPL_ENTRY, viewHandler))
	http.Handle(PATH_STATIC, http.StripPrefix(PATH_STATIC,
		http.FileServer(http.Dir(BLOG_ROOT+"/"+STATIC_DIR))))

	http.ListenAndServe(":8080", nil)
}
