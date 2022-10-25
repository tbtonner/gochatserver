package server

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
}
