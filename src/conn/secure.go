// Package conn manages secure and insecure (redirect) connections
package conn

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/mariomac/goblog/src/logr"
)

func ListenAndServeTLS(port int, cert, key string, handler http.Handler) error {
	if key == "" || cert == "" {
		logr.Get().Warn("creating insecure local certificates for localhost development only")
		cert, key = createLocalCerts()
	}
	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	return server.ListenAndServeTLS(cert, key)
}

func InsecureRedirection(hostName string, redirectPort int) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		oldURL := req.URL.String()
		req.URL.Scheme = "https"
		// if the redirection port is the standard HTTPS port, we don't attach it
		// to the redirection URL
		if redirectPort == 443 {
			req.URL.Host = hostName
		} else {
			req.URL.Host = fmt.Sprintf("%v:%d", hostName, redirectPort)
		}
		newURL := req.URL.String()
		log.Println("redirecting:", oldURL, "->", newURL)
		rw.Header().Set("Location", newURL)
		rw.WriteHeader(http.StatusMovedPermanently)
	}
}
