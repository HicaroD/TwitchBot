package main

import "testing"

func Test_parser(t *testing.T) {
	input := ":hicaro____!hicaro____@hicaro____.tmi.twitch.tv PRIVMSG #hicaro____ :message"
	parser := new_parser(input)
	username, message, err := parser.parse()

	if err != nil || username != "hicaro____" || message != "message" {
		t.Fatalf("Username should be 'hicaro____', not '%s'. Message should be 'message', not '%s'", username, message)
	}
}
