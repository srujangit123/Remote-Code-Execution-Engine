package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Language string

// Supported languages
const (
	Golang Language = "golang"
	Cpp    Language = "cpp"
)

var (
	BaseCodePath        string
	ResourceConstraints bool
)

type LanguageConfig struct {
	Extension string `yaml:"extension"`
	Image     string `yaml:"image"`
	Command   string `yaml:"command"`
}

type ResourceConfig struct {
	Memory    int64  `yaml:"memory"`
	NanoCPUs  int64  `yaml:"nanocpus"`
	Nproc     int64  `yaml:"nproc"`
	Nofile    int64  `yaml:"nofile"`
	Fsize     int64  `yaml:"fsize"`
	Core      int64  `yaml:"core"`
	Timeout   int64  `yaml:"timeout"`
}

type Config struct {
	Languages map[Language]LanguageConfig `yaml:",inline"`
	Resources ResourceConfig             `yaml:"resources"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the config file: %w", err)
	}

	return &config, nil
}

func (c *Config) GetLanguageConfig(lang Language) LanguageConfig {
	langConfig := c.Languages[lang]
	envImage := os.Getenv(fmt.Sprintf("%s_IMAGE", lang))
	if envImage != "" {
		langConfig.Image = envImage
	}
	return langConfig
}

func (c *Config) IsLanguageSupported(lang Language) bool {
	_, ok := c.Languages[lang]
	return ok
}

func (c *Config) GetSupportedLanguages() []Language {
	languages := make([]Language, 0, len(c.Languages))
	for lang := range c.Languages {
		languages = append(languages, lang)
	}
	return languages
}

func GetHostLanguageCodePath(lang Language) string {
	return filepath.Join(BaseCodePath, string(lang))
}

func IsResourceConstraintsEnabled() bool {
	return ResourceConstraints
}
