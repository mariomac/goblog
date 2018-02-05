package bentry

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TEST_RESOURCES = "../../testresources/testentries"

func TestBlogContent_LoadEntries(t *testing.T) {
	blog := new(BlogContent)
	blog.Load(TEST_RESOURCES)

	assert := assert.New(t)

	assert.Equal(3, len(blog.entries))

	assert.Equal("201710281345_gurbai.md", blog.entries[0].FileName)
	assert.Equal("Gurbai!", blog.entries[0].Title)
	assert.Equal(time.Date(2017, 10, 28, 13, 45, 0, 0, location), *blog.entries[0].Time)

	assert.Equal("201709281345_hello-my-frens.md", blog.entries[1].FileName)
	assert.Equal("Hello my frens!", blog.entries[1].Title)
	assert.Equal(time.Date(2017, 9, 28, 13, 45, 0, 0, location), *blog.entries[1].Time)

	assert.Equal("201610281345_hello_guy.md", blog.entries[2].FileName)
	assert.Equal("Hello guy!", blog.entries[2].Title)
	assert.Equal(time.Date(2016, 10, 28, 13, 45, 0, 0, location), *blog.entries[2].Time)
}

func TestBlogContent_LoadAll(t *testing.T) {
	blog := new(BlogContent)
	blog.Load(TEST_RESOURCES)

	assert := assert.New(t)

	assert.Equal(5, len(blog.all))

	assert.Equal("201710281345_gurbai.md", blog.all["201710281345_gurbai.md"].FileName)
	assert.Equal("Gurbai!", blog.all["201710281345_gurbai.md"].Title)
	assert.Equal(time.Date(2017, 10, 28, 13, 45, 0, 0, location),
		*blog.all["201710281345_gurbai.md"].Time)

	assert.Equal("201709281345_hello-my-frens.md",
		blog.all["201709281345_hello-my-frens.md"].FileName)
	assert.Equal("Hello my frens!", blog.all["201709281345_hello-my-frens.md"].Title)
	assert.Equal(time.Date(2017, 9, 28, 13, 45, 0, 0, location),
		*blog.all["201709281345_hello-my-frens.md"].Time)

	assert.Equal("201610281345_hello_guy.md", blog.all["201610281345_hello_guy.md"].FileName)
	assert.Equal("Hello guy!", blog.all["201610281345_hello_guy.md"].Title)
	assert.Equal(time.Date(2016, 10, 28, 13, 45, 0, 0, location),
		*blog.all["201610281345_hello_guy.md"].Time)

	assert.Equal("201610281345_hello_guy.md", blog.all["201610281345_hello_guy.md"].FileName)
	assert.Equal("Hello guy!", blog.all["201610281345_hello_guy.md"].Title)
	assert.Equal(time.Date(2016, 10, 28, 13, 45, 0, 0, location),
		*blog.all["201610281345_hello_guy.md"].Time)

	assert.Equal("gurbai.md", blog.all["gurbai.md"].FileName)
	assert.Equal("Gurbai page!", blog.all["gurbai.md"].Title)
	assert.Nil(blog.all["gurbai.md"].Time)

	assert.Equal("hello.md", blog.all["hello.md"].FileName)
	assert.Equal("Hello page!", blog.all["hello.md"].Title)
	assert.Nil(blog.all["hello.md"].Time)
}

func TestExtractTime(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(
		time.Date(1979, 5, 25, 06, 07, 0, 0, location),
		extractTime("197905250607"), "YYYYMMDDHHMM dates should be parsed correctly")
}

func TestGetTitleBodyAndPreview(t *testing.T) {
	assert := assert.New(t)

	title, body, preview := getTitleBodyAndPreview([]byte(
		`# This is a title

This is a paragraph

This is another paragraph`))

	assert.Equal("This is a title", title, "Title is not well extracted")
	assert.True(strings.Contains(string(body), "This is a paragraph"))
	assert.False(strings.Contains(string(body), "This is a title"), "Title should have been removed")
	assert.False(strings.Contains(string(body), "<h1>"), "H1 should have been removed")

	assert.True(strings.Contains(string(preview), "This is a paragraph"))
	assert.False(strings.Contains(string(preview), "This is a title"), "Title should have been removed")
	assert.False(strings.Contains(string(preview), "This is another paragraph"), "Only first paragraph should be in preview")

}
