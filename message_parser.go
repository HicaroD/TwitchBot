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

func (parser *Parser) get_nickname() (string, error) {
	exclamation_mark_index := strings.Index(parser.raw_message, "!")
	if exclamation_mark_index == -1 {
		return "", fmt.Errorf("Exclamation mark not found!")
	}
	nickname := parser.raw_message[1:exclamation_mark_index]
	return nickname, nil
}

func (parser *Parser) get_message() (string, error) {
	hashtag_index := strings.Index(parser.raw_message, "#")
	if hashtag_index == -1 {
		return "", fmt.Errorf("Hashtag not found!")
	}
	reimainder := parser.raw_message[hashtag_index+1:]
	colon_index := strings.Index(reimainder, ":")

	if colon_index == -1 {
		return "", fmt.Errorf("Hashtag not found!")
	}
	message := reimainder[colon_index+1:]
	return message, nil
}

func (parser *Parser) parse() (string, string, error) {
	if parser.is_user_message() {
		nickname, err := parser.get_nickname()
		if err != nil {
			return "", "", err
		}

		message, err := parser.get_message()
		if err != nil {
			return "", "", err
		}

		return nickname, message, nil
	}
	return "", "", nil
}
