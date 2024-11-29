package codecontainer

import "time"

// TODO: Is it better to store this as environment variables?
const (
	GolangContainerImage = "alpine:latest"
	cppContainerImage    = "gcc:4.9"

	MAX_EXECUTION_TIME = 20 * time.Second
	MAX_SLEEP_TIME     = 400 * time.Second

	// The server running this will check for every 10 minutes whether there are zombie containers - completed containers and removes them.
	GarbageCollectionTimeWindow = 10 * time.Minute

	baseCodePath = "/home/sbharadwajtn/personal/Remote-Code-Execution-Engine/code/"
	// Make this part of env variables. There should be a setup scripts that adds these values to env variable.
	CppCodePath    = baseCodePath + "cpp"
	GolangCodePath = baseCodePath + "golang"

	// Path where the code files are mounted.
	TargetMountPath = "/container/code"

	BaseContainerCodeExecutablePath = TargetMountPath

	// How many times we should probe to see if there's an output file or not
	MAX_RETRIES = 10
)
