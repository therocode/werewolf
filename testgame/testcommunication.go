package testgame

import "fmt"

type TestCommunication struct{}

func NewTestCommunication() *TestCommunication {
	return &TestCommunication{}
}

func (instance *TestCommunication) SendToChannel(format string, params ...interface{}) {
	if len(params) == 0 {
		fmt.Print(format + "\n")
	} else {
		fmt.Printf(format+"\n", params...)
	}
}

func (instance *TestCommunication) SendToPlayer(player string, format string, params ...interface{}) {
	fmt.Printf("PM for %s: ", player)
	if len(params) == 0 {
		fmt.Print(format + "\n")
	} else {
		fmt.Printf(format+"\n", params...)
	}
}

func (instance *TestCommunication) Request(requestFrom string, promptFormat string, params ...interface{}) (string, bool) {
	fmt.Printf(promptFormat+"\n", params...)
	var text string
	_, err := fmt.Scanln(&text)

	return text, err != nil
}

func (*TestCommunication) Respond(string, string) {}

func (*TestCommunication) MutePlayer(string) {}

func (*TestCommunication) UnmutePlayer(string) {}

func (*TestCommunication) MuteChannel() {}

func (*TestCommunication) Leave() {}
