package plugins

import (
	"os/exec"
	"strings"
	"time"

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

func (b *Bspwm) Stream() chan string {
	out := make(chan string, 1)
	updates := subscribeUpdates()
	go func() {
		lastStatus := ""
		for {
			select {
			case update := <-updates:
				lastStatus = parseStatus(update)
				out <- lastStatus
				glog.Infof("[bspwm] received %s", update)
			default:
				out <- lastStatus
				time.Sleep(time.Duration(50 * time.Millisecond))
			}
		}
	}()
	return out
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
