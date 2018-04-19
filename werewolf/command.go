package werewolf

import (
	"errors"
	"strings"
)

type command struct {
	Channel string
	Nick    string
	Command string
	Args    []string
}

func ParseCommand(channel string, nick string, raw string) (cmd command, err error) {
	if raw[0] != '!' {
		return command{}, errors.New("Command strings must start with '!'")
	}

	tokens := strings.Fields(raw)
	cmd = command{
		Command: tokens[0][1:], // Remove leading '!' from command
		Args:    tokens[1:],    // Arguments are all words following the command
		Nick:    nick,
		Channel: channel,
	}

	if len(cmd.Command) == 0 {
		return command{}, errors.New("Command cannot be empty")
	}

	return cmd, nil
}
