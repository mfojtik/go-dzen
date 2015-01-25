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
	Path = "/usr/sbin/dzen2"
)

var defaultArgs = []string{
	"-fn", "Terminus:size=8",
	"-fg", "#aaaaaa",
	"-bg", "#111111",
	"-y", "0",
	"-h", "16",
}

type Bar struct {
	ID      string
	Plugins []plugins.Plugin
	Updates chan plugins.UpdateMessage
	Output  chan string
	Args    []string
}

func NewBar(x, width int, align string) *Bar {
	return &Bar{
		ID: align,
		Args: append(defaultArgs, []string{
			"-x", fmt.Sprintf("%d", x),
			"-w", fmt.Sprintf("%d", width),
			"-ta", align,
		}...),
		Output: make(chan string),
	}
}

func (b *Bar) Add(p plugins.Plugin) {
	b.Plugins = append(b.Plugins, p)
}

func (b *Bar) outputToStdin(stdin io.WriteCloser) {
	for {
		if s := <-b.Output; len(s) > 0 {
			io.Copy(stdin, bytes.NewBufferString(s))
			continue
		}
		break
	}
}

func (b *Bar) ReceivePluginUpdates() {
	pluginNames := make([]string, len(b.Plugins))
	sources := make(map[string]string, len(b.Plugins))
	b.Updates = make(chan plugins.UpdateMessage, len(b.Plugins))

	// Tell plugins where they should send the updates
	for i, p := range b.Plugins {
		p.Send(b.Updates)
		pluginNames[i] = p.String()
	}

	glog.Infof("Bar[%s] Plugins: %s", b.ID, strings.Join(pluginNames, ","))

	// Wait 'till all plugins send an update, then craft the final result line
	for {
		update := <-b.Updates
		sources[update.Src.String()] = update.Message

		waitForUpdates := false
		for name, _ := range sources {
			if len(sources[name]) == 0 {
				glog.Infof("error!")
				waitForUpdates = true
				continue
			}
		}
		if waitForUpdates {
			glog.Infof("Bar[%s] Waiting for plugins to send updates", b.ID)
			continue
		}

		result := []string{}
		for _, p := range b.Plugins {
			result = append(result, sources[p.String()])
		}
		b.Output <- strings.Join(result, plugins.Separator) + "\n"
	}
}

func (b *Bar) Start() {
	cmd := exec.Command(Path, b.Args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		glog.Fatalf(err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer func() {
			stdin.Close()
			wg.Done()
		}()
		b.outputToStdin(stdin)
	}()

	if err := cmd.Start(); err != nil {
		glog.Fatalf(err.Error())
	}

	go func() { b.ReceivePluginUpdates(); wg.Done() }()

	// Wait for stdin to close, dzen2 will shut down automatically as there will
	// be no stdin
	wg.Wait()
	if err := cmd.Wait(); err != nil {
		glog.Fatalf(err.Error())
	}
}
