package install

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigOverride(t *testing.T) {
	tmp, err := os.CreateTemp("", "testconfig")
	require.NoError(t, err)
	_, err = tmp.WriteString(`domain: blog.com
httpsPort: 443
httpPort: 80
`)
	require.NoError(t, err)
	require.NoError(t, tmp.Close())
	require.NoError(t, os.Setenv("GOBLOG_HTTP_PORT", "81"))

	cfg, err := ReadConfig(tmp.Name())
	require.NoError(t, err)
	// verify that YAML overrides defaults and Env overrides YAML
	assert.Equal(t, "./sample", cfg.RootPath)
	assert.Equal(t, "blog.com", cfg.Domain)
	assert.Equal(t, 443, cfg.TLSPort)
	assert.Equal(t, 81, cfg.InsecurePort)
}
