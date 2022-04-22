package WorldServer

import (
	"time"
)

var lastUpdateTime time.Time

func GameLoop() {
	lastUpdateTime = time.Now()

	for {
		now := time.Now()
		// timeSinceLastUpdate := now.Sub(lastUpdateTime)
		// fmt.Printf("Time since update %d ns\n", timeSinceLastUpdate.Nanoseconds())

		lastUpdateTime = now
	}
}
