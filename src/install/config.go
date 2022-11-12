package install

import (
	"io/ioutil"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"
)

// Config of the blog installation. Via file or env vars.
type Config struct {
	RootPath       string            `env:"GOBLOG_ROOT" yaml:"rootPath"`
	TLSPort        int               `env:"GOBLOG_HTTPS_PORT" yaml:"httpsPort"`
	InsecurePort   int               `env:"GOBLOG_HTTP_PORT" yaml:"httpPort"`
	Domain         string            `env:"GOBLOG_DOMAIN" yaml:"domain"`
	TLSCertPath    string            `env:"GOBLOG_TLS_CERT" yaml:"tlsCertPath"`
	TLSKeyPath     string            `env:"GOBLOG_TLS_KEY" yaml:"tlsKeyPath"`
	Redirect       map[string]string `env:"GOBLOG_REDIRECT" yaml:"redirect"`
	CacheSizeBytes int               `env:"GOBLOG_CACHE_SIZE_BYTES" yaml:"cacheSizeBytes"`
	EntriesPerPage int               `env:"GOBLOG_ENTRIES_PER_PAGE" yaml:"entriesPerPage"`
}

// ReadConfig gets a Config object from the environment and the provided yamlPath (optional)
func ReadConfig(yamlPath string) (Config, error) {
	// default values
	cfg := Config{
		RootPath:       "./",
		TLSPort:        8443,
		InsecurePort:   8080,
		TLSKeyPath:     "",
		TLSCertPath:    "",
		CacheSizeBytes: 32 * 1024 * 1024, // 32 MB
		EntriesPerPage: 5,
	}

	// override them with YAML
	if yamlPath != "" {
		yf, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			return cfg, err
		}
		if err := yaml.Unmarshal(yf, &cfg); err != nil {
			return cfg, err
		}
	}

	// override them with env values
	err := env.Parse(&cfg, env.Options{})
	return cfg, err
}
