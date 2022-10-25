package server

import "net"

type client struct {
	con         net.Conn
	username    string
	msgChan     chan string
	currentRoom *room
}

func newClient(con net.Conn, username string, msgChan chan string, currentRoom *room) *client {
	return &client{con: con, username: username, msgChan: msgChan, currentRoom: currentRoom}
}

func (c *client) setUsername(username string) {
	c.username = username
}
