// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"./env"
	"./filesearch"
	"./btemplate"
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

type Page struct {
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := BLOG_ROOT + "/" + ENTRY_DIR + title + ENTRY_EXT
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string, template string) {
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Todo: redirect or template 404
		return
	}
	templates.Render(w, template, p)
}

// TODO: retrigger when directory changes
var entries = filesearch.ScanTimestampedEntries(BLOG_ROOT + "/" + ENTRY_DIR)
func getEntries() []*filesearch.Entry {
	return entries
}

var templates = btemplate.Templates{}

var validPage = regexp.MustCompile("^([_\\-a-zA-Z0-9]+)$")

func makeIndexHandler(rootPath string, template string, fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//m := validPage.FindStringSubmatch(r.URL.Path[len(rootPath):])
		//if m == nil {
		//	http.NotFound(w, r)
		//	return
		//}
		templates.Render(w, template, nil)
	}
}

func makePageHandler(rootPath string, template string, fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPage.FindStringSubmatch(r.URL.Path[len(rootPath):])
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[0], template)
	}
}


func main() {
	templates.Load(BLOG_ROOT+"/"+TMPL_DIR, getEntries)
	http.HandleFunc(PATH_INDEX, makeIndexHandler(PATH_INDEX, TMPL_INDEX, viewHandler))
	http.HandleFunc(PATH_ENTRY, makePageHandler(PATH_ENTRY, TMPL_ENTRY, viewHandler))
	http.Handle(PATH_STATIC, http.StripPrefix(PATH_STATIC,
		http.FileServer(http.Dir(BLOG_ROOT + "/" + STATIC_DIR))))

	http.ListenAndServe(":8080", nil)
}