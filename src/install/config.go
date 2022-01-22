package install

import (
	"io/ioutil"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v2"
)

// Config of the blog installation. Via file or env vars.
type Config struct {
	RootPath     string `env:"GOBLOG_ROOT" yaml:"rootPath"`
	TLSPort      int    `env:"GOBLOG_HTTPS_PORT" yaml:"httpsPort"`
	InsecurePort int    `env:"GOBLOG_HTTP_PORT" yaml:"httpPort"`
	Domain       string `env:"GOBLOG_DOMAIN" yaml:"domain"`
	TLSCertPath  string `env:"GOBLOG_TLS_CERT" yaml:"tlsCertPath"`
	TLSKeyPath   string `env:"GOBLOG_TLS_KEY" yaml:"tlsKeyPath"`
}

// ReadConfig gets a Config object from the environment and the provided yamlPath (optional)
func ReadConfig(yamlPath string) (Config, error) {
	// default values
	cfg := Config{
		RootPath:     "./sample",
		TLSPort:      8443,
		InsecurePort: 8080,
		TLSKeyPath:   "",
		TLSCertPath:  "",
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
