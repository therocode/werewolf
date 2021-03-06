package ircgame

import (
	"log"

	"github.com/therocode/werewolf/irc"
	"github.com/therocode/werewolf/logic/components"
	"github.com/therocode/werewolf/logic/roles/seer"
	"github.com/therocode/werewolf/logic/roles/villager"
	"github.com/therocode/werewolf/logic/roles/werewolf"
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
//nolint: gocyclo
func (lobby *IrcLobby) HandleMessage(channel string, nick string, message string) {
	log.Printf("IRC [channel=%s, nick=%s, message=%s]", channel, nick, message)

	// Ignore any message not sent to the lobby channel or any active game channels or the bot
	if _, contains := lobby.games[channel]; !contains && channel != lobby.channel && channel != lobby.botname {
		return
	}

	if cmd, err := ParseCommand(channel, nick, message); err == nil {
		// Ignore commands sent directly to the bot
		if channel == lobby.botname {
			return
		}

		switch {
		case cmd.Command == "newgame":
			lobby.handleNewGame(cmd)
		case cmd.Command == "help" && channel == lobby.channel:
			lobby.message("Lobby commands:")
			lobby.message("!help - Displays this message.")
			lobby.message("!newgame #channel - Starts a new game in the specified channel. The bot must be an op in the game channel. Note that the # is required.")
			lobby.message("Game channel commands:")
			lobby.message("!join - Sign up to join the game before it begins.")
			lobby.message("!start - Start the game. Only the owner can use this command.")
			lobby.message("!cancel - Cancel a game before it starts. Only the owner can do this.")
		case channel == lobby.channel:
			lobby.message("%s is not a recognized command in the lobby channel. Join a game channel to run game-specific commands.", cmd.Command)
		case cmd.Command == "join":
			lobby.handleJoin(cmd)
		case cmd.Command == "start":
			lobby.handleStart(cmd)
		case cmd.Command == "cancel":
			lobby.handleCancel(cmd)

		default:
			if entry, contains := lobby.games[channel]; contains {
				entry.game.communication.SendToChannel("%s is not a recognized command in a game channel", cmd.Command)
			} else {
				log.Printf("Warning: Received command on channel that is not a game channel!")
			}
		}
	} else {
		// Handle player input
		if gameChannel, contains := lobby.gamePerPlayer[nick]; contains {
			lobby.games[gameChannel].game.communication.Respond(nick, message)
		}
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

	if entry, contains := lobby.games[cmd.Args[0]]; contains && !entry.game.IsFinished() {
		lobby.message("A game is already in progress in %s", cmd.Args[0])
		return
	}

	if cmd.Args[0] == lobby.channel {
		lobby.message("Cannot start a game in the lobby channel.")
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

func (lobby *IrcLobby) handleCancel(cmd Command) {
	game := lobby.games[cmd.Channel].game
	owner := lobby.games[cmd.Channel].owner

	if cmd.Nick != owner {
		game.communication.SendToChannel("Only the owner, %s, can cancel the game", owner)
		return
	}

	if game.IsRunning() {
		game.communication.SendToChannel("Cannot cancel a game that already started.")
		return
	}

	game.communication.SendToChannel("Game was cancelled by %s", cmd.Nick)
	game.EndGame()
	delete(lobby.games, cmd.Channel)
}

func newIrcGame(communication *irc.Irc) *IrcGame {
	game := NewIrcGame(communication)

	communication.SendToChannel("New game started!")

	lynch := components.NewLynch(game, communication)
	killVote := components.NewVote("kill")

	game.AddRole(villager.NewVillager(game, communication, lynch))
	game.AddRole(werewolf.NewWerewolf(game, communication, killVote, lynch))
	game.AddRole(seer.NewSeer(game, communication, lynch))

	return game
}
