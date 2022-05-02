package AuthorisationServer

import (
	"CitadelCore/AuthorisationServer/Session"
	"CitadelCore/Shared/Communication"
	"fmt"
	"net"
)

func Start() {
	fmt.Println("Starting authserver...")
	Session.StartServer()
}

// Accepts connection and adds it to the client handler.
func (s AuthorisationSessionHandler) StartServer() {
	fmt.Println("Session manager started")
	srp.InitializaSRP()
	s.clientHandler = Communication.CreateClientHandler(s, 4)

	server, _ := net.Listen("tcp", "127.0.0.1:3724")

	for {
		connection, error := server.Accept()

		if error != nil {
			fmt.Printf("Error when accepting connection: %s\n", error)
		}

		s.clientHandler.AddClient(Communication.Client{Connection: connection})
	}
}
