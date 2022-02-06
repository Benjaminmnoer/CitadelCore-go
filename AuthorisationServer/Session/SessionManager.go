package Session

import (
	model "CitadelCore/AuthorisationServer/Model"
	"encoding/gob"
	"fmt"
	"net"
)

func Test() {
	fmt.Println("Session manager started")

	server, _ := net.Listen("tcp", "127.0.0.1:3724")

	// var buffer []byte
	for {
		connection, _ := server.Accept()
		buffer := make([]byte, 1024)
		size, error := connection.Read(buffer)

		if error != nil {
			fmt.Println("Error in reading message!")
			continue
		}

		fmt.Printf("Size of message received %d\n", size)

		logonchallenge := model.AuthorisationLogonChallenge{}
		decoder := gob.NewDecoder(buffer)
		decoder.Decode(logonchallenge)

		if error != nil {
			fmt.Println("Error in unmarshalling")
			continue
		}

		fmt.Printf("Account name%s\n", string(logonchallenge.Cmd))
		connection.Close()
	}
}
