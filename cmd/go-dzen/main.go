package main

import (
	"flag"
	"os"

	"github.com/golang/glog"
	"github.com/mfojtik/go-dzen/pkg/dzen"
	"github.com/mfojtik/go-dzen/pkg/plugins"
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
		&plugins.Network{Interfaces: []string{"wlp3s0", "enp0s25", "tun0"}},
		&plugins.Battery{},
		&plugins.SimpleDate{},
		&plugins.Bspwm{},
	}
}

func main() {
	dzenCmd := &cobra.Command{
		Use:  "go-dzen",
		Long: "Run dzen2",
		Run: func(cmd *cobra.Command, args []string) {
			d := dzen.NewCommand()
			d.Plugins = configurePlugins()
			d.Run()
		},
	}
	setupGlog(dzenCmd.PersistentFlags())
	err := dzenCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
