package server

import (
	"bufio"
	"fmt"
	"io"
	"sync"
)

type client struct {
	// con         net.Conn
	uid         string
	reader      io.Reader
	writer      io.Writer
	username    string
	msgChan     chan string
	currentRoom *room
}

func newClient(
	uid string, reader io.Reader, writer io.Writer, username string, msgChan chan string, currentRoom *room,
) *client {
	return &client{uid: uid, reader: reader, writer: writer, username: username, msgChan: msgChan,
		currentRoom: currentRoom}
}

func (cli *client) setUsername(username string) {
	cli.username = username
}

// func to get username from a client - called first before main handle connections
func (cli *client) getUsername() (string, error) {
	fmt.Println("waiting for usernme...")
	input, err := bufio.NewReader(cli.reader).ReadString('\n')
	if err != nil {
		return "", err
	}
	username := FormatStringInput(input)
	fmt.Printf("username for client %q set to %q\n", cli.uid, username)
	return username, nil
}

// goroutine to monitor the msgChan channel and send to clients
func (cli *client) monitorMsgChan() {
	for {
		msg := <-cli.msgChan // await change in msgChan

		// send msg to everone in room except the client who sent it
		for _, client := range cli.currentRoom.clients {
			if cli != client {
				err := writeToCon(client.writer, msg)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

// goroutine to read from port
func (cli *client) handleCon() {
	for {
		msg, err := bufio.NewReader(cli.reader).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Server received: %q from user: %q in room: %q\n", msg, cli.username, cli.currentRoom.name)

		cli.msgChan <- cli.username + ": " + msg
	}
}

func (cli *client) newClientSetup(rooms *rooms) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	username, err := cli.getUsername()
	if err != nil {
		fmt.Println(err)
	}
	cli.setUsername(username)
	writeToCon(
		cli.writer,
		fmt.Sprintf("Welcome %q to the %q channel", cli.username, cli.currentRoom.name),
	)
	cli.msgChan <- username + " has entered the chat"

	// start go routine for monitoring the socket
	wg.Add(1)
	go func() {
		defer wg.Done()
		cli.handleCon()
	}()

	// start goroutine for monitoring the message channel for new client
	wg.Add(1)
	go func() {
		defer wg.Done()
		cli.monitorMsgChan()
	}()
}
