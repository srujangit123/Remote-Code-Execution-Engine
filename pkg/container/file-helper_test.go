package codecontainer

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"remote-code-engine/pkg/config"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestCreateFile(t *testing.T) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	))

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")
	content := "Hello, World!"
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))

	fileName, err := createFile(filePath, encodedContent, logger)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	if fileName != "testfile.txt" {
		t.Errorf("expected file name 'testfile.txt', got '%s'", fileName)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(data) != content {
		t.Errorf("expected file content '%s', got '%s'", content, string(data))
	}
}

func TestCreateFileInvalidBase64(t *testing.T) {
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	))

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile.txt")
	invalidEncodedContent := "invalid_base64_content"

	_, err := createFile(filePath, invalidEncodedContent, logger)
	if err == nil {
		t.Fatal("expected an error due to invalid base64 content, but got none")
	}
}
func TestGetCodeAndInputFilePathsHost(t *testing.T) {
	code := &Code{
		Language: "golang",
		LanguageConfig: config.LanguageConfig{
			Extension: ".go",
		},
	}

	codeFilePath, inputFilePath := getCodeAndInputFilePathsHost(code)

	expectedCodeDir := config.GetHostLanguageCodePath(code.Language)
	expectedCodeFilePath := filepath.Join(expectedCodeDir, filepath.Base(codeFilePath))
	expectedInputFilePath := filepath.Join(expectedCodeDir, filepath.Base(inputFilePath))

	if codeFilePath != expectedCodeFilePath {
		t.Errorf("expected code file path '%s', got '%s'", expectedCodeFilePath, codeFilePath)
	}

	if inputFilePath != expectedInputFilePath {
		t.Errorf("expected input file path '%s', got '%s'", expectedInputFilePath, inputFilePath)
	}
}

func TestGetFilePathHost(t *testing.T) {
	hostCodeDirectoryPath := "/tmp/code"
	fileName := "testfile.go"
	expectedFilePath := "/tmp/code/testfile.go"

	result := getFilePathHost(hostCodeDirectoryPath, fileName)

	if result != expectedFilePath {
		t.Errorf("expected file path '%s', got '%s'", expectedFilePath, result)
	}
}

func TestGetFilePathContainer(t *testing.T) {
	mountPath := "/mnt/container"
	fileName := "testfile.txt"
	expectedFilePath := "/mnt/container/testfile.txt"

	result := getFilePathContainer(mountPath, fileName)

	if result != expectedFilePath {
		t.Errorf("expected file path '%s', got '%s'", expectedFilePath, result)
	}
}
