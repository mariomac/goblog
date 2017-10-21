package bentry

import (
	"regexp"
	"html/template"
	"../fs"
	"log"
)

// Pages, without timestamp
type Page struct {
	UrlPath string
	Title string
	Html *template.HTML
}

// Timestamped entries
type Entry struct {
	Page
	Time string
}

type BlogContent struct {
	entries map[string]Entry
	pages map[string]Page
}

// YYYYMMDDHHMMsome-text_here.md
var entryFormat = regexp.MustCompile("^[0-9]{12}[_\\-a-zA-Z0-9]+\\.md$")
var allFormat = regexp.MustCompile("^[_\\-a-zA-Z0-9]+\\.md$")

func (blog *BlogContent) Load(folder string) {
	log.Printf("Scanning for entries in folder %s...", folder)
	files := fs.Search(folder, allFormat)
	for _, entry := range files {
		log.Printf("Entry found: %s", entry)
	}
}
