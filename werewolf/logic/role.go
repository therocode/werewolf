package logic

import (
	"github.com/therocode/werewolf/werewolf/logic/timeline"
)

// Role is the interface for any game role.
type Role interface {
	// Name of the role
	Name() string

	// Implements the timeline.Generator interface
	Generate() []timeline.Event

	// Handle game event for a specific player. Use hasTerminated to inform control thread that event handler has terminated.
	Handle(player string, event timeline.Event, hasTerminated chan bool)

	// HasComponent checks if the role has the specified component
	HasComponent(Component) bool
}
