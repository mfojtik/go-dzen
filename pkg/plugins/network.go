package plugins

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/golang/glog"
)

type Network struct {
	Interfaces []string
}

type Interface struct {
	Name string
}

func (n *Network) Stream() chan string {
	out := make(chan string)
	timer := time.NewTicker(time.Duration(1 * time.Second))
	go func() {
		for {
			<-timer.C
			glog.Infof("Sending network stats")
			result := ""
			for _, i := range n.Interfaces {
				r := n.Interface(i)
				result += fmt.Sprintf("%s(%s)", i, r.Status())
			}
			out <- result
		}
	}()
	return out
}

func (n *Network) Interface(name string) *Interface {
	return &Interface{Name: name}
}

func (i *Interface) Status() string {
	device := fmt.Sprintf("/sys/class/net/%s/carrier", i.Name)
	if out, err := ioutil.ReadFile(device); err == nil {
		if strings.TrimSpace(string(out)) == "1" {
			return "up"
		}
	}
	return "down"
}
