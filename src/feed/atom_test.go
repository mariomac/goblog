package feed

import (
	"testing"

	"github.com/mariomac/goblog/src/blog"
	assert2 "github.com/stretchr/testify/assert"
)

const testResources = "../../testresources/testentries"

func TestBuildAtomFeed(t *testing.T) {
	assert := assert2.New(t)

	entries := blog.Content{}
	entries.Load(testResources)

	atomxml := BuildAtomFeed(entries.GetEntries(), "www.superblog.com", "/entry/")

	expected := "" +
		"<feed xmlns=\"http://www.w3.org/2005/Atom\">" +
		"<title>Entries for www.superblog.com</title>" +
		"<id>www.superblog.com</id>" +
		"<link href=\"http://www.superblog.com\"></link>" +
		"<updated>2017-10-28T13:45:00+02:00</updated>" +
		"<entry>" +
		"<title>Gurbai!</title><id>1509191100</id>" +
		"<link href=\"http://www.superblog.com/entry/201710281345_gurbai.md\"></link>" +
		"<published>2017-10-28T13:45:00+02:00</published><updated></updated>" +
		"<summary type=\"text/html\">&lt;p&gt;Gurbai!&lt;/p&gt;</summary>" +
		"</entry>" +
		"<entry>" +
		"<title>Hello my frens!</title>" +
		"<id>1506599100</id>" +
		"<link href=\"http://www.superblog.com/entry/201709281345_hello-my-frens.md\"></link>" +
		"<published>2017-09-28T13:45:00+02:00</published><updated></updated>" +
		"<summary type=\"text/html\">&lt;p&gt;&lt;img src=&#34;/static/img.png&#34; alt=&#34;Image&#34;/&gt;&lt;/p&gt;</summary>" +
		"</entry>" +
		"<entry>" +
		"<title>Hello guy!</title>" +
		"<id>1477655100</id>" +
		"<link href=\"http://www.superblog.com/entry/201610281345_hello_guy.md\"></link>" +
		"<published>2016-10-28T13:45:00+02:00</published><updated></updated>" +
		"<summary type=\"text/html\">&lt;p&gt;Paragraph of hello guy&lt;/p&gt;</summary>" +
		"</entry>" +
		"</feed>"

	assert.Equal(expected, atomxml, "Generated atom feed does not match")
}
