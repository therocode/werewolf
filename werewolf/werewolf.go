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
	cmdJoin  = "join"
	cmdStart = "start"
)

// gamestates
const ( // iota is reset to 0
	GameStateInvite = "GameStateInvite" // c0 == 0
	GameStateDay    = "GameStateDay"    // c1 == 1
	GameStateNight  = "GameStateNight"  // c2 == 2
)

type Player struct {
	irc  *irc.Connection
	name string
	role Role
}

type Werewolf struct {
	irc          *irc.Connection
	config       Config
	mainChannel  IRCChannel
	state        string
	owner        string
	participants map[string]*Player
}

// Create new game instance based on a configuration
func NewWerewolf(irc *irc.Connection, config Config, mainChannel string) (instance *Werewolf) {
	instance = &Werewolf{irc: irc, config: config}
	instance.mainChannel = newIRCChannel(irc, mainChannel)
	instance.participants = make(map[string]*Player)
	instance.mainChannel.message("New werewolf game started.")
	return
}

// Parse message into a game command with arguments
func (instance *Werewolf) HandleMessage(channel string, nick string, cmd Command) {
	log.Printf("ch[%s], nick[%s]| command: '%s' [%s]", channel, nick, cmd.Command, strings.Join(cmd.Args, ","))
	instance.handleCommand(channel, nick, cmd)
}

func (instance *Werewolf) logInvalidCommand(channel string, command string) {
	log.Printf(channel, "command '%s' is either unknown or not applicable during %s", command, instance.state)
	instance.irc.Privmsgf(channel, "command '%s' is either unknown or not applicable during %s", command, instance.state)
}

func (instance *Werewolf) handleCommand(channel string, nick string, cmd Command) {
	switch instance.state {
	case GameStateInvite:
		switch cmd.Command {
		case cmdJoin:
			if instance.getPlayer(nick) == nil {
				instance.irc.Privmsgf(channel, "%s has joined the game!", nick)
				instance.playerJoin(nick)
				instance.participants[instance.owner].message("you are the leader of the pack")
			} else {
				instance.irc.Privmsgf(channel, "You cannot join, %s. You've already joined.", nick)
			}
		case cmdStart:
			if nick != instance.owner {
				instance.irc.Privmsgf(channel, "%s is the owner and only they can start the game.", instance.owner)
			} else {
				instance.irc.Privmsgf(channel, "The game starts now!")
				instance.mainChannel.message("the game has now been started!")
				instance.startGame()
			}
		default:
			instance.logInvalidCommand(channel, cmd.Command)
		}
	default:
		instance.logInvalidCommand(channel, cmd.Command)
	}
}
