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

type ImageConfig map[Language]LanguageConfig

func LoadConfig(configPath string) (*ImageConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the config file: %w", err)
	}

	var config ImageConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the config file: %w", err)
	}

	return &config, nil
}

func (c *ImageConfig) GetLanguageConfig(lang Language) LanguageConfig {
	return (*c)[lang]
}

func (c *ImageConfig) IsLanguageSupported(lang Language) bool {
	_, ok := (*c)[lang]
	return ok
}

func (c *ImageConfig) GetSupportedLanguages() []Language {
	var languages []Language
	for lang := range *c {
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
