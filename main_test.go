package main

import "testing"

func Test_if_gets_error_when_command_is_empty(t *testing.T) {
	irc, _ := new_irc()
	inputs := [][]string{{"", "some data"}, {"some data", ""}, {"", ""}}

	for _, input := range inputs {
        err := irc.send_command(input[0], input[1])
		if err == nil {
			t.Fatalf("Command or body shouldn't be empty!")
		}
	}
}

func Test_if_gets_error_when_message_is_empty(t *testing.T){
	irc, _ := new_irc()
	inputs := [][]string{{"", "some data"}, {"some data", ""}, {"", ""}}

	for _, input := range inputs {
        err := irc.send_message(input[0], input[1])
		if err == nil {
			t.Fatalf("Message shouldn't be empty!")
		}
	}
}
