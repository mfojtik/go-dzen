package plugins

import (
	"os/exec"
	"strings"

	"github.com/golang/glog"
)

const bspcPath = "/usr/sbin/bspc"

type Bspwm struct{}

type ChannelWriter struct {
	Output chan string
}

func (c *ChannelWriter) Write(p []byte) (n int, err error) {
	c.Output <- strings.TrimSpace(string(p))
	return len(p), nil
}

func (b *Bspwm) String() string {
	return "bspwm"
}

func subscribeUpdates() chan string {
	writer := ChannelWriter{make(chan string)}
	cmd := exec.Command(bspcPath, "control", "--subscribe")
	cmd.Stdout = &writer
	go func() {
		cmd.Start()
		err := cmd.Wait()
		if err != nil {
			glog.Errorf("bspc command failed: %v", err)
		}
	}()
	return writer.Output
}

func (b *Bspwm) Send(out chan UpdateMessage) {
	updates := subscribeUpdates()
	go func() {
		for {
			out <- UpdateMessage{b, parseStatus(<-updates)}
		}
	}()
}

func parseStatus(s string) string {
	// WM1:o1/i:o1/ii:O1/iii:o1/iv:f1/v:LT
	parts := strings.Split(s, ":")
	result := ""
	for _, p := range parts {
		if strings.HasPrefix(p, "O") {
			result = strings.TrimPrefix(p, "O")
		}
	}
	return "[" + result + "]"
}
