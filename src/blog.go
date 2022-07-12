// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/mariomac/goblog/src/install"
	"github.com/mariomac/goblog/src/legacy"
	"github.com/mariomac/goblog/src/logr"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/conn"
	"github.com/mariomac/goblog/src/feed"
	"github.com/mariomac/goblog/src/visual"
)

var log = logr.Get()

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
		http.Error(w, "Render not found "+fileName, http.StatusNotFound) // Todo: redirect or template 404
		return
	}
	templates.Render(w, template, entry)
}

func makeIndexHandler(template string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates.Render(w, template, nil)
	}
}

var validPagePath = regexp.MustCompile(`^([_\-a-zA-Z0-9]+)\.md$`)

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
	cfgPath := flag.String("cfg", "", "Path of the YAML configuration file")
	help := flag.Bool("h", false, "This help")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}
	yamlConfig := *cfgPath
	if env, ok := os.LookupEnv("GOBLOG_CONFIG"); ok {
		yamlConfig = env
	}
	cfg, err := install.ReadConfig(yamlConfig)
	if err != nil {
		panic(err)
	}

	log.Printf("Configuration: %#v", cfg)

	log.Print("Starting GoBlog...")

	// Load blog entries
	entries = blog.Content{}
	entries.Load(cfg.RootPath + "/" + dirEntry)

	// Create Atom XML feed
	atomxml := bytes.NewBufferString(
		feed.BuildAtomFeed(entries.GetEntries(), cfg.Domain, pathEntry)).Bytes()

	// Load templates
	templates.Load(cfg.RootPath+"/"+dirTemplate, entries.GetEntries)

	mux := http.NewServeMux()
	mux.HandleFunc(pathIndex, makeIndexHandler(templateIndex))
	mux.HandleFunc(pathAtom, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/xml")
		writer.Write(atomxml)
	})
	mux.HandleFunc(pathEntry, makePageHandler(pathEntry, templateEntry, viewHandler))
	mux.Handle(pathStatic, http.StripPrefix(pathStatic,
		http.FileServer(http.Dir(cfg.RootPath+"/"+dirStatic))))

	var globalHandler http.Handler
	if len(cfg.Redirect) == 0 {
		globalHandler = mux
	} else {
		globalHandler = legacy.NewRedirector(cfg.Redirect, mux)
	}

	log.Printf("Redirecting insecure traffic from port %v", cfg.InsecurePort)
	go func() {
		panic(http.ListenAndServe(fmt.Sprintf(":%d", cfg.InsecurePort),
			conn.InsecureRedirection(cfg.Domain, cfg.TLSPort)))
	}()

	log.Printf("GoBlog is listening at port %v", cfg.TLSPort)
	panic(conn.ListenAndServeTLS(cfg.TLSPort, cfg.TLSCertPath, cfg.TLSKeyPath, globalHandler))
}
