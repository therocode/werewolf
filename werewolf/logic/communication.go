package logic

type Communication interface {
	SendToChannel(format string, params ...interface{})
	SendToPlayer(player string, format string, params ...interface{})
	Request(requestFrom string, promptFormat string, params ...interface{}) string
}
