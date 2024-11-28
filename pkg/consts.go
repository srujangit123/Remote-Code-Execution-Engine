package docker

import "time"

// TODO: Is it better to store this as environment variables?
const (
	GolangContainerImage = "alpine:latest"
	cppContainerImage    = "gcc:4.9"

	MAX_EXECUTION_TIME = 20 * time.Second
	MAX_SLEEP_TIME     = 400 * time.Second

	// The server running this will check for every 10 minutes whether there are zombie containers - completed containers and removes them.
	GarbageCollectionTimeWindow = 10 * time.Minute

	CppCodePath = "/home/srujan/Documents/code/cpp"

	GolangCodePath = "/home/srujan/Documents/code/golang"
)
