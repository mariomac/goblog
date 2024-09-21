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
	"github.com/mariomac/goblog/src/conn"
	"github.com/mariomac/goblog/src/fs"
	"github.com/mariomac/goblog/src/install"
	"github.com/mariomac/goblog/src/legacy"
	"github.com/mariomac/goblog/src/logr"
)

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
	logr.Init(cfg.LogLevel)
	log := logr.Get()
	log.Debug("configuration loaded", "config", fmt.Sprintf("%#v", cfg))

	log.Info("Starting GoBlog...")
	// TODO: allow insecure traffic
	mux, err := assets.NewCachedHandler(&cfg, true)
	if err != nil {
		log.Error("can't start blog handler. ABORTING",
			"error", err,
			"rootPath", cfg.RootPath,
			"domain", cfg.Domain,
		)
		os.Exit(1)
	}

	if err := fs.NotifyChanges(cfg.RootPath, mux.Reload); err != nil {
		log.Warn("could not listen for file changes. Your blog won't be "+
			"automatically updated if you change any file", "error", err)
	}

	var globalHandler http.HandlerFunc
	if len(cfg.Redirect) == 0 {
		globalHandler = mux.ServeHTTP
	} else {
		globalHandler = legacy.NewRedirector(cfg.Redirect, mux).ServeHTTP
	}

	if cfg.MaxRequests.Number > 0 && cfg.MaxRequests.Period > 0 {
		globalHandler = conn.ClientRateLimitHandler(globalHandler,
			cfg.MaxRequests.Number, cfg.MaxRequests.Period, cfg.MaxRequests.Period)
	}

	go func() {
		if cfg.HTTPSRedirect {
			log.Info("Redirecting insecure traffic",
				"srcPort", cfg.InsecurePort, "dstPort", cfg.TLSPort)
			panic(http.ListenAndServe(fmt.Sprintf(":%d", cfg.InsecurePort),
				conn.InsecureRedirection(cfg.Domain, cfg.TLSPort)))
		} else if cfg.InsecurePort > 0 {
			log.Info("Working with insecure port", "port", cfg.InsecurePort)
			panic(http.ListenAndServe(fmt.Sprintf(":%d", cfg.InsecurePort), globalHandler))
		}
	}()

	if cfg.TLSPort > 0 {
		log.Info("Working with secure port", "port", cfg.TLSPort)
		panic(conn.ListenAndServeTLS(cfg.TLSPort, cfg.TLSCertPath, cfg.TLSKeyPath, globalHandler))
	}
	<-make(chan struct{})
}
