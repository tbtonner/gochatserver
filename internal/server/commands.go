package server

import (
	"fmt"
	"strconv"
	"strings"
)

// func to see if room name provided already exists
func roomExist(rname string, rooms *rooms) bool {
	for _, room := range rooms.rms {
		if room.name == rname {
			return true
		}
	}

	return false
}

// function to join a client to a new room (leave old one and update 'rooms' to match new state)
func (cli *client) joinNewRoom(newRoom *room) {
	cli.currentRoom.removeCli(cli)
	cli.currentRoom = newRoom

	newRoom.addCliToRoom(cli)

	writeToCon(
		cli.writer,
		fmt.Sprintf("Welcome %q to the %q channel", cli.username, cli.currentRoom.name),
	)
	cli.msgChan <- cli.username + " has entered the chat"
}

// command to create a new room
func (cli *client) comCreate(args []string, rooms *rooms) {
	fmt.Printf("create room command: %v\n", args)

	if args[0] != "" && !roomExist(args[0], rooms) {
		newRoom := newRoom(args[0])
		rooms.addRoom(newRoom)

		cli.joinNewRoom(newRoom)
		writeToCon(cli.writer, fmt.Sprintf("Room created successfully"))
	} else {
		writeToCon(cli.writer, fmt.Sprintf("Room name invalid or already taken"))
	}
}

// command to join a new room
func (cli *client) comJoin(args []string, rooms *rooms) {
	fmt.Printf("join room command: %v\n", args)
	for _, room := range rooms.rms {
		if room.name == args[0] {
			cli.joinNewRoom(room)
			return
		}
	}

	writeToCon(cli.writer, fmt.Sprintf("Room: %q does not exist", args[0]))
}

// command to shout a given message (CAPS)
func (cli *client) comShout(args []string) {
	fmt.Printf("shout command: %v\n", args)
	msg := strings.ToUpper(strings.Join(args, " "))
	cli.msgChan <- cli.username + ": " + msg
}

// TODO: command to whisper a message to a specific user
func (cli *client) comWhisper(args []string) {
	fmt.Printf("whisper command: %v\n", args)
	writeToCon(cli.writer, fmt.Sprintf("This feature is yet to be implimented :("))
}

// TODO: command to kick a user from the current room
func (cli *client) comKick(args []string) {
	fmt.Printf("kick command: %v\n", args)
	writeToCon(cli.writer, fmt.Sprintf("This feature is yet to be implimented :("))
}

// command to spam a message x times
func (cli *client) comSpam(args []string) {
	fmt.Printf("spam command: %v\n", args)

	num, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		writeToCon(cli.writer, fmt.Sprintf("%q is not a number", args[0]))
		return
	}

	spamMsg := ""
	for i := 0; i < int(num); i++ {
		spamMsg += fmt.Sprintf("%s: %s\n", cli.username, args[1])
	}
	cli.msgChan <- strings.TrimSuffix(spamMsg, "\n")
}

// the help command - lists all the commands to user's screen
func (cli *client) comHelp() {
	fmt.Printf("help command\n")
	writeToCon(cli.writer, fmt.Sprintf(`
'/create' <room_name>: creates a new room with a given room_name if possible and joins it
'/join' <room_name>: joins room with given room_name if possible
'/shout' <msg>: capitalises msg and sends to current room
'/whisper' <person> <msg>: sends msg only to person given
'/kick' <person>: kicks person from room
'/spam' <times> <msg>: spams a msg a number of times specified in current room
`))
}

// func to print to user when user enters a command not found
func (cli *client) comNoneFound(com string) {
	fmt.Printf("Command %q not found\n", com)
	writeToCon(cli.writer, fmt.Sprintf("Command %q not found", com))
}
