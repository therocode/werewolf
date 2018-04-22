package ircgame

import (
	"github.com/therocode/werewolf/werewolf"
	"github.com/therocode/werewolf/werewolf/irc"
	"github.com/therocode/werewolf/werewolf/logic/components"
	"github.com/therocode/werewolf/werewolf/logic/roles"
	ircevent "github.com/thoj/go-ircevent"
)

type gameEntry struct {
	owner         string
	communication *irc.Irc
	game          *IrcGame
}

// IrcLobby keeps track of all ongoing IRC games
type IrcLobby struct {
	channel string
	irccon  *ircevent.Connection
	games   map[string]gameEntry
}

// NewIrcLobby creates a new IrcLobby instance
func NewIrcLobby(channel string, irccon *ircevent.Connection) *IrcLobby {
	lobby := &IrcLobby{}
	lobby.channel = channel
	lobby.irccon = irccon
	lobby.games = map[string]gameEntry{}
	return lobby
}

func (lobby *IrcLobby) message(format string, params ...interface{}) {
	lobby.irccon.Privmsgf(lobby.channel, format, params...)
}

// HandleMessage handles all incoming IRC messages
func (lobby *IrcLobby) HandleMessage(channel string, nick string, message string) {
	// Ignore any message not sent to the lobby channel or any active game channels
	if _, contains := lobby.games[channel]; !contains && channel != lobby.channel {
		return
	}

	if cmd, err := werewolf.ParseCommand(channel, nick, message); err == nil {
		if cmd.Command == "newgame" {
			if len(cmd.Args) != 1 {
				lobby.message("Usage: !%s [channel]", cmd.Command)
				return
			}

			if cmd.Args[0][0] != '#' {
				lobby.message("Channel names must start with '#")
				return
			}

			if gameEntry, contains := lobby.games[cmd.Args[0]]; contains && !gameEntry.game.IsFinished() {
				lobby.message("A game is already in progress in %s", cmd.Args[0])
				return
			}

			lobby.message("Starting a new game in %s", cmd.Args[0])
			communication := irc.NewIrc(lobby.irccon, cmd.Args[0])
			lobby.games[cmd.Args[0]] = gameEntry{nick, communication, newIrcGame(communication)}
			lobby.irccon.Join(cmd.Args[0])
		} else if channel == lobby.channel {
			lobby.message("%s is not a recognized command in the lobby channel. Join a game channel to run game-specific commands.", cmd.Command)
		} else {
			game := lobby.games[channel].game
			owner := lobby.games[channel].owner

			switch cmd.Command {
			case "join":
				game.AddPlayer(nick)
			case "start":
				if nick != owner {
					game.communication.SendToChannel("Only the owner, %s, can start the game", owner)
					return
				}

				if len(game.players) < 4 {
					game.communication.SendToChannel("Cannot start game with fewer than 4 players!")
					return
				}

				go game.Run()
			}
		}
	} else if game, contains := lobby.games[channel]; contains {
		// Handle player input
		game.communication.Respond(nick, message)
	}
}

func newIrcGame(communication *irc.Irc) *IrcGame {
	game := NewIrcGame(communication)

	communication.SendToChannel("New game started!")

	lynchVote := components.NewVote("lynch")
	killVote := components.NewVote("kill")

	game.AddRole(roles.NewVillager(game, communication, lynchVote))
	game.AddRole(roles.NewWerewolf(game, communication, killVote, lynchVote))

	return game
}
