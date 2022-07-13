// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/mariomac/goblog/src/install"
	"github.com/mariomac/goblog/src/legacy"
	"github.com/mariomac/goblog/src/logr"
	"github.com/sirupsen/logrus"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/conn"
	"github.com/mariomac/goblog/src/feed"
	"github.com/mariomac/goblog/src/visual"
)

var log *logrus.Entry

func init() {
	// TODO: make log level configurable
	logrus.SetLevel(logrus.DebugLevel)
	log = logr.Get()
}


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

var entries blog.Entries
var templates = visual.Templater{}

func viewHandler(w http.ResponseWriter, _ *http.Request, fileName string) {
	// TODO: extra fields. E.g. source IP
	log.WithField("fileName", fileName).Debug("view handler")
	entry, found := entries.Get(fileName)
	if !found {
		http.Error(w, "Entry not found "+fileName, http.StatusNotFound) // Todo: redirect or template 404
		return
	}
	if err := templates.Render(visual.EntryTemplate, entry, w); err != nil {
		log.WithFields(logrus.Fields{
			logrus.ErrorKey: err,
			"fileName": fileName,
		}).Error("rendering entry template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func makeIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: pass entries here instead of embedding it as a "getTemplate" function
		if err := templates.Render(visual.IndexTemplate, nil, w); err != nil {
			log.WithError(err).Error("rendering index template")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

var validPagePath = regexp.MustCompile(`^([_\-a-zA-Z0-9]+)\.md$`)

func makePageHandler(
	rootPath string,
	fn func(http.ResponseWriter, *http.Request, string),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPagePath.FindStringSubmatch(r.URL.Path[len(rootPath):])
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[0])
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
	entries, err = blog.PreloadEntries(path.Join(cfg.RootPath, dirEntry))
	if err != nil {
		log.WithError(err).WithField("directory", path.Join(cfg.RootPath, dirEntry)).
			Fatal("can't load log entries")
	}

	// Create Atom XML feed
	atomxml := bytes.NewBufferString(
		feed.BuildAtomFeed(entries.Sorted(0, math.MaxInt), cfg.Domain, pathEntry)).Bytes()

	templates, err = visual.LoadTemplates(path.Join(cfg.RootPath, dirTemplate), func() []*blog.Entry {
		// TODO: paginate well
		return entries.Sorted(0, math.MaxInt)
	})
	if err != nil {
		log.WithError(err).WithField("directory", path.Join(cfg.RootPath, dirTemplate)).
			Fatal("can't load templates")
	}

	mux := http.NewServeMux()
	mux.HandleFunc(pathIndex, makeIndexHandler())
	mux.HandleFunc(pathAtom, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/xml")
		// TODO: handle error
		writer.Write(atomxml)
	})
	mux.HandleFunc(pathEntry, makePageHandler(pathEntry, viewHandler))
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
