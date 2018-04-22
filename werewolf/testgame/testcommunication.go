package testgame

import "fmt"

// TestCommunication is a test implementation of the Communication interface which uses the command line
type TestCommunication struct{}

// NewTestCommunication creates a new TestCommunication instance
func NewTestCommunication() *TestCommunication {
	return &TestCommunication{}
}

// SendToChannel implements the Communication interface
func (instance *TestCommunication) SendToChannel(format string, params ...interface{}) {
	if len(params) == 0 {
		fmt.Print(format + "\n")
	} else {
		fmt.Printf(format+"\n", params...)
	}
}

// SendToPlayer implements the Communication interface
func (instance *TestCommunication) SendToPlayer(player string, format string, params ...interface{}) {
	fmt.Printf("PM for %s: ", player)
	if len(params) == 0 {
		fmt.Print(format + "\n")
	} else {
		fmt.Printf(format+"\n", params...)
	}
}

// RequestName implements the Communication interface
func (instance *TestCommunication) RequestName(requestFrom string, promptFormat string, params ...interface{}) string {
	fmt.Printf(promptFormat+"\n", params...)
	var text string
	fmt.Scanln(&text)

	return text
}
