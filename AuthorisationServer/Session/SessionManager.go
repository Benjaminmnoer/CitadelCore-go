package Session

import (
	srp "CitadelCore/AuthorisationServer/SRP"
	"fmt"
	"net"
)

func StartServer() {
	fmt.Println("Session manager started")
	srp.InitializaSRP()

	server, _ := net.Listen("tcp", "127.0.0.1:3724")

	for {
		connection, error := server.Accept()

		if error != nil {
			fmt.Printf("Error when accepting connection: %s\n", error)
		}

		HandleSession(connection)
	}
}
