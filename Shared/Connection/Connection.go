package Connection

import (
	"net"
	"time"
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
	conn := connection.Connection
	durr, _ := time.ParseDuration(connection.DeadlineDuration)
	conn.SetWriteDeadline(time.Now().Add(durr))
	return conn.Write(data)
}

func (connection TcpConnection) Read(data *[]byte) (n int, error error) {
	conn := connection.Connection
	durr, _ := time.ParseDuration(connection.DeadlineDuration)
	conn.SetReadDeadline(time.Now().Add(durr))
	return conn.Read(*data)
}
