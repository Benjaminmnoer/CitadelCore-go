package Session

import (
	Handlers "CitadelCore/AuthorisationServer/Handler"
	model "CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Repository"
	"CitadelCore/AuthorisationServer/SRP"
	"CitadelCore/Shared/Connection"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

var repository = Repository.InitializeAuthorisationRepository()

func HandleSession(conn net.Conn) {
	connection := Connection.CreateTcpConnection(conn, "500ms")
	endsession := false
	srp := SRP.NewSrp()

	fmt.Println("Session started")
	for !endsession {
		fmt.Println("Reading data")
		buffer := make([]byte, 256)
		size, error := connection.Read(&buffer)

		if error != nil {
			fmt.Printf("Error in reading message! %s\n", error)
			return
		}

		fmt.Printf("Size of message read %d\n", size)
		fmt.Printf("Command received %d\n", buffer[0])
		response, done := delegateCommand(buffer[0], buffer, srp)
		endsession = done
		output := new(bytes.Buffer)
		error = binary.Write(output, binary.BigEndian, response)

		if error != nil {
			fmt.Printf("Error in serializing response! %s\n", error)
			return
		}

		size, error = connection.Write(output.Bytes())

		if error != nil {
			fmt.Printf("Error in writing response! %s\n", error)
			return
		}

		fmt.Printf("Wrote %d bytes\n", size)

		if done {
			// TODO: Move logic elsewhere
			connection.Connection.Close()
			// session.done = true
		}
	}

	fmt.Println("Session finished")
}

func delegateCommand(cmd uint8, data []byte, srp *SRP.SRP6) (interface{}, bool) {
	// response := interface{}
	switch cmd {
	case model.AuthLogonChallenge:
		fmt.Println("AuthlogonChallenge registered")
		logonchallenge := model.LogonChallenge{}
		convertData(data, &logonchallenge)

		response := Handlers.HandleLogonChallenge(logonchallenge, repository, srp)

		return response, false // Expect more
	case model.AuthLogonProof:
		fmt.Println("AuthlogonProof registered")
		logonproof := model.LogonProof{}
		convertData(data, &logonproof)

		response, err := Handlers.HandleLogonProof(logonproof, srp)
		srp.PrintSRP()

		if err != nil {
			fmt.Printf("Error in handling logon proof: %e\n", err)
			return nil, true
		}

		return response, false // Expect realmlist command after proof.
	case model.AuthReconnectChallenge:
		fmt.Println("AuthReconnectChallenge registered")
		return nil, true
	case model.AuthReconnectProof:
		fmt.Println("AuthReconnectProof registered")
		return nil, true // Dont expect anymore after this. Perhaps realmlist?
	case model.RealmList:
		fmt.Println("Realmlist registered")
		return nil, true //Model.RealmListResponse
	}

	return nil, true
}

func convertData(data []byte, result interface{}) {
	reader := bytes.NewReader(data)
	error := binary.Read(reader, binary.LittleEndian, result)

	if error != nil {
		fmt.Printf("Error in binary conversion: %s\n", error)
		panic(error)
	}
}
