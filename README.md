# Go Tasks - Task 2: Chat server
## Basic
Create a chat server which two clients (A and B) can connect to over a TCP socket (see the first two examples in [the go docs](https://pkg.go.dev/net) as a starting point). The chat client should be a program which you can run for each client which prompts the user for input and sends it to the server, and prints anything the server sends it. 

Any message A sends, the server should send to B, and vice versa. If a third client tries to connect to the server the server should reject it and close its connection.

#### Hints: 

* Each connection to the server should have a goroutine that reads from the socket
* Each connection to the server should have a goroutine that reads from a channel
* Once a message terminated by a newline is read from the socket the goroutine should send it to the other connection’s channel
* The goroutine that reads from a channel should send what was on the channel to its socket
* The client should use a [bufio.Scanner](https://pkg.go.dev/bufio#Scanner) to read input from stdin
* Is a Scanner the best type to use in this scenario?
* What other alternatives are there (hint use the documentation)?
* The client should have a goroutine that reads out of the socket and prints it on the screen

## Medium
When starting the server have a password passed as a command line argument. The client should prompt the user for a password and send it to the server. If the password is wrong then the connection should be rejected. Otherwise it should be accepted and the behaviour should be as above.

Would you use this method of authentication on a real server?

Try to use a command line parsing library like [cobra](https://github.com/spf13/cobra) to get the password.

#### Hints
* After connecting the first thing the client should do is send the password followed by a newline
* The first thing the server should do after accepting a connection is read until a newline

## Medium II
Talking to one person is fun, but talking to lots of people is even more fun. When a client connects it should tell the server what chat room it wants to join and their name. Now when a client sends a message it should be sent to every client in the room (apart from itself), specifying the name.

#### Hints
* You need a thread-safe way of mapping a room name to the connections and their names

## Advanced
Give your code to one of your colleagues. Have them run it and try to make it crash or find a bug. Fix the crash/bug if any.

Then implement a `“whisper”` command. When a client sends a message like:

```/whisper matt why are you so cool and great```

The client who has identified themselves as matt, and only them, should receive the message “why are you so cool and great”.

Implement other commands. For example:
* `/spam <username> <times> <message>` - send <message> to <username> <times` times
* `/shout <message>` - send <message> to the room but ALL IN CAPS
* `/kick <username>` - have the server close <username>’s connection
* Something of your own design
