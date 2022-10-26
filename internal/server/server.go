package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type rooms struct {
	rms []room
	mu  sync.Mutex
}

func (rooms *rooms) addRoomToRooms(room room) {
	rooms.mu.Lock()
	defer rooms.mu.Unlock()

	rooms.rms = append(rooms.rms, room)
}

const (
	PORT     string = ":8080"
	PROTOCOL string = "tcp"
)

// func to format string input before sending to clients
func FormatStringInput(s string) string {
	return strings.Trim(s, "\r\n")
}

// func to send a given msg to a client
func writeToCon(writer io.Writer, msg string) error {
	_, err := writer.Write([]byte(msg + "\n"))
	return err
}

// main for server (to be called from runServer)
func RunServer() {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	defaultRoom := newDefaultRoom()
	rooms := rooms{}
	rooms.addRoomToRooms(*defaultRoom)

	// set up listener on port
	ln, err := net.Listen(PROTOCOL, PORT)
	defer ln.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to port: " + PORT + " via: " + PROTOCOL)

	// loop to accept new net.dials (from clients)
	for {
		// confirgure connection
		con, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		// create new client instance and add to the clients slice (set to default room)
		newClient := newClient(
			con.RemoteAddr().String(),
			con,
			con,
			"annoymous",
			make(chan string, 1),
			defaultRoom,
		)
		defaultRoom.addCliToRoom(newClient)
		fmt.Printf("client: %q connected\n", newClient.uid)

		// // for now only allow two clients to join -> close connection if 3rd one joins
		// if len(defaultRoom.clients) == 3 {
		// 	fmt.Println("Too many clients - closing connection...")
		// 	con.Close()
		// 	fmt.Println("Connection closed")
		// 	break
		// }

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer con.Close()
			newClient.newClientSetup(&rooms)
		}()
	}

	fmt.Println("closing server...")
}
