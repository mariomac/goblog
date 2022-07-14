package blog

import (
	"errors"
	"io/fs"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testResources = "../../testresources/testblog/entries"

func TestBlogContent_LoadEntries(t *testing.T) {
	type testCase struct {
		file string
		title string
		ts time.Time
	}
	for _, tc := range []testCase{
		{file: "201710281345_gurbai.md", title: "Gurbai!", ts: time.Date(2017, 10, 28, 13, 45, 0, 0, location)},
		{file: "201709281345_hello-my-frens.md", title: "Hello my frens!", ts: time.Date(2017, 9, 28, 13, 45, 0, 0, location)},
		{file: "201610281345_hello_guy.md", title: "Hello guy!", ts: time.Date(2016, 10, 28, 13, 45, 0, 0, location)},
		{file: "hello.md", title: "Hello page!"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			entry, err := LoadEntry(path.Join(testResources, tc.file))
			require.NoError(t, err)
			assert.Equal(t, tc.title, entry.Title)
			assert.Equal(t, tc.file, entry.FileName)
			assert.Equal(t, tc.ts, entry.Time)
		})
	}
}

func TestFileNotFound(t *testing.T) {
	_, err := LoadEntry("foobar.md")
	require.True(t, errors.Is(err, fs.ErrNotExist))
}

// TODO: test not found

func TestExtractTime(t *testing.T) {
	assert.Equal(t,
		time.Date(1979, 5, 25, 06, 07, 0, 0, location),
		extractTime("197905250607"), "YYYYMMDDHHMM dates should be parsed correctly")
}

func TestGetTitleBodyAndPreview(t *testing.T) {
	title, body, preview := getTitleBodyAndPreview([]byte(
		`# This is a title

This is a paragraph

This is another paragraph`))

	assert.Equal(t, "This is a title", title, "Title is not well extracted")
	assert.True(t, strings.Contains(string(body), "This is a paragraph"))
	assert.False(t, strings.Contains(string(body), "This is a title"), "Title should have been removed")
	assert.False(t, strings.Contains(string(body), "<h1>"), "H1 should have been removed")

	assert.True(t, strings.Contains(string(preview), "This is a paragraph"))
	assert.False(t, strings.Contains(string(preview), "This is a title"), "Title should have been removed")
	assert.False(t, strings.Contains(string(preview), "This is another paragraph"), "Only first paragraph should be in preview")

}
