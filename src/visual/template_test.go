package visual

import (
	"testing"

	"github.com/mariomac/goblog/src/logr"
)
import (
	"github.com/mariomac/goblog/src/blog"
)

var tlog = logr.Get()

const testResources = "../../testresources/testset1"

func getEntries() []blog.Entry {
	return make([]blog.Entry, 0)
}

func TestTemplates_Load(t *testing.T) {
	templates := Templates{}
	templates.Load(testResources, getEntries)

	expected := map[string]bool{
		"golog_templates": true,
		"thing2.html":     true, "thing3.html": true,
		"testsub2/thing3": true, "testsub/thing2": true,
		"test1.html": true, "test2.html": true,
	}

	actual := templates.Templates()

	if len(expected) != len(actual) {
		t.Errorf("Failed loading templates. Expected: %d. Got: %d", len(expected), len(actual))
		for _, o := range actual {
			tlog.Println(o.Name())
		}
	}
}
