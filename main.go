package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"strings"
)

type Bot struct {
	conn                                  *textproto.Conn
	server, port, channel, nick, username string
	done                                  chan bool
	reader                                *bufio.Reader
}

func NewBot(server, port, channel, nick, username string) *Bot {
	return &Bot{
		server:   server,
		port:     port,
		channel:  channel,
		nick:     nick,
		username: username,
	}
}

func (bot Bot) Connect() {
	fmt.Printf("Connecting to %s .. ", bot.server)
	connection_string := net.JoinHostPort(bot.server, bot.port)
	conn, err := textproto.Dial("tcp", connection_string)
	if err != nil {
		fmt.Println(err)
	}

	bot.conn = conn
	go bot.ListenToMessage()
	fmt.Println(bot)
	bot.conn.PrintfLine("NICK %s", bot.nick)
	bot.conn.PrintfLine("USER %s 0 * :%s", bot.nick, bot.username)
	bot.conn.PrintfLine("JOIN %s", bot.channel)
	bot.SendMessage("Hi")
	bot.SendUserInput()
	<-bot.done
	fmt.Println("Quitting")
}

func (bot Bot) SendMessage(message string) {
	bot.conn.PrintfLine("PRIVMSG %s :%s", bot.channel, message)
}

func (bot Bot) SendCommand(message string) {
	bot.Quit()
}

func (bot Bot) Quit() {
	bot.conn.PrintfLine("QUIT :bye")
	bot.done <- true
	bot.conn.Close()
}

func (bot Bot) ListenToMessage() {
	for {
		message, err := bot.conn.ReadLine()
		if err != nil {
			bot.conn.Close()
		}
		fmt.Println(message)
		index := strings.Index(message, "PING")
		if index == 0 {
			parts := strings.Split(message, ":")
			response := "PONG " + parts[1]
			fmt.Println("Received PING, responding with ", response)
			bot.conn.PrintfLine(response)
		}
	}
}

func (bot Bot) SendUserInput() {
	bot.reader = bufio.NewReader(os.Stdin)
	for {

		user_input, _ := bot.reader.ReadString('\n')
		// if user input starts with / that means it is a command.
		index := strings.Index(user_input, "/")
		if index == 0 {
			bot.SendCommand(user_input)
		}
		bot.SendMessage(user_input)
	}
}

func main() {

	bot := NewBot("irc.freenode.net", "6667", "#botwar", "ragsag1", "Botty")
	bot.Connect()

}
