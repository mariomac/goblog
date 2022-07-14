// Package visual holds the presentation layer of the blog (this is, template)
package visual

import (
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strings"

	"github.com/mariomac/goblog/src/fs"
	"github.com/mariomac/goblog/src/logr"
	"github.com/sirupsen/logrus"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

var log = logr.Get()

type TemplateType string

const (
	EntryTemplate = "entry.html"
	IndexTemplate = "index.html"
)

var mandatoryTemplates = []TemplateType{EntryTemplate, IndexTemplate}

// Templater wraps and extends the functionality of Go's template
type Templater struct {
	templates *template.Template
}

var validTemplate = regexp.MustCompile(`\.html$`)

func LoadTemplates(
	folder string,
) (Templater, error) {
	tlog := log.WithField("folder", folder)
	tlog.Info("Scanning for template")

	templateFiles, err := fs.Search(folder, validTemplate)
	if err != nil {
		return Templater{}, fmt.Errorf("scanning for template in folder %s: %w", folder, err)
	}
	if tlog.Level <= logrus.DebugLevel {
		for _, f := range templateFiles {
			tlog.WithField("file", f).Debug("Template file found")
		}
	}
	templates, err := template.New("golog_templates").
		Funcs(template.FuncMap{
			"md2html": md2html(),
		}).
		ParseFiles(templateFiles...)
	if err != nil {
		return Templater{}, fmt.Errorf("parsing template files: %w", err)
	}

	for _, mt := range mandatoryTemplates {
		if templates.Lookup(string(mt)) == nil {
			return Templater{}, fmt.Errorf("missing mandatory template: %s", mt)
		}
	}

	return Templater{templates: templates}, nil
}

func (t *Templater) Render(template TemplateType, data interface{}, dest io.Writer) error {
	return t.templates.ExecuteTemplate(dest, string(template), data)
}

// TODO: remove
func md2html() func(mdText []byte) template.HTML {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.Highlighting,
		),
	)
	return func(mdText []byte) template.HTML {
		sb := strings.Builder{}
		if err := markdown.Convert(mdText, &sb); err != nil {
			// TODO: properly log/manage blogerr
			return template.HTML(`<h1>Error rendering content</h1><p>` + err.Error() + `</p>`)
		}
		return template.HTML(sb.String())
	}
}
