package dzen

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/mfojtik/go-dzen/pkg/plugins"
)

const (
	Path         = "/usr/sbin/dzen2"
	XrandrPath   = "/usr/sbin/xrandr"
	DefaultWidth = "1366"
)

type Dzen struct {
	Args    []string
	Plugins []plugins.Plugin
}

func getCurrentWidth() string {
	out, err := exec.Command(XrandrPath, "-q").Output()
	if err != nil {
		glog.Fatalf("Unable to get screen width: %v", err)
	}
	currentLine := strings.SplitN(string(out), "\n", 2)[0]
	parts := strings.Split(currentLine, " ")
	for i, p := range parts {
		if strings.TrimSpace(p) == "current" {
			glog.Infof("Discovered screen width: %spx", strings.TrimSpace(parts[i+1]))
			return strings.TrimSpace(parts[i+1])
		}
	}
	return DefaultWidth
}

func NewCommand() *Dzen {
	return &Dzen{
		Args: []string{
			"-fn", "Terminus:size=8",
			"-fg", "#aaaaaa",
			"-bg", "#111111",
			"-ta", "right",
			"-x", "0",
			"-y", "0",
			"-w", getCurrentWidth(),
			"-h", "16",
		},
	}
}

func (d *Dzen) collectPluginStreams() chan string {
	out := make(chan string)
	pluginChannels := make([]chan string, len(d.Plugins))

	pluginNames := []string{}
	for i, p := range d.Plugins {
		pluginChannels[i] = p.Stream()
		pluginNames = append(pluginNames, p.String())
	}

	glog.Infof("Enabled plugins: [%v]", strings.Join(pluginNames, ","))

	// Listen to updates from streams
	go func() {
		for {
			result := ""
			for _, c := range pluginChannels {
				result += fmt.Sprintf("%s%s", plugins.Separator, <-c)
			}
			//glog.Info("Sending inputs to dzen...")
			out <- result + "\n"
		}
	}()
	return out
}

func streamToInput(stream <-chan string, stdin io.WriteCloser) {
	for {
		if s := <-stream; s != "exit" {
			io.Copy(stdin, bytes.NewBufferString(s))
			continue
		}
		break
	}
}

func (d *Dzen) Run() {
	cmd := exec.Command(Path, d.Args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		glog.Fatalf(err.Error())
	}

	w := sync.WaitGroup{}
	w.Add(1)
	go func() {
		defer func() {
			stdin.Close()
			w.Done()
		}()
		streamToInput(d.collectPluginStreams(), stdin)
	}()

	if err := cmd.Start(); err != nil {
		glog.Fatalf(err.Error())
	}

	// Wait for stdin to close, dzen2 will shut down automatically as there will
	// be no stdin
	w.Wait()
	if err := cmd.Wait(); err != nil {
		glog.Fatalf(err.Error())
	}
}
