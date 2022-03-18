package main

import (
	"fmt"
	"strings"
)

type Parser struct {
	raw_message string
}

func new_parser(raw_message string) *Parser {
	return &Parser{raw_message}
}

func (parser *Parser) is_user_message() bool {
	return strings.Contains(parser.raw_message, "PRIVMSG")
}

// TODO(Hícaro): Fix and refactor parser function(a lot of indentation)
func (parser *Parser) get_message() (string, string, error) {
	nickname, message := "", ""
	if parser.is_user_message() {
		exclamation_mark_index := strings.Index(parser.raw_message, "!")
		if exclamation_mark_index != -1 {
			nickname = parser.raw_message[1:exclamation_mark_index]
			remainder_message := parser.raw_message[exclamation_mark_index+1:]
			hashtag_index := strings.Index(remainder_message, "#")

			if hashtag_index != -1 {
				remainder_message = remainder_message[hashtag_index+1:]
				colon_index := strings.Index(parser.raw_message, ":")
				message = remainder_message[colon_index:]
			} else {
				return "", "", fmt.Errorf("Hashtag not found!")
			}
		} else {
			return "", "", fmt.Errorf("Exclamation mark not found!")
		}
	}
	return nickname, message, nil
}
