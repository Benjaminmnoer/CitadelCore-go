package Communication

import (
	"CitadelCore/Shared/DataStructures"
	"fmt"
	"time"
)

type ClientHandler struct {
	clients         DataStructures.ConcurrentQueue
	sessionHandler  SessionHandler
	numberOfThreads int
}

func CreateClientHandler(sessionHandler SessionHandler, nThreads int) *ClientHandler {
	return &ClientHandler{clients: DataStructures.ConcurrentQueue{}, sessionHandler: sessionHandler, numberOfThreads: nThreads}
}

func (ch *ClientHandler) AddClient(client Client) {
	ch.clients.Enqueue(client)
}

func (ch *ClientHandler) getClient() (Client, error) {
	client := ch.clients.Dequeue()

	if client == nil {
		return Client{}, fmt.Errorf("the queue is empty")
	}

	return client.(Client), nil
}

func (ch *ClientHandler) Start() error {
	for i := 0; i < ch.numberOfThreads; i++ {
		go ch.handleConnection()
	}

	return nil
}

func (ch *ClientHandler) handleConnection() error {
	for {
		client, err := ch.getClient()
		if err != nil {
			time.Sleep(5 * time.Millisecond)
			continue
		}

	READ:
		data, err := client.Read()

		if err != nil {
			// TODO: Log error and continue.
			continue
		}

		response, status, err := ch.sessionHandler.HandleSession(client, data)

		if err != nil {
			// TODO: Log error and handle connection.
			client.Dispose()
			continue
		}

		err = client.Write(response)

		if err != nil {
			// TODO: Log error and handle connection.
			client.Dispose()
			continue
		}

		if status == KeepClient {
			goto READ
		} else if status == EndConnection {
			client.Dispose()
		} else if status == EnqueueClient {
			ch.AddClient(client)
		} else {
			client.Dispose()
			return fmt.Errorf("unhandled client status message")
		}
	}
}
