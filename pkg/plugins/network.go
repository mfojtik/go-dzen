package plugins

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/mfojtik/go-dzen/pkg/util"
)

const (
	IpPath      = "/usr/sbin/ip"
	NetworkUp   = ".icons/dzen2/net_up_03.xbm"
	NetworkDown = ".icons/dzen2/net_down_03.xbm"
)

var InterfaceNames = map[string]string{
	"wlp3s0":  "wifi",
	"enp0s25": "eth",
	"tun0":    "vpn",
}

type Network struct {
	Interfaces []string
}

func (n Network) String() string { return "network" }

type Interface struct {
	Name string
}

func (n *Network) Stream() chan string {
	out := make(chan string)
	util.RunEvery(time.Duration(1*time.Second), func() error {
		result := ""
		for _, i := range n.Interfaces {
			r := n.Interface(i)
			result += r.String()
		}
		out <- result
		return nil
	})
	return out
}

func (n *Network) Interface(name string) *Interface {
	return &Interface{Name: name}
}

func (i *Interface) String() string {
	statusColor := "red"
	ipAddress := ""
	networkIcon := icon(NetworkDown)
	if i.IsUp() {
		statusColor = "green"
		ipAddress = "(" + i.IP() + ")"
		networkIcon = icon(NetworkUp)
	}
	return fmt.Sprintf("%s%s^fg(%s)%s^fg()%s", Separator, networkIcon, statusColor, InterfaceNames[i.Name], ipAddress)
}

func (i *Interface) IP() string {
	out, err := exec.Command(IpPath, "-o", "-4", "addr", "show", i.Name).Output()
	if err != nil {
		glog.Errorf(err.Error())
		return ""
	}
	parts := strings.Split(string(out), " ")
	result := ""
	for i, s := range parts {
		if strings.TrimSpace(s) == "inet" {
			result = strings.TrimSpace(parts[i+1])
		}
	}
	if len(result) == 0 {
		return ""
	}
	return strings.Split(result, "/")[0]
}

func (i *Interface) IsUp() bool {
	device := fmt.Sprintf("/sys/class/net/%s/carrier", i.Name)
	if out, err := ioutil.ReadFile(device); err == nil {
		if strings.TrimSpace(string(out)) == "1" {
			return true
		}
	}
	return false
}
