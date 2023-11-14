package main

import "time"

type VirtualClock struct {
	conf          ConfigClock
	virtualTime   time.Time
	isStarted     bool
	stopRequested bool
}

func NewVirtualClock(conf *Configuration) *VirtualClock {

	vc := &VirtualClock{}

	vc.conf = conf.Clock
	vc.isStarted = false
	vc.stopRequested = false

	return vc
}

func (vc *VirtualClock) Start() error {

	if vc.isStarted {

		return nil
	}

	go func() {

		vc.virtualTime = time.Now()

		for !vc.stopRequested {

			time.Sleep(time.Duration(vc.conf.Interval))
			vc.virtualTime = vc.virtualTime.Add(time.Duration(vc.conf.Interval*vc.conf.Multiplier) * time.Millisecond)

			// simulate
			// send data

		}

		vc.isStarted = false
		vc.stopRequested = false
	}()

	return nil

}

func (vc *VirtualClock) Stop() {

	if vc.isStarted {

		done := make(chan bool)

		vc.stopRequested = true

		done <- !vc.isStarted

		<-done
	}
}

func (vc *VirtualClock) IsStarted() bool {

	return vc.isStarted
}
