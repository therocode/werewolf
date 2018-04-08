package werewolf

import (
	"log"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

type Config struct {
	dayLength   int `json:"day_length"`   //in minutes
	nightLength int `json:"night_length"` //in minutes
}

// commands
const (
	cmdJoin = "join"
)

// gamestates
const ( // iota is reset to 0
	GameStateInvite = iota // c0 == 0
	GameStateDay    = iota // c1 == 1
	GameStateNight  = iota // c2 == 2
)

type Player struct {
}

type Werewolf struct {
	irc          *irc.Connection
	config       Config
	state        int
	participants map[string]Player
}

// Create new game instance based on a configuration
func NewWerewolf(irc *irc.Connection, config Config) (instance *Werewolf) {
	instance = &Werewolf{irc: irc, config: config}
	instance.state = GameStateInvite
	instance.participants = make(map[string]Player)
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
	if instance.state == GameStateInvite {
		switch command {
		case cmdJoin:
			if instance.getPlayer(nick) == nil {
				instance.irc.Privmsgf(channel, "%s has joined the game!", nick)
				instance.playerJoin(nick)
			} else {
				instance.irc.Privmsgf(channel, "You cannot join, %s. You've already joined.", nick)
			}
		default:
			log.Printf(channel, "unknown command '%s'", command)
			instance.irc.Privmsgf(channel, "Unknown command '%s'", command)
		}
	}
}
