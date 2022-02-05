package WorldServer

import (
	"CitadelCore/Shared/Socket"
	"fmt"
)

func Start() {
	fmt.Println("Starting authserver...")
	Socket.Start(8828)
}
