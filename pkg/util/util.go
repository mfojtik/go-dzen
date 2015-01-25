package util

import (
	"os/exec"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	xrandrPath         = "/usr/sbin/xrandr"
	defaultScreenWidth = "1366"
)

func RunEvery(d time.Duration, f func() error) {
	timer := time.NewTicker(d)
	go func() {
		for {
			<-timer.C
			if err := f(); err != nil {
				glog.Errorf(err.Error())
			}
		}
	}()
}

func ScreenWidth() string {
	out, err := exec.Command(xrandrPath, "-q").Output()
	if err != nil {
		glog.Fatalf("Unable to get screen width: %v", err)
	}
	currentLine := strings.SplitN(string(out), "\n", 2)[0]
	parts := strings.Split(currentLine, " ")
	for i, p := range parts {
		if strings.TrimSpace(p) == "current" {
			return strings.TrimSpace(parts[i+1])
		}
	}
	return defaultScreenWidth
}
