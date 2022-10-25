package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	PORT     string = ":8080"
	PROTOCOL string = "tcp"
)

func FormatStringInput(s string) string {
	return strings.Trim(s, "\r\n")
}

// goroutine to monitor the msgChan channel and send to clients
func monitorMsgChan(cli *client) {
	for {
		msg := <-cli.msgChan // await change in msgChan

		roomToSend := cli.currentRoom
		for i := 0; i < len(roomToSend.clients); i++ {
			if roomToSend.clients[i].con != cli.con {
				_, err := roomToSend.clients[i].con.Write([]byte(cli.username + ": " + msg + "\n"))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

// goroutine to read from port
func handleCon(cli *client) {
	defer cli.con.Close()
	for {
		msg, err := bufio.NewReader(cli.con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Server received:" + msg)

		cli.msgChan <- msg
	}
}

func getUsername(cli *client) (string, error) {
	fmt.Println("waiting for usernme...")
	input, err := bufio.NewReader(cli.con).ReadString('\n')
	if err != nil {
		return "", err
	}
	username := FormatStringInput(input)
	fmt.Println("username for client set to " + username)
	return username, nil
}

// main for server (to be called from runServer)
func RunServer() {
	wg := sync.WaitGroup{}
	defaultRoom := newDefaultRoom()
	// rooms := []room{defaultRoom} // TODO: need for later! -> make concurrent

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
			con,
			"annoymous",
			make(chan string, 1),
			defaultRoom,
		)
		defaultRoom.addCliToRoom(newClient)

		username, err := getUsername(newClient)
		if err != nil {
			fmt.Println(err)
		}
		newClient.setUsername(username)

		// start go routine for monitoring the socket
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleCon(newClient)
		}()

		// start goroutine for monitoring the message channel for new client
		wg.Add(1)
		go func() {
			defer wg.Done()
			monitorMsgChan(newClient)
		}()

		// for now only allow two clients to join -> close connection if 3rd one joins
		if len(defaultRoom.clients) == 3 {
			fmt.Println("Too many clients - closing connection...")
			con.Close()
			fmt.Println("Connection closed")
		}
	}

	wg.Wait()
	fmt.Println("closing server...")
}
