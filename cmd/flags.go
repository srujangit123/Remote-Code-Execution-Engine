package main

import "flag"

type CommandLineFlags struct {
	Platform            string
	MaxDockerContainers int
}

func ParseFlags() CommandLineFlags {
	var flags CommandLineFlags

	flag.StringVar(&flags.Platform, "platform", "amd64", "Platform can either be arm64 or amd64")
	flag.IntVar(&flags.MaxDockerContainers, "max-containers", 1, "Max number of docker containers running at any instant")

	flag.Parse()
	return flags
}
