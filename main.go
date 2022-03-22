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
	channel_name string
	bot_name     string
	oauth_token  string
	client       net.Conn
}

func new_irc(channel_name, bot_name, oauth_token string) (*IRC, error) {
	connection, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return nil, err
	}

	if channel_name == "" || bot_name == "" || oauth_token == "" {
		return nil, fmt.Errorf("Fields should not be empty")
	}

	return &IRC{channel_name, bot_name, oauth_token, connection}, nil
}

func (irc *IRC) send_command(command, body string) error {
	if command == "" || body == "" {
		return fmt.Errorf("command or body shouldn't be empty")
	}
	_, err := fmt.Fprintf(irc.client, "%s %s\n", command, body)
	return err
}

func (irc *IRC) send_message(message string) error {
	if message == "" {
		return fmt.Errorf("message or channel name should not be empty")
	}
	err := irc.send_command("PRIVMSG "+irc.channel_name, ":"+message)
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
		log.Fatal(err)
	}

	var (
		OAUTH_TOKEN  = os.Getenv("OAUTH_TOKEN")
		BOT_NAME     = os.Getenv("BOT_NAME")
		CHANNEL_NAME = "#" + os.Getenv("CHANNEL_NAME")
	)

	commands := map[string]string{
		"list_of_commands": "!me, !bot, !socials, !projects, !colors",
		"me":               "My name is Hícaro, I don't stream that much, but I hope you like it <3",
		"socials":          "Twitter: https://twitter.com/DanrlleyHicaro",
		"projects":         "All my projects are open-source. You can find them on https://github.com/HicaroD",
		"colors":           "https://github.com/HicaroD/Icarus",
		"bot":              "This bot is one of my projects and it was written in Go. You can find it here: https://github.com/HicaroD/TwitchBot",
		"today":            "No tasks today.",
	}

	irc, err := new_irc(CHANNEL_NAME, BOT_NAME, OAUTH_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(2)
	go func() {
		irc.send_command("PASS", irc.oauth_token)
		irc.send_command("NICK", irc.bot_name)
		irc.send_command("JOIN", irc.channel_name)
		wg.Done()
	}()

	go func() {
		for {
			received_data := make([]byte, BUFFER_SIZE)
			received_data_size, err := irc.client.Read(received_data)
			if err != nil {
				log.Fatal(err)
			}

			raw_message := strings.Split(string(received_data), "\n")[0]

			if received_data_size > 0 {
				parser := new_parser(raw_message)
				nickname, parsed_message, err := parser.parse()

				if err != nil {
					log.Fatal(err)
				}
				if strings.HasPrefix(raw_message, "PING") {
					go func() {
						err := irc.send_pong_to_server()
						if err != nil {
							log.Fatal(err)
						}
					}()
				}

				if strings.HasPrefix(parsed_message, "!commands") {
					go func() {
						err := irc.send_message(commands["list_of_commands"])
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!bot") {
					go func() {
						err := irc.send_message(commands["bot"])
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!me") {
					go func() {
						err := irc.send_message(commands["me"])
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!projects") {
					go func() {
						err := irc.send_message(commands["projects"])
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!socials") {
					go func() {
						err := irc.send_message(commands["socials"])
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!colors") {
					go func() {
						err := irc.send_message(commands["colors"])
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!today") {
					go func() {
						if nickname == "hicaro____" {
							message := strings.TrimPrefix(parsed_message, "!today")
							if strings.TrimSpace(message) == "" {
								err := irc.send_message(commands["today"])
								if err != nil {
									log.Fatal(err)
								}
							} else {
								commands["today"] = message
							}
						} else {
							err := irc.send_message(commands["today"])
							if err != nil {
								log.Fatal(err)
							}
						}
					}()
				}
			}
		}
		wg.Done()
	}()
	wg.Wait()
}
