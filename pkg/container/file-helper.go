package codecontainer

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"remote-code-engine/pkg/config"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func getFilePathContainer(mountPath, fileName string) string {
	return filepath.Join(mountPath, fileName)
}

func getFilePathHost(hostCodeDirectoryPath, fileName string) string {
	return filepath.Join(hostCodeDirectoryPath, fileName)
}

func getCodeAndInputFilePathsHost(code *Code) (string, string) {
	filename := uuid.New().String()
	codeDirectoryPathHost := config.GetHostLanguageCodePath(code.Language)
	codeFileName := filename + code.Extension
	inputFileName := filename + ".txt"

	return getFilePathHost(codeDirectoryPathHost, codeFileName),
		getFilePathHost(codeDirectoryPathHost, inputFileName)
}

func createFile(filePath, base64FileContent string, logger *zap.Logger) (string, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create the file: %w", err)
	}

	data, err := base64.StdEncoding.DecodeString(base64FileContent)
	if err != nil {
		return filepath.Base(filePath), fmt.Errorf("failed to decode the file content: %w", err)
	}

	n, err := f.Write(data)
	if err != nil {
		return filepath.Base(filePath), fmt.Errorf("failed to write the content to the file: %w", err)
	}
	logger.Info("wrote the file content to the file",
		zap.String("file path", filePath),
		zap.Int("bytes", n),
	)

	return filepath.Base(filePath), nil
}

func createCodeAndInputFilesHost(code *Code, logger *zap.Logger) (string, string, error) {
	codeFilePath, inputFilePath := getCodeAndInputFilePathsHost(code)
	codeFileName, err := createFile(codeFilePath, code.EncodedCode, logger)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the code file: %w", err)
	}

	inputFileName, err := createFile(inputFilePath, code.EncodedInput, logger)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the input file: %w", err)
	}

	return codeFileName, inputFileName, nil
}
