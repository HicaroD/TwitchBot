package main

import "testing"

func test_if_gets_error_when_command_is_empty(t *testing.T) {
	var err error
	irc, _ := new_irc()
	inputs := [][]string{{"", "some data"}, {"some data", ""}, {"", ""}}

	for _, input := range inputs {
		err = irc.send_command(input[0], input[1])
		if err != nil {
			t.Errorf("It shouldn't be empty!")
		}

		err = irc.send_message(input[0], input[1])
		if err != nil {
			t.Errorf("It shouldn't be empty!")
		}
	}
}
