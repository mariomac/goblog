// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"github.com/shurcooL/github_flavored_markdown"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	"./env"
	"./filesearch"
)

const ENV_GOBLOG_ROOT = "GOBLOG_ROOT"

const TMPL_INDEX = "index"
const TMPL_ENTRY = "entry"

const TMPL_DIR = "template/"
const TMPL_EXT = ".html"
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
	renderTemplate(w, template, p)
}

var templates = template.Must(
	template.New(TMPL_INDEX).Funcs(template.FuncMap{"md2html": md2html}).ParseFiles(
		filesearch.GetTemplates(BLOG_ROOT + "/" + TMPL_DIR)...))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl + TMPL_EXT, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPage = regexp.MustCompile("^([_\\-a-zA-Z0-9]+)$")

func makeIndexHandler(rootPath string, template string, fn func(http.ResponseWriter, *http.Request, string, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//m := validPage.FindStringSubmatch(r.URL.Path[len(rootPath):])
		//if m == nil {
		//	http.NotFound(w, r)
		//	return
		//}
		renderTemplate(w, template, nil)
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

func md2html(mdText []byte) template.HTML {
	return template.HTML(github_flavored_markdown.Markdown(mdText))
}

func main() {
	// todo: mezclar handler
	indexHandler := makeIndexHandler(PATH_INDEX, TMPL_INDEX, viewHandler)
	filesearch.ScanTimestampedEntries(BLOG_ROOT + "/" + ENTRY_DIR)
	http.HandleFunc(PATH_INDEX, indexHandler)
	http.HandleFunc(PATH_ENTRY, makePageHandler(PATH_ENTRY, TMPL_ENTRY, viewHandler))
	http.Handle(PATH_STATIC, http.StripPrefix(PATH_STATIC,
		http.FileServer(http.Dir(BLOG_ROOT + "/" + STATIC_DIR))))

	http.ListenAndServe(":8080", nil)
}