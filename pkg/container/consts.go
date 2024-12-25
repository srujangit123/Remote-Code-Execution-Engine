package codecontainer

import "time"

const (
	MAX_EXECUTION_TIME = 60 * time.Second

	// The server running this will check for every 10 minutes whether there are zombie containers - completed containers and removes them.
	GarbageCollectionTimeWindow = 5 * time.Minute

	// Path where the code files are mounted.
	TargetMountPath = "/container/code"
)
