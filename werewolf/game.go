package werewolf

import (
	"github.com/therocode/werewolf/werewolf/roles"
	"github.com/therocode/werewolf/werewolf/timeline"
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
const (
	GameStateLobby    = "GameStateLobby"
	GameStateStarted  = "GameStateStarted"
	GameStateGameOver = "GameStateGameOver"
)

type Player struct {
	irc  *irc.Connection
	name string
	role roles.Role
}

type Game struct {
	irc         *irc.Connection
	config      Config
	mainChannel IRCChannel
	state       string
	timeline    []timeline.Event
	owner       string
	players     map[string]*Player
	gameRole    *roles.Base
}

/*
// Create new game instance based on a configuration
func NewWerewolfGame(irc *irc.Connection, config Config, mainChannel string, owner string) (instance *Game) {
	instance = &Game{irc: irc, config: config}
	instance.mainChannel = newIRCChannel(irc, mainChannel)
	instance.players = make(map[string]*Player)
	instance.mainChannel.message("New werewolf game started.")
	instance.state = GameStateLobby
	instance.owner = owner
	instance.gameRole = roles.NewBase(instance, instance)
	return
}

func (instance *Game) logInvalidCommand(cmd command) {
	log.Printf(cmd.Channel, "command '%s' is either unknown or not applicable during %s", cmd.Command, instance.state)
	instance.irc.Privmsgf(cmd.Channel, "command '%s' is either unknown or not applicable during %s", cmd.Command, instance.state)
}

func (instance *Game) HandleCommand(cmd command) {
	switch cmd.Command {
	case cmdJoin:
		instance.handleCommandJoin(cmd)
	case cmdStart:
		instance.handleCommandStart(cmd)
	default:
		instance.logInvalidCommand(cmd)
	}
}

func (instance *Game) handleCommandJoin(cmd command) {
	if instance.state != GameStateLobby {
		instance.logInvalidCommand(cmd)
		return
	}

	if instance.getPlayer(cmd.Nick) == nil {
		instance.irc.Privmsgf(cmd.Channel, "%s has joined the game!", cmd.Nick)
		instance.playerJoin(cmd.Nick)
		instance.players[instance.owner].message("you are the leader of the pack")
	} else {
		instance.irc.Privmsgf(cmd.Channel, "You cannot join, %s. You've already joined.", cmd.Nick)
	}
}

func (instance *Game) handleCommandStart(cmd command) {
	if instance.state != GameStateLobby {
		instance.logInvalidCommand(cmd)
		return
	}

	if cmd.Nick != instance.owner {
		instance.irc.Privmsgf(cmd.Channel, "%s is the owner and only they can start the game.", instance.owner)
	} else {
		instance.irc.Privmsgf(cmd.Channel, "The game starts now!")
		instance.mainChannel.message("the game has now been started!")
		instance.startGame()
	}
}

func (instance *Game) startGame() {
	instance.state = GameStateStarted
	//go instance.runGame()
}

func (instance *Game) getRoleSet() map[timeline.Generator]bool {
	roles := map[timeline.Generator]bool{}
	roles[instance.gameRole] = true
	for _, player := range instance.players {
		roles[player.role] = true
	}
	return roles
}

func (instance *Game) GetRoles() []roles.Role {
	roles := []roles.Role{}
		roles = append(roles, instance.gameRole)
		for _, player := range instance.players {
			roles = append(roles, player.role)
		}
	return roles
}

func (instance *Game) EndGame() {
	instance.state = GameStateGameOver
}

func (instance *Game) RequestName(requestFrom string) string {
	// TODO: Implement
	return ""
}

func (instance *Game) runGame() {
	for instance.state == GameStateStarted {
		// If there are no events in the timeline, generate more
		if len(instance.timeline) == 0 {
			instance.timeline = timeline.Generate(instance.getRoleSet())
		}

		// Pop the first event in the timeline
		var event timeline.Event
		event, instance.timeline = instance.timeline[0], instance.timeline[1:]

		// Create a goroutine executing the game role handler
		gameHasTerminated := make(chan bool, 1)
		go instance.gameRole.Handle("", instance, instance, event, gameHasTerminated)

		// Create a goroutine executing the role handler for the event for each player
		var hasTerminated map[string]chan bool
		for name, player := range instance.players {
			hasTerminated[name] = make(chan bool, 1)
			go player.role.Handle(name, instance, instance, event, hasTerminated[name])
		}

		// Block until all role handlers have finished
		<-gameHasTerminated
		for _, channel := range hasTerminated {
			<-channel
		}
	}
}

func (instance *Game) SendToChannel(format string, params ...interface{}) {
	instance.mainChannel.message(format, params)
}
*/
