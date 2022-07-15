package fs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
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
	found, err := Search(testResources, regexp.MustCompile("^inexistingfile$"))
	require.NoError(t, err)
	if len(found) != 0 {
		t.Errorf("Search should be empty, but got: %s", fmt.Sprint(found))
	}
}

func TestEmptySearchWrongFolder(t *testing.T) {
	found, err := Search("wrongfolderhere", nil)
	require.NoError(t, err)
	if len(found) != 0 {
		t.Errorf("Search should be empty, but got: %s", fmt.Sprint(found))
	}
}

func TestAllSearch(t *testing.T) {
	found, err := Search(testResources, nil)
	require.NoError(t, err)
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
	found, err := Search(testResources, regexp.MustCompile(`\.md$`))
	require.NoError(t, err)
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
