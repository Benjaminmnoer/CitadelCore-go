package Communication

// Interface for delegating commands and returning a response object
type CommandDelegater interface {
	Delegate([]byte) ([]byte, error)
}
