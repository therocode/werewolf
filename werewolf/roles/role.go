package roles

import "github.com/therocode/werewolf/werewolf/timeline"

// Role is the interface for any game role.
type Role interface {
	Name() string
	Generate() []timeline.Event
	Handle(player string, event timeline.Event, hasTerminated chan bool)
}
