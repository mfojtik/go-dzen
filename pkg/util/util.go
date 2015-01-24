package util

import (
	"time"

	"github.com/golang/glog"
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
