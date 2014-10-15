package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"strings"
)

func ListenToMessage(conn *textproto.Conn, done chan bool) {
	for {
		message, err := conn.ReadLine()
		if err != nil {
			done <- true
		}
		fmt.Println(message)
		index := strings.Index(message, "PING")
		if index == 0 {
			parts := strings.Split(message, ":")
			response := "PONG " + parts[1]
			fmt.Println("Received PING, responding with ", response)
			conn.PrintfLine(response)
		}
	}
}

func SendUserInput(conn *textproto.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {

		input, _ := reader.ReadString('\n')
		conn.PrintfLine("PRIVMSG #botwar :" + input)
	}
}

func main() {
	done := make(chan bool)
	fmt.Println("Connecting to freenode")
	conn, err := net.Dial("tcp", "irc.freenode.net:6667")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	new_conn := textproto.NewConn(conn)
	go ListenToMessage(new_conn, done)

	new_conn.PrintfLine(":Macha!~macha@unaffiliated/macha PRIVMSG #botwar :Test response")
	new_conn.PrintfLine("NICK ragsagar1")
	new_conn.PrintfLine("USER ragsagar1 0 * :Bot")

	new_conn.PrintfLine("JOIN #botwar")
	new_conn.PrintfLine("PRIVMSG #botwar :Hi")
	go SendUserInput(new_conn)
	<-done
}
