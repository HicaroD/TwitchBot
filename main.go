package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	twitch_irc_client = "irc.chat.twitch.tv:6667"
	wg                = sync.WaitGroup{}
)

type IRC struct {
	channel_name string
	bot_name     string
	oauth_token  string
	client       net.Conn
}

type Commands struct {
	Commands string `json:"commands"`
	Me       string `json:"me"`
	Socials  string `json:"socials"`
	Projects string `json:"projects"`
	Bot      string `json:"bot"`
	Today    string `json:"today"`
}

func get_commands() (*Commands, error) {
	var err error

	content, err := ioutil.ReadFile("./commands.json")
	if err != nil {
		return nil, err
	}

	var commands Commands
	err = json.Unmarshal(content, &commands)

	if err != nil {
		return nil, err
	}

	return &commands, nil
}

func new_irc(channel_name, bot_name, oauth_token string) (*IRC, error) {
	connection, err := net.Dial("tcp", twitch_irc_client)
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
	godotenv.Load()

	var err error
	var (
		OAUTH_TOKEN  = os.Getenv("OAUTH_TOKEN")
		BOT_NAME     = os.Getenv("BOT_NAME")
		CHANNEL_NAME = "#" + os.Getenv("CHANNEL_NAME")
	)
	fmt.Println("Joining chat!")

	irc, err := new_irc(CHANNEL_NAME, BOT_NAME, OAUTH_TOKEN)

	if err != nil {
		log.Fatal(err)
	}

	commands, err := get_commands()
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
						err := irc.send_message(commands.Commands)
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!bot") {
					go func() {
						err := irc.send_message(commands.Bot)
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!me") {
					go func() {
						err := irc.send_message(commands.Me)
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!projects") {
					go func() {
						err := irc.send_message(commands.Projects)
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!socials") {
					go func() {
						err := irc.send_message(commands.Socials)
						if err != nil {
							log.Fatal(err)
						}
					}()
				} else if strings.HasPrefix(parsed_message, "!today") {
					go func() {
						if nickname == "hicaro____" {
							message_body := parser.get_command_message_body(parsed_message, "!today")

							if message_body == "" {
								err := irc.send_message(commands.Today)
								if err != nil {
									log.Fatal(err)
								}
							} else {
								commands.Today = message_body
							}
						} else {
							err := irc.send_message(commands.Today)
							if err != nil {
								log.Fatal(err)
							}
						}
					}()
				}
			}
		}
	}()
	wg.Wait()
}
