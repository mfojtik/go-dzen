package plugins

type Plugin interface {
	Send(chan UpdateMessage)
	String() string
}

type UpdateMessage struct {
	Src     Plugin
	Message string
}
