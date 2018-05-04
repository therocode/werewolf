package logic

// Communication method for the game logic.
type Communication interface {
	SendToChannel(format string, params ...interface{})
	SendToPlayer(player string, format string, params ...interface{})

	// Request asks for input from a specific player and blocks until the input is
	// delivered, or until a timeout occurs. The second return value is true if a timeout
	// occurred.
	Request(requestFrom string, promptFormat string, params ...interface{}) (string, bool)

	// Respond is used to deliver the input requested by Request()
	Respond(player string, message string)

	MutePlayer(player string)

	UnmutePlayer(player string)

	MuteChannel()

	Leave()
}
