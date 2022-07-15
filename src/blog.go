// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on 'Writing web applications' tutorial: https://golang.org/doc/articles/wiki/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mariomac/goblog/src/assets"
	"github.com/mariomac/goblog/src/install"
	"github.com/mariomac/goblog/src/legacy"
	"github.com/mariomac/goblog/src/logr"
	"github.com/sirupsen/logrus"

	"github.com/mariomac/goblog/src/conn"
)

var log *logrus.Entry

func init() {
	// TODO: make log level configurable
	logrus.SetLevel(logrus.DebugLevel)
	log = logr.Get()
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
	// TODO: allow insecure traffic
	mux, err := assets.NewCachedHandler(cfg.RootPath, true, cfg.Domain)
	if err != nil {
		log.WithFields(logrus.Fields{
			logrus.ErrorKey: err,
			"rootPath":      cfg.RootPath,
			"domain":        cfg.Domain,
		}).Fatal("can't start blog handler")
	}

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
