package plugins

import "time"

const (
	defaultFormat = "2006-01-02 15:04"
	clockIcon     = ".icons/dzen2/clock.xbm"
)

type SimpleDate struct {
	Format string
}

func (n SimpleDate) String() string {
	return "date"
}

func (p *SimpleDate) Send(output chan UpdateMessage) {
	timer := time.NewTicker(time.Duration(1 * time.Second))
	go func() {
		for now := range timer.C {
			output <- UpdateMessage{Src: p, Message: icon(clockIcon) + now.Format(defaultFormat)}
		}
	}()
}
