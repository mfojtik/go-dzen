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

func (p *SimpleDate) Stream() chan string {
	out := make(chan string, 1)
	timer := time.NewTicker(time.Duration(1 * time.Second))
	go func() {
		for now := range timer.C {
			out <- icon(clockIcon) + now.Format(defaultFormat)
		}
	}()
	return out
}
