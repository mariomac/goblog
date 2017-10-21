package btemplate

import "testing"
import (
	"../bentry"
	"log"
)

const TEST_RESOURCES = "../../testresources/testset1"

func getEntries() []bentry.Entry {
	return make([]bentry.Entry, 0)
}

func TestTemplates_Load(t *testing.T) {
	templates := Templates{}
	templates.Load(TEST_RESOURCES, getEntries)

	expected := map[string]bool{
		"golog_templates": true,
		"thing2.html":     true, "thing3.html": true,
		"testsub2/thing3": true, "testsub/thing2": true,
		"test1.html": true, "test2.html": true,
	}

	actual := templates.entries.Templates()

	if len(expected) != len(actual) {
		t.Errorf("Failed loading templates. Expected: %d. Got: %d", len(expected), len(actual))
		for _, o := range actual {
			log.Println(o.Name())
		}
	}
}
