package feed

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/mariomac/goblog/src/blog"
	"golang.org/x/tools/blog/atom"
)

// BuildAtomFeed builds an XML Atom feed from an ordered (from new to old) list of blog entries
func BuildAtomFeed(bentries []blog.Entry, hostname string, entrypath string) string {
	entries := make([]*atom.Entry, len(bentries))

	for i, bentry := range bentries {
		entries[i] = &atom.Entry{
			Title: bentry.Title,
			ID:    fmt.Sprint(bentry.Time.Unix()),
			Link: []atom.Link{
				{Href: "http://" + hostname + entrypath + bentry.FileName},
			},
			Published: atom.Time(*bentry.Time),
			Summary: &atom.Text{
				Type: "text/html",
				Body: string(bentry.Preview),
			},
		}
	}

	feed := atom.Feed{
		Title: "Entries for " + hostname,
		ID:    hostname,
		Link: []atom.Link{
			{Href: "http://" + hostname},
		},
		Updated: atom.Time(*bentries[0].Time),
		Entry:   entries,
	}

	out := make([]byte, 0, 2048)
	buf := bytes.NewBuffer(out)
	encoder := xml.NewEncoder(buf)
	encoder.Encode(feed)
	return buf.String()
}
