package simple

import "fmt"

type Connection struct {
	*File
}

func NewConnection(f *File) (*Connection, func()) {
	connection := &Connection{File: f}
	return connection, func() {
		connection.Close()
	}
}

func (c *Connection) Close() {
	fmt.Println("Close connection", c.File.Name)
}
