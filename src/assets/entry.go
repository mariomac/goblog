package assets

import (
	"bytes"
	"fmt"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/logr"
	"github.com/mariomac/goblog/src/visual"
)

const (
	entryMimeType = indexMimeType
)

type EntryGenerator struct {
	entries   *blog.Entries
	templates visual.Templater
}

func (e *EntryGenerator) Get(urlPath string) (*WebAsset, error) {
	file := urlPath
	// TODO: extra fields. E.g. source IP
	entry, found := e.entries.Get(file)
	if !found {
		return nil, errNotFound{}
	}
	body := bytes.Buffer{}
	if err := e.templates.Render(visual.EntryTemplate, entry, &body); err != nil {
		logr.Get().Error("rendering entry template",
			"error", err,
			"urlPath", urlPath,
			"fileName", file,
		)
		return nil, internalError{cause: fmt.Errorf("rendering entry template: %w", err)}
	}
	return &WebAsset{
		MimeType: entryMimeType,
		Body:     body.Bytes(),
	}, nil
}
