package assets

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/visual"
)

const (
	indexMimeType = "text/html; charset=utf-8"
)

type IndexGenerator struct {
	entries        *blog.Entries
	entriesPerPage int
	templates      *visual.Templater
}

type IndexRender struct {
	PreviousPage int
	CurrentPage  int
	NextPage     int
	TotalPages   int
	Entries      []*blog.Entry
}

func (i *IndexGenerator) Get(page string) (*WebAsset, error) {
	ir := IndexRender{
		TotalPages:  i.entries.Len()/i.entriesPerPage + 1,
		CurrentPage: 1,
	}
	if len(page) != 0 {
		var err error
		ir.CurrentPage, err = strconv.Atoi(strings.SplitN(page, "/", 2)[0])
		if err != nil || ir.CurrentPage <= 0 {
			return nil, errNotFound{}
		}
	}
	ir.PreviousPage, ir.NextPage = ir.CurrentPage-1, ir.CurrentPage+1
	ir.Entries = i.entries.Sorted(ir.CurrentPage-1, i.entriesPerPage)
	body := bytes.Buffer{}
	if err := i.templates.Render(visual.IndexTemplate, ir, &body); err != nil {
		return nil, internalError{cause: fmt.Errorf("rendering index template: %w", err)}
	}
	return &WebAsset{
		MimeType: indexMimeType,
		Body:     body.Bytes(),
	}, nil
}
