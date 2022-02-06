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
		decoder := gob.NewDecoder(connection)
		// buffer := new(bytes.Buffer)
		// size, error := connection.Read(buffer.Bytes())

		// if error != nil {
		// 	fmt.Println("Error in reading message!")
		// 	continue
		// }

		// fmt.Printf("Size of message received %d\n", size)

		logonchallenge := model.AuthorisationLogonChallenge{}
		// reader := bytes.NewReader(buffer.Bytes())
		// error = binary.Read(reader, binary.LittleEndian, &logonchallenge)
		error := decoder.Decode(&logonchallenge)

		if error != nil {
			fmt.Printf("Error in byte conversion. Error code %s", error.Error())
			continue
		}

		fmt.Printf("Account name%s\n", string(logonchallenge.Cmd))
		connection.Close()
	}
}
