package main

import (
    "testing"
    "github.com/joho/godotenv"
    "os"
)

func Test_if_gets_error_when_command_is_empty(t *testing.T) {
    err := godotenv.Load()
    if err != nil {
		t.Fatalf("Unable to read .env file")
    }   

	var (
		OAUTH_TOKEN  = os.Getenv("OAUTH_TOKEN")
		BOT_NAME     = os.Getenv("BOT_NAME")
		CHANNEL_NAME = "#" + os.Getenv("CHANNEL_NAME")
	)

	irc, err := new_irc(CHANNEL_NAME, BOT_NAME, OAUTH_TOKEN)
    if err != nil {
		t.Fatalf("Unable to create IRC client")
    }

	inputs := [][]string{{"", "some data"}, {"some data", ""}, {"", ""}}

	for _, input := range inputs {
		err := irc.send_command(input[0], input[1])
		if err == nil {
			t.Fatalf("command or body shouldn't be empty!")
		}
	}
}

func Test_if_gets_error_when_message_is_empty(t *testing.T) {
    var err error

    err = godotenv.Load()
    if err != nil {
		t.Fatalf("Unable to read .env file")
    }   

	var (
		OAUTH_TOKEN  = os.Getenv("OAUTH_TOKEN")
		BOT_NAME     = os.Getenv("BOT_NAME")
		CHANNEL_NAME = "#" + os.Getenv("CHANNEL_NAME")
	)

	irc, err := new_irc(CHANNEL_NAME, BOT_NAME, OAUTH_TOKEN)
    if err != nil {
		t.Fatalf("Unable to create IRC client")
    }

    err = irc.send_message("")
    if err == nil {
        t.Fatalf("Message shouldn't be empty!")
    }
}
