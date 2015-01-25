package main

import (
	"flag"
	"os"
	"strconv"
	"sync"

	"github.com/golang/glog"
	"github.com/mfojtik/go-dzen/pkg/dzen"
	"github.com/mfojtik/go-dzen/pkg/plugins"
	"github.com/mfojtik/go-dzen/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func setupGlog(flags *pflag.FlagSet) {
	from := flag.CommandLine
	if fflag := from.Lookup("v"); fflag != nil {
		level := fflag.Value.(*glog.Level)
		levelPtr := (*int32)(level)
		flags.Int32Var(levelPtr, "loglevel", 0, "Set the level of log output (0-3)")
	}
	// FIXME currently glog has only option to redirect output to stderr
	// the preferred for STI would be to redirect to stdout
	flag.CommandLine.Set("logtostderr", "true")
}

func configurePlugins() []plugins.Plugin {
	return []plugins.Plugin{
		&plugins.Bspwm{},
		&plugins.Network{Interfaces: []string{"wlp3s0", "enp0s25", "tun0"}},
		&plugins.Battery{},
		&plugins.SimpleDate{},
	}
}

func main() {
	dzenCmd := &cobra.Command{
		Use:  "go-dzen",
		Long: "Run dzen2",
		Run: func(cmd *cobra.Command, args []string) {
			screenWidth, _ := strconv.Atoi(util.ScreenWidth())
			// Configure left bar
			leftBar := dzen.NewBar(0, screenWidth/2, "l")
			leftBar.Add(&plugins.Bspwm{})

			// Configure right bar
			rightBar := dzen.NewBar(screenWidth/2, screenWidth, "r")
			rightBar.Add(&plugins.Network{Interfaces: []string{"wlp3s0", "enp0s25", "tun0"}})
			rightBar.Add(&plugins.Battery{})
			rightBar.Add(&plugins.SimpleDate{})

			// Start the bars asynchronously
			var wg sync.WaitGroup
			wg.Add(2)
			go func() { leftBar.Start(); wg.Done() }()
			go func() { rightBar.Start(); wg.Done() }()
			wg.Wait()
		},
	}
	setupGlog(dzenCmd.PersistentFlags())
	err := dzenCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
