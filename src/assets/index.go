package assets

import (
	"bytes"
	"fmt"
	"math"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/visual"
)

const (
	indexMimeType = "text/html; charset=utf-8"
)

type IndexGenerator struct{
	entries *blog.Entries
	templates *visual.Templater
}

func (i *IndexGenerator) Get(_ string) (*WebAsset, error) {
	// TODO: properly paginate entries

	body := bytes.Buffer{}
	if err := i.templates.Render(visual.IndexTemplate, i.entries.Sorted(0, math.MaxInt), &body);
	err != nil {
		return nil, internalError{cause: fmt.Errorf("rendering index template: %w", err)}
	}
	return &WebAsset{
		MimeType: indexMimeType,
		Body: body.Bytes(),
	}, nil
}

