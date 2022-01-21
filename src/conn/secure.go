// Package conn manages secure and insecure (redirect) connections
package conn

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func ListenAndServeTLS(port int, handler http.Handler) error {
	cert, key := createLocalCerts()
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
