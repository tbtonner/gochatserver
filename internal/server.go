package internal

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// goroutine to read from port
func handleCon(con net.Conn) {
	defer con.Close()
	for {
		data, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(data)
	}
}

// main for server (to be called from runServer)
func RunServer() {
	var i int
	wg := sync.WaitGroup{}

	// set up listener on port
	ln, err := net.Listen("tcp", ":8080")
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

		// start go routine for monitoring the socket
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleCon(con)
		}()

		i++
		if i == 3 { // for now only allow two clients to join -> close connection if 3rd one joins
			con.Close()
		}
	}

	wg.Wait()
	fmt.Println("closing server...")
}
