package ircgame

import (
	"log"

	"github.com/therocode/werewolf/irc"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/roles"
	ircevent "github.com/thoj/go-ircevent"
)

type gameEntry struct {
	owner         string
	communication *irc.Irc
	game          *IrcGame
}

// IrcLobby keeps track of all ongoing IRC games
type IrcLobby struct {
	botname       string
	channel       string
	irccon        *ircevent.Connection
	games         map[string]gameEntry
	gamePerPlayer map[string]string
}

// NewIrcLobby creates a new IrcLobby instance
func NewIrcLobby(botname string, channel string, irccon *ircevent.Connection) *IrcLobby {
	lobby := &IrcLobby{}
	lobby.botname = botname
	lobby.channel = channel
	lobby.irccon = irccon
	lobby.games = map[string]gameEntry{}
	lobby.gamePerPlayer = map[string]string{}
	return lobby
}

func (lobby *IrcLobby) message(format string, params ...interface{}) {
	lobby.irccon.Privmsgf(lobby.channel, format, params...)
}

// HandleMessage handles all incoming IRC messages
func (lobby *IrcLobby) HandleMessage(channel string, nick string, message string) {
	log.Printf("IRC [channel=%s, nick=%s, message=%s]", channel, nick, message)

	// Ignore any message not sent to the lobby channel or any active game channels or the bot
	if _, contains := lobby.games[channel]; !contains && channel != lobby.channel && channel != lobby.botname {
		return
	}

	if cmd, err := ParseCommand(channel, nick, message); err == nil {
		switch {
		case cmd.Command == "newgame":
			lobby.handleNewGame(cmd)
		case channel == lobby.channel:
			lobby.message("%s is not a recognized command in the lobby channel. Join a game channel to run game-specific commands.", cmd.Command)
		case cmd.Command == "join":
			lobby.handleJoin(cmd)
		case cmd.Command == "start":
			lobby.handleStart(cmd)
		default:
			lobby.games[channel].game.communication.SendToChannel("%s is not a recognized command in a game channel", cmd.Command)
		}
	} else {
		// Handle player input
		gameChannel := lobby.gamePerPlayer[nick]
		lobby.games[gameChannel].game.communication.Respond(nick, message)
	}
}

func (lobby *IrcLobby) handleNewGame(cmd Command) {
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
	lobby.games[cmd.Args[0]] = gameEntry{cmd.Nick, communication, newIrcGame(communication)}
	lobby.irccon.Join(cmd.Args[0])
}

func (lobby *IrcLobby) handleJoin(cmd Command) {
	game := lobby.games[cmd.Channel].game

	if otherChannel, contains := lobby.gamePerPlayer[cmd.Nick]; contains {
		otherGame := lobby.games[otherChannel].game

		if !otherGame.IsFinished() {
			game.communication.SendToChannel("You can only join one game, and %s still has a game in progress", otherChannel)
			return
		}
	}

	game.AddPlayer(cmd.Nick)
	lobby.gamePerPlayer[cmd.Nick] = cmd.Channel
}

func (lobby *IrcLobby) handleStart(cmd Command) {
	game := lobby.games[cmd.Channel].game
	owner := lobby.games[cmd.Channel].owner

	if cmd.Nick != owner {
		game.communication.SendToChannel("Only the owner, %s, can start the game", owner)
		return
	}

	if len(game.players) < 4 {
		game.communication.SendToChannel("Cannot start game with fewer than 4 players!")
		return
	}

	go game.Run()
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
