package server

import "fmt"

type room struct {
	name    string
	clients []*client
}

func newRoom(name string) *room {
	return &room{name: name, clients: make([]*client, 0)}
}

func newDefaultRoom() *room {
	return newRoom("main")
}

func (r *room) addCliToRoom(cli *client) {
	r.clients = append(r.clients, cli)
	fmt.Println(r.clients)
}

func (r *room) removeCli(cli *client) {
	for i, client := range r.clients {
		if client.uid == cli.uid {
			r.clients = append(r.clients[:i], r.clients[i+1:]...)
		}
	}
}
