package main

import (
	"flag"
	"remote-code-engine/pkg/config"

	"go.uber.org/zap"
)

func ParseFlags() (string) {
	flag.StringVar(&config.BaseCodePath, "code-dir", "/tmp/", "Base path to store the code files")
	flag.BoolVar(&config.ResourceConstraints, "resource-constraints", false, "Enable resource constraints (default false)")
	configFilePath := flag.String("config-file", "config.yml", "Path to the config file")
	help := flag.Bool("help", false, "Display help")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	logger.Info("parsed the flags",
		zap.String("code-dir", config.BaseCodePath),
		zap.Bool("resource-constraints", config.ResourceConstraints),
		zap.String("config-file", *configFilePath),
	)

	return *configFilePath
}

	if *help {
		flag.Usage()
		return
	}

	logger.Info("parsed the flags",
		zap.String("code-dir", config.BaseCodePath),
		zap.Bool("resource-constraints", config.ResourceConstraints),
	)
}
