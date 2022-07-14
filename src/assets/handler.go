package assets

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/logr"
	"github.com/mariomac/goblog/src/visual"
	"github.com/sirupsen/logrus"
)

// Path names
const (
	pathStatic = "/static/"
	pathEntry  = "/entry/"
	pathIndex  = "/"
	pathAtom   = "/atom.xml"

	dirTemplate = "template/"
	dirEntry    = "entries/"
)

var alog = logr.Get()

type WebAsset struct {
	MimeType string
	Body     []byte
}

type webAssetGenerator interface {
	Get(urlPath string) (*WebAsset, error)
}

type route struct {
	Prefix    string
	Generator webAssetGenerator
}

type CachedHandler struct {
	routes []route
}

func CreateCachedHandler(rootPath string, isTLS bool, hostName string) (CachedHandler, error) {
	entries, err := blog.PreloadEntries(path.Join(rootPath, dirEntry))
	if err != nil {
		return CachedHandler{}, fmt.Errorf("loading blog entries: %w", err)
	}

	templates, err := visual.LoadTemplates(path.Join(rootPath, dirTemplate))
	if err != nil {
		return CachedHandler{}, fmt.Errorf("loading templates: %w", err)
	}
	protocol := "http://"
	if isTLS {
		protocol = "https://"
	}
	return CachedHandler{routes: []route{
		{Prefix: pathStatic, Generator: &FileAssetGenerator{rootPath: rootPath}},
		{Prefix: pathEntry, Generator: &EntryGenerator{templates: templates, entries: &entries}},
		{Prefix: pathAtom, Generator: &AtomGenerator{
			urlProtocol: protocol, hostName: hostName, entryPath: pathEntry, entries: &entries}},
		{Prefix: pathIndex, Generator: &IndexGenerator{entries: &entries, templates: &templates}},
	}}, nil
}

func (c *CachedHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// TODO: test e.g. leading spaces in URL
	// TODO: check cache
	path := request.URL.Path
	for _, r := range c.routes {
		if strings.HasPrefix(path, r.Prefix) {
			asset, err := r.Generator.Get(path)
			if err != nil {
				switch err.(type) {
				case errNotFound:
					writeErr(http.StatusNotFound, err, writer, request)
				default:
					writeErr(http.StatusInternalServerError, err, writer, request)
				}
				return
			}
			writer.Header().Set("Content-Type", asset.MimeType)
			if _, err := writer.Write(asset.Body); err != nil {
				alog.WithFields(logrus.Fields{
					logrus.ErrorKey: err,
					"url":           request.URL,
					"contentType":   asset.MimeType,
				}).Error("couldn't write response")
			}
			// TODO: update cache
			return
		}
	}
	writeErr(http.StatusNotFound, errNotFound{url: request.URL.String()}, writer, request)
}

func writeErr(code int, err error, writer http.ResponseWriter, request *http.Request) {
	// TODO: provide a proper internal error page
	writer.WriteHeader(code)
	if _, werr := writer.Write([]byte(err.Error())); werr != nil {
		alog.WithFields(logrus.Fields{
			logrus.ErrorKey: werr,
			"cause":         err,
			"url":           request.URL,
			"statusCode":    code,
		}).Warn("couldn't write response error message")
	}
}
