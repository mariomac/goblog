package fs

import (
	"fmt"
	"regexp"
	"testing"
)

func equals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

const testResources = "../../testresources/testset1"

func TestEmptySearchNoMatches(t *testing.T) {
	found := Search(testResources, regexp.MustCompile("^inexistingfile$"))
	if len(found) != 0 {
		t.Errorf("Search should be empty, but got: %s", fmt.Sprint(found))
	}
}

func TestEmptySearchWrongFolder(t *testing.T) {
	found := Search("wrongfolderhere", nil)
	if len(found) != 0 {
		t.Errorf("Search should be empty, but got: %s", fmt.Sprint(found))
	}
}

func TestAllSearch(t *testing.T) {
	found := Search(testResources, nil)
	expected := []string{
		"../../testresources/testset1/entry.html",
		"../../testresources/testset1/index.html",
		"../../testresources/testset1/testsub/thing.md",
		"../../testresources/testset1/testsub/thing2.html",
		"../../testresources/testset1/testsub2/thing3.html",
		"../../testresources/testset1/testsub2/thing4.md",
		"../../testresources/testset1/zztest.md",
	}
	if !equals(found, expected) {
		t.Errorf("Search does not contain expected files. Expecting:%s\nGot:%s",
			fmt.Sprint(expected), fmt.Sprint(found))
	}
}

func TestMarkdownSearch(t *testing.T) {
	found := Search(testResources, regexp.MustCompile(`\.md$`))
	expected := []string{
		"../../testresources/testset1/testsub/thing.md",
		"../../testresources/testset1/testsub2/thing4.md",
		"../../testresources/testset1/zztest.md",
	}
	if !equals(found, expected) {
		t.Errorf("Search does not contain expected files. Expecting:%s\nGot:%s",
			fmt.Sprint(expected), fmt.Sprint(found))
	}
}
