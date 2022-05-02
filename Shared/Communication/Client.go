package Communication

import "net"

// Type holding the connection to the client.
type Client struct {
	Connection net.Conn
}

func (c Client) GetEndpoint() string {
	return c.Connection.RemoteAddr().String()
}

func (c Client) Write(message []byte) error {
	_, err := c.Connection.Write(message)
	return err
}

func (c Client) Read() ([]byte, error) {
	buffer := make([]byte, 1024)
	_, err := c.Connection.Read(buffer)
	return buffer, err
}

func (c Client) Dispose() error {
	return c.Connection.Close()
}
