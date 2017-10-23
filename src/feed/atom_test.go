package feed

import (
	"testing"
	"../bentry"
	assert2 "github.com/stretchr/testify/assert"
)

const TEST_RESOURCES = "../../testresources/testentries"

func TestBuildAtomFeed(t *testing.T) {
	assert := assert2.New(t)

	entries := bentry.BlogContent{}
	entries.Load(TEST_RESOURCES)

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
		"</entry>" +
		"<entry>" +
		"<title>Hello my frens!</title>" +
		"<id>1506599100</id>" +
		"<link href=\"http://www.superblog.com/entry/201709281345_hello-my-frens.md\"></link>" +
		"<published>2017-09-28T13:45:00+02:00</published><updated></updated>" +
		"</entry>" +
		"<entry>" +
		"<title>Hello guy!</title>" +
		"<id>1477655100</id>" +
		"<link href=\"http://www.superblog.com/entry/201610281345_hello_guy.md\"></link>" +
		"<published>2016-10-28T13:45:00+02:00</published><updated></updated>" +
		"</entry>" +
		"</feed>"

	assert.Equal(expected, atomxml, "Generated atom feed does not match")
}
