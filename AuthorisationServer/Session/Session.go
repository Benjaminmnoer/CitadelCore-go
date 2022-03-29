package Session

import (
	Handlers "CitadelCore/AuthorisationServer/Handler"
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Repository"
	"CitadelCore/AuthorisationServer/SRP"
	"CitadelCore/Shared/Binary"
	"CitadelCore/Shared/Connection"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

var repository = Repository.InitializeAuthorisationRepository()
var accountname = ""
var reconnectproof [16]byte

func HandleSession(conn net.Conn) {
	connection := Connection.CreateTcpConnection(conn, "500ms")
	endsession := false
	srp := SRP.NewSrp()

	fmt.Println("Session started")
	for !endsession {
		fmt.Println("Reading data")
		buffer := make([]byte, 256)
		_, err := connection.Read(&buffer)

		if err != nil {
			fmt.Printf("Error in reading message! %s\n", err)
			break
		}

		response, err, done := delegateCommand(buffer[0], buffer, srp)

		if err != nil {
			// TODO: Log instead of print to console
			fmt.Printf("Error in auth session.\n%s\n", err)
			break
		}

		endsession = done
		bytes, err := Binary.Serialize(response)

		if err != nil {
			fmt.Printf("Error in serializing response! %s\n", err)
			return
		}

		_, err = connection.Write(bytes)

		if err != nil {
			fmt.Printf("Error in writing response! %s\n", err)
			return
		}

		if done {
			// TODO: Move logic elsewhere
			connection.Connection.Close()
		}
	}

	fmt.Println("Session finished")
}

func delegateCommand(cmd uint8, data []byte, srp *SRP.SRP6) (interface{}, error, bool) {
	// response := interface{}
	switch cmd {
	case Model.AuthLogonChallenge:
		fmt.Println("AuthlogonChallenge registered")
		logonchallenge := Model.LogonChallenge{}
		convertData(data, &logonchallenge)
		accountname = logonchallenge.GetAccountName()

		response := Handlers.HandleLogonChallenge(logonchallenge, repository, srp)

		return response, nil, false // Expect logon proof
	case Model.AuthLogonProof:
		fmt.Println("AuthlogonProof registered")
		logonproof := Model.LogonProof{}
		convertData(data, &logonproof)

		response, err := Handlers.HandleLogonProof(logonproof, srp)

		if err != nil {
			return nil, fmt.Errorf("Error in handling logon proof: %e\n", err), true
		}

		return response, nil, false // Expect realmlist command after proof.
	case Model.AuthReconnectChallenge:
		fmt.Println("AuthReconnectChallenge registered")
		reconnectChallenge := Model.LogonChallenge{}
		convertData(data, &reconnectChallenge)

		response, err := Handlers.HandleReconnectChallenge(reconnectChallenge)
		reconnectproof = response.Salt

		if err != nil {
			return nil, err, true
		}

		return response, nil, false
	case Model.AuthReconnectProof:
		fmt.Println("AuthReconnectProof registered")
		reconnectProof := Model.ReconnectProof{}
		convertData(data, &reconnectProof)

		response, err := Handlers.HandleReconnectProof(reconnectProof, accountname, reconnectproof)

		if err != nil {
			return nil, err, true
		}

		return response, nil, true // Do we expect realm list after this? We will see i guess.
	case Model.RealmList:
		fmt.Println("Realmlist registered")

		realmlist, err := Handlers.HandleRealmList(repository)

		if err != nil {
			fmt.Printf("Error getting realmlist: %e", err)
			return nil, err, true
		}

		fmt.Println(realmlist)

		return realmlist, nil, true //Model.RealmListResponse
	}

	return nil, fmt.Errorf("No matching command was found"), true
}

// TODO: Move to more reasonable place.
func convertData(data []byte, result interface{}) {
	reader := bytes.NewReader(data)
	error := binary.Read(reader, binary.LittleEndian, result)

	if error != nil {
		fmt.Printf("Error in binary conversion: %s\n", error)
		panic(error)
	}
}
