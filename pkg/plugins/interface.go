package plugins

type Plugin interface {
	Stream() chan string
	String() string
}
