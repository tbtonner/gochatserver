package server

import (
	"fmt"
	"strings"
)

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

func (cli *client) comCreate(args []string, rooms *rooms) {
	fmt.Printf("create room command: %v\n", args)
	newRoom := newRoom(args[0])
	rooms.addRoom(newRoom)

	cli.joinNewRoom(newRoom)

	writeToCon(cli.writer, fmt.Sprintf(""))
}

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

func (cli *client) comShout(args []string) {
	fmt.Printf("shout command: %v\n", args)
	msg := strings.ToUpper(strings.Join(args, " "))
	cli.msgChan <- cli.username + ": " + msg
}

func (cli *client) comWhisper(args []string) {
	fmt.Printf("whisper command: %v\n", args)
	writeToCon(cli.writer, fmt.Sprintf(""))
}

func (cli *client) comKick(args []string) {
	fmt.Printf("kick command: %v\n", args)
	writeToCon(cli.writer, fmt.Sprintf(""))
}

func (cli *client) comSpam(args []string) {
	fmt.Printf("spam command: %v\n", args)
	writeToCon(cli.writer, fmt.Sprintf(""))
}

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

func (cli *client) comNoneFound(com string) {
	fmt.Printf("Command %q not found\n", com)
	writeToCon(cli.writer, fmt.Sprintf("Command %q not found", com))
}
