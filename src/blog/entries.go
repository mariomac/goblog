package blog

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/mariomac/goblog/src/fs"
	"github.com/mariomac/goblog/src/logr"
)

var elog = logr.Get()

var anyPageFormat = regexp.MustCompile(`\.md$`)

type Entries struct {
	pageSize int
	sorted   []*Entry          // only timestamped entries, sorted from new to old
	all      map[string]*Entry // all the entries and pages, accessible by FileName
}

func (e *Entries) Sorted(pageNum, pageSize int) []*Entry {
	startIdx := pageNum * pageSize
	if startIdx >= len(e.sorted) {
		return nil
	}
	endIdx := startIdx + pageSize
	if endIdx > len(e.sorted) {
		endIdx = len(e.sorted)
	}
	return e.sorted[startIdx:endIdx]
}

func (e *Entries) Get(fileName string) (*Entry, bool) {
	entry, ok := e.all[fileName]
	return entry, ok
}

func PreloadEntries(directory string) (Entries, error) {
	e := Entries{all: map[string]*Entry{}}
	plog := elog.WithField("dir", directory)
	plog.Info("loading all blog entries and pages")

	files, err := fs.Search(directory, anyPageFormat)
	if err != nil {
		return e, fmt.Errorf("loading pages from directory %s: %w", directory, err)
	}
	for _, file := range files {
		plog.WithField("filePath", file).Debug("found file entry")
		entry, err := LoadEntry(file)
		if err != nil {
			plog.WithError(err).Warn("can't load blog entry. Ignoring")
			continue
		}
		e.all[entry.FileName] = entry
		// Timestamped entries will be sorted as blog entries
		if !entry.Time.IsZero() {
			e.sorted = append(e.sorted, entry)
		}
	}

	sort.Slice(e.sorted, func(i, j int) bool {
		return e.sorted[i].Time.Sub(e.sorted[j].Time) > 0
	})

	return e, nil
}
