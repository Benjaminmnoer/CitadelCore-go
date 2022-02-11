package Session

import (
	Handlers "CitadelCore/AuthorisationServer/Handler"
	model "CitadelCore/AuthorisationServer/Model"
	"CitadelCore/Shared/Connection"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type session struct {
	build      string
	connection Connection.TcpConnection
}

func StartSession(connection net.Conn) {
	session := session{
		connection: Connection.CreateTcpConnection(connection, "500ms"),
	}

	handleSession(session)
}

func handleSession(session session) {
	buffer := make([]byte, 256)
	size, error := session.connection.Read(&buffer)

	if error != nil {
		fmt.Println("Error in reading message!")
		return
	}

	fmt.Printf("Size of message received %d\n", size)

	response := delegateCommand(buffer[0], buffer)
	output := new(bytes.Buffer)
	error = binary.Write(output, binary.LittleEndian, response)

	if error != nil {
		fmt.Println("Error in reading message!")
		return
	}

	session.connection.Write(output.Bytes())
}

func delegateCommand(cmd uint8, data []byte) interface{} {
	switch cmd {
	case model.AuthLogonChallenge:
		fmt.Println("AuthlogonChallenge registered")
		logonchallenge := convertData(data, model.LogonChallenge{}).(model.LogonChallenge)

		fmt.Printf("Gamename: %s\n", logonchallenge.Gamename)
		fmt.Printf("Accountname: %s\n", logonchallenge.Accountname)

		response, error := Handlers.HandleLogonChallenge(logonchallenge)

		if error != nil {
			fmt.Printf("Error handling logon challenge: %s\n", error)
		}

		return response
	case model.AuthLogonProof:
		fmt.Println("AuthlogonProof registered")
		return nil
	case model.AuthReconnectChallenge:
		fmt.Println("AuthReconnectChallenge registered")
		return nil
	case model.AuthReconnectProof:
		fmt.Println("AuthReconnectProof registered")
		return nil
	case model.RealmList:
		fmt.Println("Realmlist registered")
		return nil
	}

	return nil
}

func convertData(data []byte, result interface{}) interface{} {
	reader := bytes.NewReader(data)
	error := binary.Read(reader, binary.LittleEndian, &result)

	if error != nil {
		fmt.Printf("Error in binary conversion: %s\n", error)
		return nil
	}

	return result
}
