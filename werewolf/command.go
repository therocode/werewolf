package werewolf

import (
	"errors"
	"strings"
)

type Command struct {
	Command string
	Args    []string
}

func ParseCommand(raw string) (command Command, err error) {
	if raw[0] != '!' {
		return Command{}, errors.New("Command strings must start with '!'")
	}

	tokens := strings.Fields(raw)
	command = Command{
		Command: tokens[0][1:], // Remove leading '!' from command
		Args:    tokens[1:],    // Arguments are all words following the command
	}

	if len(command.Command) == 0 {
		return Command{}, errors.New("Command cannot be empty")
	}

	return command, nil
}
