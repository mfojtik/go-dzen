package dzen

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/golang/glog"
	"github.com/mfojtik/go-dzen/pkg/plugins"
)

const Path = "/usr/sbin/dzen2"

type Dzen struct {
	Args    []string
	Plugins []plugins.Plugin
}

func NewCommand() *Dzen {
	return &Dzen{
		Args: []string{
			"-fn", "Liberation Mono:size=9",
			"-fg", "#aaaaaa",
			"-bg", "#111111",
			"-ta", "right",
			"-x", "0",
			"-y", "0",
			"-w", "1366",
			"-h", "18",
		},
	}
}

func (d *Dzen) collectPluginStreams() chan string {
	out := make(chan string)

	// Collect streams from all plugins
	glog.Infof("Registering %d plugins", len(d.Plugins))
	pluginChannels := make([]chan string, len(d.Plugins))
	for i, p := range d.Plugins {
		pluginChannels[i] = p.Stream()
	}

	glog.Infof("Plugins: %v", pluginChannels)

	// Listen to updates from streams
	go func() {
		for {
			result := ""
			for _, c := range pluginChannels {
				result = fmt.Sprintf("%s%s", result, <-c)
			}
			out <- result + "\n"
		}
	}()
	return out
}

func inputToStdin(inp <-chan string, stdin io.WriteCloser) {
	for {
		if s := <-inp; s != "exit" {
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
		inputToStdin(d.collectPluginStreams(), stdin)
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
