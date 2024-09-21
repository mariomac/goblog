package blog

import (
	"log/slog"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mariomac/goblog/src/logr"
)

func TestBlogContent_LoadAll(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	blog, err := PreloadEntries(testResources)
	require.NoError(t, err)

	require.Len(t, blog.all, 5)

	assert.Equal(t, "201710281345_gurbai.md", blog.all["201710281345_gurbai.md"].FileName)
	assert.Equal(t, "Gurbai!", blog.all["201710281345_gurbai.md"].Title)
	assert.Equal(t, time.Date(2017, 10, 28, 13, 45, 0, 0, location),
		blog.all["201710281345_gurbai.md"].Time)

	assert.Equal(t, "201709281345_hello-my-frens.md",
		blog.all["201709281345_hello-my-frens.md"].FileName)
	assert.Equal(t, "Hello my frens!", blog.all["201709281345_hello-my-frens.md"].Title)
	assert.Equal(t, time.Date(2017, 9, 28, 13, 45, 0, 0, location),
		blog.all["201709281345_hello-my-frens.md"].Time)

	assert.Equal(t, "201610281345_hello_guy.md", blog.all["201610281345_hello_guy.md"].FileName)
	assert.Equal(t, "Hello guy!", blog.all["201610281345_hello_guy.md"].Title)
	assert.Equal(t, time.Date(2016, 10, 28, 13, 45, 0, 0, location),
		blog.all["201610281345_hello_guy.md"].Time)

	assert.Equal(t, "gurbai.md", blog.all["gurbai.md"].FileName)
	assert.Equal(t, "Gurbai page!", blog.all["gurbai.md"].Title)
	assert.True(t, blog.all["gurbai.md"].Time.IsZero())

	assert.Equal(t, "hello.md", blog.all["hello.md"].FileName)
	assert.Equal(t, "Hello page!", blog.all["hello.md"].Title)
	assert.True(t, blog.all["hello.md"].Time.IsZero())

	// Check that entries are sorted by timestamp
	require.Len(t, blog.sorted, 3)
	assert.Equal(t, "201710281345_gurbai.md", blog.sorted[0].FileName)
	assert.Equal(t, "201709281345_hello-my-frens.md", blog.sorted[1].FileName)
	assert.Equal(t, "201610281345_hello_guy.md", blog.sorted[2].FileName)
}

func TestPager(t *testing.T) {
	logr.Init(slog.LevelDebug)
	entries := Entries{sorted: []*Entry{
		{Time: time.Date(2021, 10, 12, 0, 0, 0, 0, location)},
		{Time: time.Date(2021, 10, 11, 0, 0, 0, 0, location)},
		{Time: time.Date(2021, 10, 10, 0, 0, 0, 0, location)},
		{Time: time.Date(2021, 10, 9, 0, 0, 0, 0, location)},
		{Time: time.Date(2021, 10, 8, 0, 0, 0, 0, location)},
		{Time: time.Date(2021, 10, 7, 0, 0, 0, 0, location)},
		{Time: time.Date(2021, 10, 6, 0, 0, 0, 0, location)},
		{Time: time.Date(2021, 10, 5, 0, 0, 0, 0, location)},
	}}
	assert.Equal(t, entries.sorted[0:3], entries.Sorted(0, 3))
	assert.Equal(t, entries.sorted[3:6], entries.Sorted(1, 3))
	assert.Equal(t, entries.sorted[6:8], entries.Sorted(2, 3))
	assert.Empty(t, entries.Sorted(3, 3))
}
