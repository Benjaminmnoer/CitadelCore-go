package Communication

import (
	"CitadelCore/Shared/DataStructures"
	"fmt"
	"time"
)

type ClientHandler struct {
	clients         DataStructures.ConcurrentQueue
	numberOfThreads int
}

func CreateClientHandler(nThreads int) ClientHandler {
	return ClientHandler{clients: DataStructures.ConcurrentQueue{}, numberOfThreads: nThreads}
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

		data, err := client.Read()

		if err != nil {
			// TODO: Log error and continue
			continue
		}

		size := len(data)

		fmt.Printf("%d", size)
	}
}
