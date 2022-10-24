package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// goroutine to read from port
func monitorSocket(con net.Conn) {
	for {
		msg, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read input from the server:", err.Error())
			os.Exit(1) // TODO: handle instead of exiting
		}
		msg = strings.Trim(msg, "\r\n")

		fmt.Printf("\033[2K\rother: %s\n", msg)
		fmt.Print("> ")
	}
}

// goroutine to send message to port
func sendToSocket(con net.Conn) {
	// repl start
	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
		}

		_, e := con.Write([]byte(scanner.Text() + "\n")) // write to socket
		if err != nil {
			fmt.Println(e)
		}
	}
}

// main for client (to be called from runClient)
func RunClient() {
	wg := sync.WaitGroup{}

	con, err := net.Dial(PROTOCOL, PORT)
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

	// start go routine for sending to the socket
	wg.Add(1)
	go func() {
		defer wg.Done()
		sendToSocket(con)
	}()

	wg.Wait()
}
