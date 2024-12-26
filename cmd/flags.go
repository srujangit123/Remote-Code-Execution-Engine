package main

import (
	"flag"
	"remote-code-engine/pkg/config"

	"go.uber.org/zap"
)

func ParseFlags() {
	flag.StringVar(&config.BaseCodePath, "code-dir", "/tmp/", "Base path to store the code files")
	flag.BoolVar(&config.ResourceConstraints, "resource-constraints", false, "Enable resource constraints (default false)")
	help := flag.Bool("help", false, "Display help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	logger.Info("parsed the flags",
		zap.String("code-dir", config.BaseCodePath),
		zap.Bool("resource-constraints", config.ResourceConstraints),
	)
}
