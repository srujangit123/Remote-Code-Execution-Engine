package main

import "remote-code-engine/pkg/config"

type Request struct {
	EncodedCode  string          `json:"code"`
	EncodedInput string          `json:"input"`
	Language     config.Language `json:"language"`
}

type Response struct {
	Output string `json:"output"`
}
