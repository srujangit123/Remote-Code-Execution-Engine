package codecontainer

import (
	"remote-code-engine/pkg/config"
	"testing"
)

func TestGetContainerCommand(t *testing.T) {
	tests := []struct {
		name            string
		code            *Code
		codeFileName    string
		inputFileName   string
		expectedCommand []string
	}{
		{
			name: "golang command",
			code: &Code{
				Language: "golang",
				LanguageConfig: config.LanguageConfig{
					Extension: ".go",
					Command:   "go run {{FILE}} < {{INPUT}}",
				},
			},
			codeFileName:  "main.go",
			inputFileName: "input.txt",
			expectedCommand: []string{
				"sh", "-c",
				"go run /container/code/main.go < /container/code/input.txt",
			},
		},
		{
			name: "cpp command",
			code: &Code{
				Language: "cpp",
				LanguageConfig: config.LanguageConfig{
					Extension: ".cpp",
					Command:   "g++ {{FILE}} -o a.out && a.out < {{INPUT}}",
				},
			},
			codeFileName:  "main.cpp",
			inputFileName: "input.txt",
			expectedCommand: []string{
				"sh", "-c",
				"g++ /container/code/main.cpp -o a.out && a.out < /container/code/input.txt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			command := getContainerCommand(tt.code, tt.codeFileName, tt.inputFileName)
			if len(command) != len(tt.expectedCommand) {
				t.Errorf("expected command length %d, got %d", len(tt.expectedCommand), len(command))
			}
			for i := range command {
				if command[i] != tt.expectedCommand[i] {
					t.Errorf("expected command '%s', got '%s'", tt.expectedCommand[i], command[i])
				}
			}
		})
	}
}
