// Package visual holds the presentation layer of the blog (this is, templates)
package visual

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/mariomac/goblog/src/blog"
	"github.com/mariomac/goblog/src/fs"
	"github.com/mariomac/goblog/src/logr"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

var log = logr.Get()

// Templates wraps and extends the functionality of Go's template.Template
type Templates struct {
	*template.Template
}

var validTemplate = regexp.MustCompile(`.*\.html$`)

const templateExtension = ".html"

// Load gets all the pre-loaded templates from a given folder, populated with the entries
// returned by the getEntries function
func (t *Templates) Load(folder string, getEntries func() []blog.Render) {
	log.Printf("Scanning for templates in folder %s...\n", folder)

	templateFiles := fs.Search(folder, validTemplate)
	for _, f := range templateFiles {
		log.Printf("Template file found: %s\n", f)
	}

	t.Template = template.Must(
		template.New("golog_templates").Funcs(
			template.FuncMap{"entries": getEntries, "md2html": md2html}).ParseFiles(
			templateFiles...))
}

// Render renders the given template, with the given data, through the http.ResponseWriter
func (t *Templates) Render(w http.ResponseWriter, template string, data interface{}) {
	err := t.Template.ExecuteTemplate(w, template+templateExtension, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO: remove
func md2html(mdText []byte) template.HTML {
	// TODO: proper caching of goldmark
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.Highlighting,
		),
	)
	sb := strings.Builder{}
	if err := markdown.Convert(mdText, &sb); err != nil {
		// TODO: properly log/manage blogerr
		return template.HTML(`<h1>Error parsing markdown</h1><p>` + err.Error() + `</p>`)
	}
	return template.HTML(sb.String())
}
