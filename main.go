package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}

type IRC struct {
	client net.Conn
}

func new_irc() (*IRC, error) {
	connection, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return nil, err
	}
	return &IRC{connection}, nil
}

func (irc *IRC) send_command(command, body string) error {
	if command == "" || body == "" {
		return fmt.Errorf("Command or body shouldn't be empty")
	}
	_, err := fmt.Fprintf(irc.client, "%s %s\n", command, body)
	return err
}

func (irc *IRC) send_message(channel_name, message string) error {
	if message == "" || channel_name == "" {
		return fmt.Errorf("Message or channel name should not be empty")
	}
	err := irc.send_command("PRIVMSG "+channel_name, ":"+message)
	return err
}

func (irc *IRC) send_pong_to_server() error {
	err := irc.send_command("PONG", ":tmi.twitch.tv")
	return err
}

const BUFFER_SIZE = 2040

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file")
	}

	var (
		OAUTH_TOKEN  = os.Getenv("OAUTH_TOKEN")
		BOT_NAME     = os.Getenv("BOT_NAME")
		CHANNEL_NAME = "#" + os.Getenv("CHANNEL_NAME")
	)

	irc, err := new_irc()
	if err != nil {
		log.Fatal("Unable to establish connection to IRC server")
	}

	wg.Add(2)
	go func() {
		irc.send_command("PASS", OAUTH_TOKEN)
		irc.send_command("NICK", BOT_NAME)
		irc.send_command("JOIN", CHANNEL_NAME)
		wg.Done()
	}()

	go func() {
		for {
			received_data := make([]byte, BUFFER_SIZE)
			received_data_size, err := irc.client.Read(received_data)
			if err != nil {
				log.Fatal("Unable to read data from socket")
			}

			message := strings.Split(string(received_data), "\n")[0]

			if received_data_size > 0 {
				parser := new_parser(message)
				nickname, message, err := parser.get_message()

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(nickname)
				fmt.Println(message)
				if strings.HasPrefix(message, "PING") {
					err := irc.send_pong_to_server()
					if err != nil {
						log.Fatal("Unable to send PONG")
					}
				} else if strings.HasPrefix(message, ":") {
					err := irc.send_message(CHANNEL_NAME, "Testing my bot")
					if err != nil {
						log.Fatal("Unable to send message")
					}
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
