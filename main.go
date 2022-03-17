package main

import (
    "fmt"
    "log"
    "os"
    "net"
    "strings"
    "sync"
    "github.com/joho/godotenv"
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
    fmt.Fprintf(irc.client, "%s %s\n", command, body)
    return nil
}

const BUFFER_SIZE = 2040;

func main(){
    var err error

    err = godotenv.Load()
    if err != nil {
        log.Fatal("Error loading the .env file")
    }

    var (
        OAUTH_TOKEN = os.Getenv("OAUTH_TOKEN")
        BOT_NAME = os.Getenv("BOT_NAME")
        CHANNEL_NAME = os.Getenv("CHANNEL_NAME")
    )
    fmt.Println(OAUTH_TOKEN, BOT_NAME, CHANNEL_NAME)

    connection, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")

    defer connection.Close()

    if err != nil {
        log.Fatal("Unable to establish connection to IRC server")
    }

    fmt.Fprintf(connection, "PASS %s\nNICK %s\nJOIN %s\n", OAUTH_TOKEN, BOT_NAME, "#" + CHANNEL_NAME)
    for {
        received_data := make([]byte, BUFFER_SIZE)
        received_data_size, err := connection.Read(received_data)
        if err != nil {
            log.Fatal("Unable to read data from socket")
        }

        message := string(received_data)

        if received_data_size > 0 {
            fmt.Println(message)
            if strings.HasPrefix(message, "PING") {
                fmt.Fprintf(connection, "PONG :tmi.twitch.tv")
            }
        } 
    }
}
