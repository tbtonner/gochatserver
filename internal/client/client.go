package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/tbtonner/gochatserver/internal/server"
)

// TODO: make useful if only one line
func printMsg(msg string) {
	// print msg (remove current line and replace)
	fmt.Printf("\033[2K\r%s\n> ", server.FormatStringInput(msg))
}

// goroutine to read from port
func monitorSocket(con net.Conn) {
	for {
		msg, err := bufio.NewReader(con).ReadString('\n')
		if err != nil {
			fmt.Println("Unable to read input from the server:", err.Error())
			os.Exit(1) // TODO: handle instead of exiting
		}
		printMsg(msg)
	}
}

// read evaluate and print (used in rep loop to get user input)
func rep(con net.Conn, prefix string, scanner bufio.Scanner) error {
	fmt.Printf("%s ", prefix)
	scanner.Scan()

	err := scanner.Err()
	if err != nil {
		return err
	}

	_, err = con.Write([]byte(scanner.Text() + "\n")) // write to socket
	if err != nil {
		return err
	}

	return nil
}

// goroutine to send message to port
func sendToSocket(con net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	err := rep(con, "Please enter your username:", *scanner)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Welcome to the 'main' chatroom!")

	// repl start
	for {
		err = rep(con, ">", *scanner)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// main for client (to be called from runClient)
func RunClient() {
	wg := sync.WaitGroup{}

	con, err := net.Dial(server.PROTOCOL, server.PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to server")

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
