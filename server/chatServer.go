package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

// goroutine to read from port
func monitorSocket(con net.Conn) {
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
func Run() {
	wg := sync.WaitGroup{}

	// set up listener on port
	ln, err := net.Listen("tcp", ":8080")
	defer ln.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

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
		monitorSocket(con)
	}()

	wg.Wait()
}
