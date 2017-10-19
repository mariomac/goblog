// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

const INDEX_TMPL = "index"
const TMPL_DIR = "templates/"
const TMPL_EXT = ".html"
const PAGE_DIR = "pages/"
const PAGE_EXT = ".md"

type Page struct {
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := PAGE_DIR + title + PAGE_EXT
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Todo: redirect or template 404
		return
	}
	renderTemplate(w, INDEX_TMPL, p)
}

var templates = template.Must(template.ParseFiles(TMPL_DIR + INDEX_TMPL + TMPL_EXT))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + TMPL_EXT, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/([_a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[0])
	}
}

func main() {
	http.HandleFunc("/", makeHandler(viewHandler))

	http.ListenAndServe(":8080", nil)
}