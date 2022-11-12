package assets

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/logr"
	"github.com/mariomac/goblog/src/visual"
	"github.com/mariomac/guara/pkg/cache"
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

var unsupportedMethodErr = errors.New("unsupported method")

var alog = logr.Get()

type WebAsset struct {
	MimeType string
	Body     []byte
}

func (w *WebAsset) SizeBytes() int {
	return len(w.Body)
}

type webAssetGenerator interface {
	// Get returns a webasset given the urlPath. The urlPath function removed the parent route
	// that led to the given asset generator
	Get(urlPath string) (*WebAsset, error)
}

type route struct {
	Prefix    string
	Generator webAssetGenerator
}

type CachedHandler struct {
	rootPath      string
	hostName      string
	maxCacheBytes int
	tls           bool
	assets        *cache.LRU[string, *WebAsset]
	routes        []route
}

const entriesPerPage = 5

// TODO pass "routedHandler" as argument and remove router logic from here
func NewCachedHandler(rootPath string, isTLS bool, hostName string, maxCacheBytes int) (*CachedHandler, error) {
	cc := &CachedHandler{
		rootPath:      rootPath,
		hostName:      hostName,
		tls:           isTLS,
		maxCacheBytes: maxCacheBytes,
	}
	if err := cc.Reload(); err != nil {
		return nil, fmt.Errorf("loading resources: %w", err)
	}
	return cc, nil
}

func (c *CachedHandler) Reload() error {
	entries, err := blog.PreloadEntries(path.Join(c.rootPath, dirEntry))
	if err != nil {
		return fmt.Errorf("loading blog entries: %w", err)
	}

	templates, err := visual.LoadTemplates(path.Join(c.rootPath, dirTemplate))
	if err != nil {
		return fmt.Errorf("loading template: %w", err)
	}
	protocol := "http://"
	if c.tls {
		protocol = "https://"
	}
	c.assets = cache.NewLRU[string, *WebAsset](c.maxCacheBytes)
	c.routes = []route{
		{Prefix: pathStatic, Generator: &FileAssetGenerator{rootPath: c.rootPath}},
		{Prefix: pathEntry, Generator: &EntryGenerator{templates: templates, entries: &entries}},
		{Prefix: pathAtom, Generator: &AtomGenerator{
			urlProtocol: protocol, hostName: c.hostName, entryPath: pathEntry, entries: &entries}},
		{Prefix: pathIndex, Generator: &IndexGenerator{entries: &entries, templates: &templates, entriesPerPage: entriesPerPage}},
	}
	return nil
}

func (c *CachedHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// TODO: instrument cache size in bytes
	alog := alog.WithFields(logrus.Fields{
		"method":     request.Method,
		"url":        request.URL,
		"remoteAddr": request.RemoteAddr,
	})
	alog.Debug("new request")
	if request.Method != http.MethodGet {
		writeErr(http.StatusBadRequest, unsupportedMethodErr, writer, request)
		return
	}
	fileUrlPath := path.Clean(request.URL.Path)
	if asset, ok := c.assets.Get(fileUrlPath); ok {
		alog.Debug("found cached copy")
		writeAsset(writer, asset, alog)
		return
	}
	for _, r := range c.routes {
		if strings.HasPrefix(fileUrlPath, r.Prefix) {
			asset, err := r.Generator.Get(fileUrlPath[len(r.Prefix):])
			if err != nil {
				switch e := err.(type) {
				case errNotFound:
					e.url = fileUrlPath
					writeErr(http.StatusNotFound, err, writer, request)
				default:
					writeErr(http.StatusInternalServerError, err, writer, request)
				}
				return
			}
			writeAsset(writer, asset, alog)
			c.assets.Put(fileUrlPath, asset)
			return
		}
	}
	writeErr(http.StatusNotFound, errNotFound{url: request.URL.String()}, writer, request)
}

func writeAsset(writer http.ResponseWriter, asset *WebAsset, alog *logrus.Entry) {
	writer.Header().Set("Content-Type", asset.MimeType)
	if _, err := writer.Write(asset.Body); err != nil {
		alog.WithFields(logrus.Fields{
			logrus.ErrorKey: err,
			"contentType":   asset.MimeType,
		}).Error("couldn't write response")
	}
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
