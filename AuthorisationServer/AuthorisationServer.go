package AuthorisationServer

import (
	"CitadelCore/AuthorisationServer/Session"
	"fmt"
)

func Start() {
	fmt.Println("Starting authserver...")
	Session.Test()
}
