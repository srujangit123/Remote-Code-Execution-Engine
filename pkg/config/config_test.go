package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Define a sample config
	sampleConfig := ImageConfig{
		Golang: {
			Extension: ".go",
			Image:     "golang:latest",
			Command:   "go build",
		},
		Cpp: {
			Extension: ".cpp",
			Image:     "cpp:latest",
			Command:   "g++ -o main",
		},
	}

	// Marshal the sample config to YAML
	data, err := yaml.Marshal(&sampleConfig)
	if err != nil {
		t.Fatalf("failed to marshal sample config: %v", err)
	}

	// Write the sample config to a temporary file
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("failed to write sample config to file: %v", err)
	}

	// Load the config using the LoadConfig function
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Compare the loaded config with the sample config
	if len(*loadedConfig) != len(sampleConfig) {
		t.Fatalf("expected %d languages, got %d", len(sampleConfig), len(*loadedConfig))
	}

	for lang, expectedConfig := range sampleConfig {
		loadedLangConfig, ok := (*loadedConfig)[lang]
		if !ok {
			t.Fatalf("language %s not found in loaded config", lang)
		}

		if loadedLangConfig != expectedConfig {
			t.Errorf("expected config for language %s: %+v, got: %+v", lang, expectedConfig, loadedLangConfig)
		}
	}
}

func TestGetHostLanguageCodePath(t *testing.T) {
	BaseCodePath = "/base/path"

	tests := []struct {
		lang     Language
		expected string
	}{
		{Golang, "/base/path/golang"},
		{Cpp, "/base/path/cpp"},
	}

	for _, tt := range tests {
		t.Run(string(tt.lang), func(t *testing.T) {
			result := GetHostLanguageCodePath(tt.lang)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
