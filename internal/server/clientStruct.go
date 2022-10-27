package server

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

type client struct {
	uid         string
	reader      io.Reader
	writer      io.Writer
	username    string
	msgChan     chan string
	comChan     chan string
	currentRoom *room
}

// constructor for client struct
func newClient(
	uid string, reader io.Reader, writer io.Writer, username string, msgChan, comChan chan string, currentRoom *room,
) *client {
	return &client{uid: uid, reader: reader, writer: writer, username: username, msgChan: msgChan, comChan: comChan,
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

// func to send msg to everone in room except the client who sent it
func (cli *client) sendToAllBarMe(msg string) {
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

// goroutine to monitor the msgChan channel
func (cli *client) monitorMsgChan() {
	for {
		msg, ok := <-cli.msgChan // await change in msgChan
		if !ok {
			fmt.Printf("msgChan not OK for %q\n", cli.uid)
			return
		}

		cli.sendToAllBarMe(msg)
	}
}

// goroutine to monitor the comChan channel - then execute said command
func (cli *client) monitorComChan(rooms *rooms) {
	for {
		com, ok := <-cli.comChan // await change in comChan
		if !ok {
			fmt.Printf("comChan not OK for %q\n", cli.uid)
			return
		}

		splitList := strings.Split(com, " ")
		switch splitList[0] {
		case "create":
			cli.comCreate(splitList[1:], rooms)
		case "join":
			cli.comJoin(splitList[1:], rooms)
		case "shout":
			cli.comShout(splitList[1:])
		case "whisper":
			cli.comWhisper(splitList[1:])
		case "help":
			cli.comHelp()
		case "kick":
			cli.comKick(splitList[1:])
		case "spam":
			cli.comSpam(splitList[1:])
		default:
			cli.comNoneFound(splitList[0])
		}

	}
}

// goroutine to read from port
func (cli *client) handleCon() {
	for {
		data, err := bufio.NewReader(cli.reader).ReadString('\n')
		if err != nil {
			fmt.Printf("Server lost connection to client: %q\n", cli.uid)
			cli.currentRoom.removeCli(cli)
			close(cli.msgChan)
			return
		}
		fmt.Printf("Server received: %q from user: %q in room: %q\n", data, cli.username, cli.currentRoom.name)

		// if / found at start (ie. a command) then send to command channel, else send to default message channel
		if data[0:1] == "/" {
			cli.comChan <- FormatStringInput(data[1:])
		} else {
			cli.msgChan <- cli.username + ": " + data
		}
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

	// start goroutine for monitoring the commands channel for new client
	wg.Add(1)
	go func() {
		defer wg.Done()
		cli.monitorComChan(rooms)
	}()

}
