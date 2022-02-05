package AuthorisationServer

import (
	"CitadelCore/Shared/Socket"
	"fmt"
)

func Start() {
	fmt.Println("Starting authserver...")
	Socket.Start(3724)
}
