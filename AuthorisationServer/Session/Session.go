package Session

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Repository"
	"CitadelCore/AuthorisationServer/SRP"
	Handlers "CitadelCore/AuthorisationServer/Session/Handler"
	"CitadelCore/Shared/Communication"
	"CitadelCore/Shared/Connection"
	"CitadelCore/Shared/Helpers/Binary"
	"bytes"
	"encoding/binary"
	"fmt"
)

var accountname = ""
var reconnectproof [16]byte

type Session struct {
	repository     Repository.AuthorisationRepository
	accountName    string
	reconnectProof [16]byte
}

func StartSession(repo Repository.AuthorisationRepository) Session {
	return Session{repository: repo}
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

func (s Session) HandleSession(client Communication.Client) {
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

func (s Session) DelegateCommand(cmd uint8, data []byte, srp *SRP.SRP6) (interface{}, error, bool) {
	switch cmd {
	case Model.AuthLogonChallenge:
		fmt.Println("AuthlogonChallenge registered")
		logonchallenge := Model.LogonChallenge{}
		convertData(data, &logonchallenge)
		s.accountName = logonchallenge.GetAccountName()

		response := Handlers.HandleLogonChallenge(logonchallenge, s.repository, srp)

		return response, nil, false // Expect logon proof
	case Model.AuthLogonProof:
		fmt.Println("AuthlogonProof registered")
		logonproof := Model.LogonProof{}
		convertData(data, &logonproof)

		response, err := Handlers.HandleLogonProof(logonproof, srp)

		if err != nil {
			return nil, fmt.Errorf("error in handling logon proof: %e", err), true
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

		response, err := Handlers.HandleReconnectProof(reconnectProof, accountname, s.reconnectProof)

		if err != nil {
			return nil, err, true
		}

		return response, nil, true // Do we expect realm list after this? We will see i guess.
	case Model.RealmList:
		fmt.Println("Realmlist registered")

		realmlist, err := Handlers.HandleRealmList(s.repository)

		if err != nil {
			fmt.Printf("Error getting realmlist: %e", err)
			return nil, err, true
		}

		fmt.Println(realmlist)

		return realmlist, nil, true //Model.RealmListResponse
	}

	return nil, fmt.Errorf("no matching command was found"), true
}
