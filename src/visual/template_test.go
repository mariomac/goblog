package visual

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testResources = "../../testresources/testset1"

func TestTemplates_Load(t *testing.T) {
	templater, err := LoadTemplates(testResources)
	require.NoError(t, err)

	actual := map[string]struct{}{}
	for _, template := range templater.templates.Templates() {
		actual[template.Name()] = struct{}{}
	}

	assert.Equal(t, map[string]struct{}{
		"golog_templates": {},
		"thing2.html":     {},
		"thing3.html":     {},
		"testsub2/thing3": {},
		"testsub/thing2":  {},
		"entry.html":      {},
		"index.html":      {},
	}, actual)

}
