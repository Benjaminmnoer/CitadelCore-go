package Session

import (
	"fmt"
	"net"
)

func Test() {
	fmt.Println("Session manager started")

	server, _ := net.Listen("tcp", "127.0.0.1:3724")

	for {
		connection, error := server.Accept()

		if error != nil {
			fmt.Printf("Error when accepting connection: %s\n", error)
		}

		StartSession(connection)
	}
}
