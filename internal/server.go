package internal

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type client struct {
	con      net.Conn
	username string
	msgChan  chan string
}

var (
	clients []client
)

const (
	PORT     string = ":8080"
	PROTOCOL string = "tcp"
)

func handleSendMsg(cli client) {
	for {
		msg := <-cli.msgChan
		for i := 0; i < len(clients); i++ {
			if clients[i].con != cli.con {
				_, err := clients[i].con.Write([]byte(msg + "\n"))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

// goroutine to read from port
func handleCon(cli client) {
	defer cli.con.Close()
	for {
		data, err := bufio.NewReader(cli.con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Server received:" + data)

		cli.msgChan <- data
	}
}

// main for server (to be called from runServer)
func RunServer() {
	var i int
	wg := sync.WaitGroup{}

	// set up listener on port
	ln, err := net.Listen(PROTOCOL, PORT)
	defer ln.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// confirgure connection
		con, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		var newClient client
		newClient.con = con
		newClient.username = "testOther"
		newClient.msgChan = make(chan string, 1)

		clients = append(clients, newClient)

		// start go routine for monitoring the socket
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleCon(newClient)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			handleSendMsg(newClient)
		}()

		i++
		if i == 3 { // for now only allow two clients to join -> close connection if 3rd one joins
			con.Close()
		}
	}

	wg.Wait()
	fmt.Println("closing server...")
}
