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

type Authsession struct {
	build      string
	connection Connection.TcpConnection
	srp        *SRP.SRP6
	done       bool
}

func StartSession(connection net.Conn) {
	session := Authsession{
		// TODO: Reduce deadline?
		connection: Connection.CreateTcpConnection(connection, "500ms"),
		done:       false,
	}

	handleSession(&session)
}

func handleSession(session *Authsession) {
	fmt.Println("Session started")
	for !session.done {
		buffer := make([]byte, 256)
		size, error := session.connection.Read(&buffer)

		if error != nil {
			fmt.Printf("Error in reading message! %s\n", error)
			return
		}

		fmt.Printf("Size of message read %d\n", size)

		response := delegateCommand(buffer[0], buffer, session)
		output := new(bytes.Buffer)
		error = binary.Write(output, binary.LittleEndian, response)

		if error != nil {
			fmt.Printf("Error in serializing response! %s\n", error)
			return
		}

		size, error = session.connection.Write(output.Bytes())

		if error != nil {
			fmt.Printf("Error in writing response! %s\n", error)
			return
		}

		fmt.Printf("Wrote %d bytes\n", size)

		// if session.done {
		// TODO: Move logic elsewhere
		// session.connection.Connection.Close()
		// session.done = true
		// }
	}

	fmt.Println("Session finished")
}

func delegateCommand(cmd uint8, data []byte, session *Authsession) interface{} {
	switch cmd {
	case model.AuthLogonChallenge:
		fmt.Println("AuthlogonChallenge registered")
		logonchallenge := model.LogonChallenge{}
		convertData(data, &logonchallenge)

		// fmt.Printf("Gamename: %s\n", logonchallenge.Gamename)
		// fmt.Printf("Accountname: %s\n", logonchallenge.Accountname)

		response, srp := Handlers.HandleLogonChallenge(logonchallenge, repository)

		fmt.Printf("SRP pointer returned: %p\n", srp)

		// fmt.Println("Response")
		// res2p, _ := json.Marshal(response)
		// fmt.Println(string(res2p))

		// fmt.Println("SRP")
		// srp2p, _ := json.Marshal(srp)
		// fmt.Println(string(srp2p))

		session.srp = srp
		session.done = false // Expect proof
		return response
	case model.AuthLogonProof:
		fmt.Println("AuthlogonProof registered")
		logonproof := model.LogonProof{}
		convertData(data, &logonproof)
		fmt.Printf("SRP pointer returned: %p\n", session.srp)

		response := Handlers.HandleLogonProof(logonproof, session.srp)
		session.done = false // Expect realmlist command after proof.
		return response
	case model.AuthReconnectChallenge:
		fmt.Println("AuthReconnectChallenge registered")
		return nil
	case model.AuthReconnectProof:
		fmt.Println("AuthReconnectProof registered")
		session.done = true // Dont expect anymore after this. Perhaps realmlist?
		return nil
	case model.RealmList:
		fmt.Println("Realmlist registered")
		return nil
	}

	return nil
}

func convertData(data []byte, result interface{}) {
	reader := bytes.NewReader(data)
	error := binary.Read(reader, binary.LittleEndian, result)

	if error != nil {
		fmt.Printf("Error in binary conversion: %s\n", error)
		panic(error)
	}
}
