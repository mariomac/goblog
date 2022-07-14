package assets

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/mariomac/goblog/src/blog"
	"golang.org/x/tools/blog/atom"
)

// TODO: make configurable
const (
	atomMaxLength = 10
	atomMimeType  = "application/atom+xml"
)

type AtomGenerator struct {
	urlProtocol string
	hostName    string
	entryPath   string
	entries     *blog.Entries
}

func (a *AtomGenerator) Get(_ string) (*WebAsset, error) {
	bEntries := a.entries.Sorted(0, atomMaxLength)

	entries := make([]*atom.Entry, 0, len(bEntries))

	for _, bentry := range bEntries {
		entries = append(entries, &atom.Entry{
			Title: bentry.Title,
			ID:    fmt.Sprint(bentry.Time.Unix()),
			Link: []atom.Link{
				{Href: "http://" + a.hostName + a.entryPath + bentry.FileName},
			},
			Published: atom.Time(bentry.Time),
			Summary: &atom.Text{
				Type: "text/html",
				Body: string(bentry.Preview),
			},
		})
	}

	feed := atom.Feed{
		Title: "Entries for " + a.hostName,
		ID:    a.hostName,
		Link: []atom.Link{
			{Href: a.urlProtocol + a.hostName},
		},
		Updated: atom.Time(bEntries[0].Time),
		Entry:   entries,
	}

	buf := bytes.Buffer{}
	encoder := xml.NewEncoder(&buf)
	if err := encoder.Encode(feed); err != nil {
		return nil, internalError{cause: fmt.Errorf("encoding atom: %w", err)}
	}
	return &WebAsset{
		MimeType: atomMimeType,
		Body:     buf.Bytes(),
	}, nil
}
