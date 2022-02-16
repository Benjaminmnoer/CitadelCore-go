package Connection

import (
	"fmt"
	"net"
)

type TcpConnection struct {
	Connection       net.Conn
	DeadlineDuration string
}

func CreateTcpConnection(connection net.Conn, deadlineDuration string) TcpConnection {
	return TcpConnection{
		Connection:       connection,
		DeadlineDuration: deadlineDuration,
	}
}

func (connection TcpConnection) Write(data []byte) (n int, error error) {
	fmt.Printf("Writing to %s\n", connection.Connection.RemoteAddr())
	// durr, _ := time.ParseDuration(connection.DeadlineDuration)
	// conn.SetWriteDeadline(time.Now().Add(durr))
	res, err := connection.Connection.Write(data)
	// conn.SetWriteDeadline(time.Time{})
	return res, err
}

func (connection TcpConnection) Read(data *[]byte) (n int, error error) {
	fmt.Printf("Reading from %s\n", connection.Connection.RemoteAddr())
	// durr, _ := time.ParseDuration(connection.DeadlineDuration)
	// conn.SetReadDeadline(time.Now().Add(durr))
	res, err := connection.Connection.Read(*data)
	// conn.SetReadDeadline(time.Time{})
	return res, err
}
