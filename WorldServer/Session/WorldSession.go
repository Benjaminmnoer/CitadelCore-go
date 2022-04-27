package Session

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var clientQueue []net.Conn
var numberOfClients int
var nHandlers int
var queueLock = sync.Mutex{}

const maxNumberOfConnections = 3000

func StartSession() {
	clientQueue = make([]net.Conn, maxNumberOfConnections)
	nHandlers = 4

	server, _ := net.Listen("tcp", "127.0.0.1:8085")

	for i := 0; i < nHandlers; i++ {
		go handleConnections()
	}

	for {
		conn, err := server.Accept()

		if err != nil {
			fmt.Printf("Error while accepting:\n%e", err)
		}

		if numberOfClients < maxNumberOfConnections {
			clientQueue = append(clientQueue, conn)
		} else {
			// TODO: Send message, take stats
			conn.Close()
		}
	}
}

func handleConnections() {
	for {
		queueLock.Lock()
		defer queueLock.Unlock()

		if len(clientQueue) <= 0 {
			time.Sleep(1 * time.Millisecond)
			continue
		}

		connection := clientQueue[0]
		clientQueue = clientQueue[1:]

		queueLock.Unlock()

		buffer := make([]byte, 256)
		_, err := connection.Read(buffer)

		if err != nil {
			fmt.Printf("Error in connection with %s:\n%e\n", connection.RemoteAddr().String(), err)
			continue
		}

		// TODO: Handle data

		queueLock.Lock()
		defer queueLock.Unlock()
		clientQueue = append(clientQueue, connection)
	}
}

func GetQueueDepth() int {
	return numberOfClients
}
