package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

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
	_, err := writer.Write([]byte(msg + "\000"))
	return err
}

// main for server (to be called from runServer)
func RunServer() {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	defaultRoom := newDefaultRoom()
	rooms := rooms{}
	rooms.addRoom(defaultRoom)

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
			break
		}

		// create new client instance and add to the clients slice (set to default room)
		newClient := newClient(
			con.RemoteAddr().String(),
			con,
			con,
			"annoymous",
			make(chan string, 1),
			make(chan string, 1),
			defaultRoom,
		)
		defaultRoom.addCliToRoom(newClient)
		fmt.Printf("client: %q connected\n", newClient.uid)

		// go routine for this client (monitor and set username, write etc.)
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer con.Close()
			newClient.newClientSetup(&rooms)
			// connection been terminated at this point:
			newClient.sendToAllBarMe(fmt.Sprintf("%q has left the chat", newClient.username))
			fmt.Printf("closing connection: %q\n", newClient.uid)
		}()
	}

	fmt.Println("closing server...")
}
