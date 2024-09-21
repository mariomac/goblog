package install

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigOverride(t *testing.T) {
	tmp, err := os.CreateTemp("", "testconfig")
	require.NoError(t, err)
	_, err = tmp.WriteString(`domain: blog.com
httpsPort: 443
httpPort: 80
maxRequests:
  number: 30
  period: 10s
`)
	require.NoError(t, err)
	require.NoError(t, tmp.Close())
	require.NoError(t, os.Setenv("GOBLOG_HTTP_PORT", "81"))

	cfg, err := ReadConfig(tmp.Name())
	require.NoError(t, err)
	// verify that YAML overrides defaults and Env overrides YAML
	assert.Equal(t, "./", cfg.RootPath)
	assert.Equal(t, "blog.com", cfg.Domain)
	assert.Equal(t, 443, cfg.TLSPort)
	assert.Equal(t, 81, cfg.InsecurePort)
	assert.Equal(t, MaxRequestsCfg{
		Number: 30,
		Period: 10 * time.Second,
	}, cfg.MaxRequests)
}

func TestConfigOverride_Env(t *testing.T) {
	require.NoError(t, os.Setenv("GOBLOG_HTTP_PORT", "81"))
	require.NoError(t, os.Setenv("GOBLOG_MAX_REQUESTS_NUMBER", "30"))
	require.NoError(t, os.Setenv("GOBLOG_MAX_REQUESTS_PERIOD", "1m"))

	cfg, err := ReadConfig("")
	require.NoError(t, err)
	// verify that YAML overrides defaults and Env overrides YAML
	assert.Equal(t, 81, cfg.InsecurePort)
	assert.Equal(t, MaxRequestsCfg{
		Number: 30,
		Period: time.Minute,
	}, cfg.MaxRequests)
}
