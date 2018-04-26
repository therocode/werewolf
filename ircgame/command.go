package ircgame

import (
	"errors"
	"strings"
)

// Command represents a command to control the Werewolf game
type Command struct {
	Channel string
	Nick    string
	Command string
	Args    []string
}

// ParseCommand parses a raw string into a Command.
func ParseCommand(channel string, nick string, raw string) (cmd Command, err error) {
	if raw[0] != '!' {
		return Command{}, errors.New("Command strings must start with '!'")
	}

	tokens := strings.Fields(raw)
	cmd = Command{
		Command: tokens[0][1:], // Remove leading '!' from command
		Args:    tokens[1:],    // Arguments are all words following the command
		Nick:    nick,
		Channel: channel,
	}

	if len(cmd.Command) == 0 {
		return Command{}, errors.New("Command cannot be empty")
	}

	return cmd, nil
}
