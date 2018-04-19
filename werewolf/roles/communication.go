package roles

type Communication interface {
	SendToChannel(format string, params ...interface{})
	RequestName(requestFrom string, promptFormat string, params ...interface{}) string
}
