package Communication

type SessionHandler interface {
	// Handles session. Returns true for returning the connection to the back of the queue, false if it should not.
	HandleSession(Client, []byte) ([]byte, SessionStatus, error)
}

type SessionStatus int8

const (
	EndConnection = iota
	EnqueueClient
	KeepClient
)
