package main

import "testing"

func Test_parser(t *testing.T) {
	input := ":hicaro____!hicaro____@hicaro____.tmi.twitch.tv PRIVMSG #hicaro____ :message"
	parser := new_parser(input)
	username, message, err := parser.parse()

	if err != nil || username != "hicaro____" || message != "message" {
		t.Fatalf("username should be 'hicaro____', not '%s'. Message should be 'message', not '%s'", username, message)
	}
}

func Test_get_command_message_body(t *testing.T) {
	inputs := [][]string{{":hicaro____!hicaro____@hicaro____.tmi.twitch.tv PRIVMSG #hicaro____ :!bot message body here", "message body here"},
		{":hicaro____!hicaro____@hicaro____.tmi.twitch.tv PRIVMSG #hicaro____ :!bot", ""},
		{":hicaro____!hicaro____@hicaro____.tmi.twitch.tv PRIVMSG #hicaro____ :!bot       ", ""}}

	for _, message := range inputs {
		unparsed_message, expected_message_body := message[0], message[1]

		parser := new_parser(unparsed_message)
		_, parsed_message, _ := parser.parse()
		message_body := parser.get_command_message_body(parsed_message, "!bot")

		if message_body != expected_message_body {
			t.Errorf("message body should be '%s', not '%s'", expected_message_body, message_body)
		}
	}
}
