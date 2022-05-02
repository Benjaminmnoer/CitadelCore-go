package Session

import (
	"CitadelCore/AuthorisationServer/Model"
	"CitadelCore/AuthorisationServer/Session/Handlers"
	"CitadelCore/Shared/Communication"
	"CitadelCore/Shared/Helpers/Binary"
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

type AuthorisationSessionHandler struct {
	clientToSession sync.Map // Ip to session map.
}

func CreateSessionHandler() AuthorisationSessionHandler {
	return AuthorisationSessionHandler{clientToSession: sync.Map{}}
}

func (s AuthorisationSessionHandler) HandleSession(client Communication.Client, data []byte) ([]byte, Communication.SessionStatus, error) {
	session, found := s.clientToSession.Load(client.GetEndpoint())
	session = session.(AuthorisationSession)

	if !found {
		session = StartSession()
		s.clientToSession.Store(client.GetEndpoint(), session)
		fmt.Printf("Session started for endpoint %s", client.GetEndpoint())
	}

	response, status, err := s.delegateCommand(data[0], data, session.(AuthorisationSession))

	if err != nil {
		// TODO: Log instead of print to console
		return nil, Communication.EndConnection, fmt.Errorf("Error in auth session: %s\n", err)
	}

	bytes, err := Binary.Serialize(response)

	if err != nil {
		return nil, Communication.EndConnection, fmt.Errorf("Error in serialization: %s\n", err)
	}

	return bytes, status, nil
}

func (s AuthorisationSessionHandler) delegateCommand(cmd uint8, data []byte, session AuthorisationSession) ([]byte, Communication.SessionStatus, error) {
	switch cmd {
	case Model.AuthLogonChallenge:
		fmt.Println("AuthlogonChallenge registered")
		logonchallenge := Model.LogonChallenge{}
		convertData(data, &logonchallenge)
		session.AccountName = logonchallenge.GetAccountName()

		response := Handlers.HandleLogonChallenge(logonchallenge, session)
		result, err := convertToBytes(response)

		if err != nil {
			return nil, Communication.EndConnection, fmt.Errorf("Error in binary conversion: %s\n", err)
		}

		return result, Communication.KeepClient, nil // Expect logon proof
	case Model.AuthLogonProof:
		fmt.Println("AuthlogonProof registered")
		logonproof := Model.LogonProof{}
		convertData(data, &logonproof)

		response, err := Handlers.HandleLogonProof(logonproof, session)

		if err != nil {
			return nil, Communication.EndConnection, fmt.Errorf("error in handling logon proof: %e", err)
		}

		result, err := convertToBytes(response)

		if err != nil {
			return nil, Communication.EndConnection, fmt.Errorf("Error in binary conversion: %s\n", err)
		}

		return result, Communication.KeepClient, nil // Expect realmlist command after proof.
	case Model.AuthReconnectChallenge:
		fmt.Println("AuthReconnectChallenge registered")
		reconnectChallenge := Model.LogonChallenge{}
		convertData(data, &reconnectChallenge)

		response, err := Handlers.HandleReconnectChallenge(reconnectChallenge)
		session.ReconnectProof = response.Salt

		if err != nil {
			return nil, Communication.EndConnection, err
		}

		result, err := convertToBytes(response)

		if err != nil {
			return nil, Communication.EndConnection, fmt.Errorf("Error in binary conversion: %s\n", err)
		}

		return result, Communication.KeepClient, nil
	case Model.AuthReconnectProof:
		fmt.Println("AuthReconnectProof registered")
		reconnectProof := Model.ReconnectProof{}
		convertData(data, &reconnectProof)

		response, err := Handlers.HandleReconnectProof(reconnectProof, session)

		if err != nil {
			return nil, Communication.EndConnection, err
		}

		result, err := convertToBytes(response)

		if err != nil {
			return nil, Communication.EndConnection, fmt.Errorf("Error in binary conversion: %s\n", err)
		}

		return result, Communication.KeepClient, nil // Dont stop connection, it will ask for realmlist
	case Model.RealmList:
		fmt.Println("Realmlist registered")

		response, err := Handlers.HandleRealmList()

		if err != nil {
			fmt.Printf("Error getting realmlist: %e", err)
			return nil, Communication.KeepClient, fmt.Errorf("error getting realmlist: %e", err)
		}

		result, err := convertToBytes(response)

		if err != nil {
			return nil, Communication.EndConnection, fmt.Errorf("Error in binary conversion: %s\n", err)
		}

		return result, Communication.KeepClient, nil
	}

	return nil, Communication.EndConnection, fmt.Errorf("no matching command was found")
}

func convertData(data []byte, result interface{}) {
	reader := bytes.NewReader(data)
	error := binary.Read(reader, binary.LittleEndian, result)

	if error != nil {
		fmt.Printf("Error in binary conversion: %s\n", error)
		panic(error)
	}
}

// Binary conversion
func convertToBytes(input interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	error := binary.Write(buffer, binary.LittleEndian, input)

	if error != nil {
		return nil, fmt.Errorf("error converting to binary: %e", error)
	}

	return buffer.Bytes(), nil
}
