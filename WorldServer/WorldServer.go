package WorldServer

import (
	"CitadelCore/WorldServer/Session"
	"fmt"
)

func Start() {
	fmt.Println("Starting worldserver...")

	// Load config

	// Start connections in other threads.
	go Session.StartSession()

	// Game loop should occupy main thread
	GameLoop()
}
