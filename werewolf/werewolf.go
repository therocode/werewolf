package werewolf

import (
	"log"
	"strings"
)

type Config struct {
	dayLength   int `json:"day_length"`   //in minutes
	nightLength int `json:"night_length"` //in minutes
}

type Werewolf struct {
	config Config
}

// Create new game instance based on a configuration
func NewWerewolf(config Config) (instance *Werewolf) {
	instance = &Werewolf{config}
	return
}

// Parse message into a game command with arguments
func (instance Werewolf) HandleMessage(channel string, nick string, message string) {
	words := strings.Fields(message)

	if len(words) > 0 {
		command := words[0]   //command is f.e "vote" from "!vote <name>"
		command = command[1:] //remove the leading '!'

		if command == "" {
			log.Printf("ch[%s], nick[%s]| empty command", channel, nick)
			return
		}

		command_args := []string{} //commands is a list of strings with arguments for the command. f.e. [<name>] from "!vote <name>". can be empty
		if len(words) > 1 {
			command_args = words[1:]
		}

		log.Printf("ch[%s], nick[%s]| command: '%s' [%s]", channel, nick, command, strings.Join(command_args, ","))
		instance.handleCommand(channel, nick, command, command_args)
	} else {
		log.Printf("ch[%s], nick[%s]| empty message", channel, nick)
	}
}

func (instance Werewolf) handleCommand(channel string, nick string, command string, arguments []string) {
}
