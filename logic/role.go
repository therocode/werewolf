package logic

import (
	"github.com/therocode/werewolf/logic/timeline"
)

// Role is the interface for any game role.
type Role interface {
	// Name of the role
	Name() string

	// Implements the timeline.Generator interface
	Generate() []timeline.Event

	// Handle game event for a specific player. Use hasTerminated to inform control thread that event handler has terminated.
	// The general rule is that the data mutex should be locked when entering the Handle call and unlocked when exiting from it,
	// and also, the mutex should be unlocked whenever a blocking call is made, and relocked afterward.
	Handle(player string, event timeline.Event, hasTerminated chan bool)

	// HasComponent checks if the role has the specified component
	HasComponent(Component) bool
}
