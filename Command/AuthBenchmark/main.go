package main

import (
	"CitadelCore/AuthorisationServer/Model"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	ip := os.Args[1]
	port := os.Args[2]
	runtime := ""
	fmt.Printf("Starting benchmark for %s:%s\n", ip, port)

	if len(os.Args) > 3 {
		runtime = os.Args[3]
		fmt.Printf("Runtime has been specified: %s\n", runtime)
	}

	var gamename [4]byte
	copy(gamename[:], []byte("WoW"))
	var platform [4]byte
	copy(platform[:], []byte("x86"))
	var os [4]byte
	copy(os[:], []byte("Win"))
	var country [4]byte
	copy(country[:], []byte("enUS"))
	logonchallenge := Model.LogonChallenge{
		Command:         Model.AuthLogonChallenge,
		ProtocolVersion: 8,
		Size:            34,
		Gamename:        gamename,
		Major:           3,
		Minor:           3,
		Patch:           5,
		Build:           12340,
		Platform:        platform,
		Operatingsystem: os,
		Country:         country,
	}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, logonchallenge)
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", ip, port))

	if err != nil {
		fmt.Printf("Error in resolving address\n%e", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("Error in dialing address\n%e", err)
		return
	}
	conn.Write(buffer.Bytes())
}
