package btemplate

import (
	"regexp"
	"html/template"
	"../fs"
	"../filesearch"
	"net/http"
	"github.com/shurcooL/github_flavored_markdown"
	"log"
)

type Templates struct {
	entries *template.Template
}

var validTemplate = regexp.MustCompile(".*\\.html$")

const TMPL_EXT = ".html"

// Todo: retrigger on folder change
func (t *Templates) Load(folder string, getEntries func() []*filesearch.Entry) {
	log.Printf("Scanning for templates in folder %s...\n", folder)

	templateFiles := fs.Search(folder, validTemplate)
	for _, f := range templateFiles {
		log.Printf("Template file found: %s\n", f)
	}

	t.entries = template.Must(
		template.New("golog_templates").Funcs(
			template.FuncMap{"entries": getEntries, "md2html": md2html}).ParseFiles(
			templateFiles...))
}

func (t *Templates) Render(w http.ResponseWriter, template string, data interface{}) {
	err := t.entries.ExecuteTemplate(w, template + TMPL_EXT, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO: remove
func md2html(mdText []byte) template.HTML {
	return template.HTML(github_flavored_markdown.Markdown(mdText))
}
