// Package legacy implements mechanisms for backwards compatibility with legacy systems
package legacy

import "net/http"

// Redirector middleware of old URLS to their new location. If no legacy urls are matched,
// it forwards the request to its wrapped handler
type Redirector struct {
	urls map[string]string
	next http.Handler
}

func NewRedirector(urls map[string]string, next http.Handler) *Redirector {
	return &Redirector{
		urls: urls,
		next: next,
	}
}

func (red *Redirector) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if !red.tryRedirection(res, req) {
		red.next.ServeHTTP(res, req)
	}
}

func (red *Redirector) tryRedirection(res http.ResponseWriter, req *http.Request) bool {
	if len(red.urls) == 0 || req.URL == nil {
		return false
	}
	if newPath, ok := red.urls[req.URL.Path]; ok {
		req.URL.Path = newPath
		res.Header().Set("Location", req.URL.String())
		res.WriteHeader(http.StatusMovedPermanently)
		return true
	}
	return false
}
