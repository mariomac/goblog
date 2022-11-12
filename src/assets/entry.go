package assets

import (
	"bytes"
	"fmt"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/logr"
	"github.com/mariomac/goblog/src/visual"
	"github.com/sirupsen/logrus"
)

const (
	entryMimeType = indexMimeType
)

var elog = logr.Get()

type EntryGenerator struct {
	entries   *blog.Entries
	templates visual.Templater
}

func (e *EntryGenerator) Get(urlPath string) (*WebAsset, error) {
	file := urlPath[len(pathEntry):]
	// TODO: extra fields. E.g. source IP
	entry, found := e.entries.Get(file)
	if !found {
		return nil, errNotFound{}
	}
	body := bytes.Buffer{}
	if err := e.templates.Render(visual.EntryTemplate, entry, &body); err != nil {
		elog.WithFields(logrus.Fields{
			logrus.ErrorKey: err,
			"urlPath":       urlPath,
			"fileName":      file,
		}).Error("rendering entry template")
		return nil, internalError{cause: fmt.Errorf("rendering entry template: %w", err)}
	}
	return &WebAsset{
		MimeType: entryMimeType,
		Body:     body.Bytes(),
	}, nil
}
