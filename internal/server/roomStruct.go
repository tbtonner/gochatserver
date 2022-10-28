package server

import "fmt"

type room struct {
	name    string
	clients []*client
}

// constructor for the room struct
func newRoom(name string) *room {
	return &room{name: name, clients: make([]*client, 0)}
}

// set default room (the "main" rooom)
func newDefaultRoom() *room {
	return newRoom("main")
}

// adds a new client to a room
func (r *room) addCliToRoom(cli *client) {
	r.clients = append(r.clients, cli)
	fmt.Println(r.clients)
}

// removes a client from a room
func (r *room) removeCli(cli *client) {
	for i, client := range r.clients {
		if client.uid == cli.uid {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
		}
	}
}
