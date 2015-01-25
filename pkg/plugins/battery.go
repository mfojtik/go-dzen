package plugins

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/mfojtik/go-dzen/pkg/util"
)

const (
	BatteryDevice    = "/sys/class/power_supply/BAT0/capacity"
	ChargerDevice    = "/sys/class/power_supply/AC/online"
	ChargerIcon      = ".icons/dzen2/ac.xbm"
	EmptyBatteryIcon = ".icons/dzen2/bat_empty_01.xbm"
	LowBatteryIcon   = ".icons/dzen2/bat_low_01.xbm"
	FullBatteryIcon  = ".icons/dzen2/bat_full_01.xbm"
)

type Battery struct{}

func (n Battery) String() string { return "battery" }

func (n *Battery) Send(out chan UpdateMessage) {
	util.RunEvery(time.Duration(1*time.Second), func() error {
		out <- UpdateMessage{n, n.prettyValue()}
		return nil
	})
}

func (n *Battery) IsUsed() bool {
	if out, err := ioutil.ReadFile(ChargerDevice); err == nil && strings.TrimSpace(string(out)) == "1" {
		return false
	}
	return true
}

func (n *Battery) prettyValue() string {
	c := n.Capacity()
	if c == "n/a" {
		return "^fg(red)bat(n/a)^fg()"
	}
	val, err := strconv.Atoi(c)
	if err != nil {
		return "^fg(red)bat(n/a)^fg()"
	}
	switch {
	case val >= 20 && val <= 30:
		return icon(LowBatteryIcon) + "^fg(orange)" + c + "%^fg()"
	case val >= 10 && val < 20:
		return icon(EmptyBatteryIcon) + "^fg(red)" + c + "%^fg()"
	case val < 10:
		return "^bg(red)^fg(black)!connect charger!^fg()^bg()"
	default:
		return icon(FullBatteryIcon) + c + "%"
	}
}

func (n *Battery) Capacity() string {
	result := "n/a"
	if out, err := ioutil.ReadFile(BatteryDevice); err == nil {
		result = strings.TrimSpace(string(out))
	}
	return result
}
