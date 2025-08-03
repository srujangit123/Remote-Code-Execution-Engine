package codecontainer

import "time"



	// The server running this will check for every 10 minutes whether there are zombie containers.
	GarbageCollectionTimeWindow = 5 * time.Minute

	// Path where the code files are mounted.
	TargetMountPath = "/container/code"
)
