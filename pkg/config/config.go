package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Architecture string
type Language string

// Supported languages
const (
	Golang Language = "golang"
	Cpp    Language = "cpp"
)

// Supported architectures
const (
	Arm64  Architecture = "arm64"
	X86_64 Architecture = "x86_64"
)

type LanguageConfig struct {
	Extension string `yaml:"extension"`
	Image     string `yaml:"image"`
}

type ImageConfig struct {
	Arm64  map[Language]LanguageConfig `yaml:"arm64"`
	X86_64 map[Language]LanguageConfig `yaml:"x86_64"`
}

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
