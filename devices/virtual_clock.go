package main

import (
	"log"
	"time"
)

type VirtualClock struct {
	conf          *Configuration
	virtualTime   time.Time
	isStarted     bool
	stopRequested bool
	simulation    *Simulation
	finished      chan bool
}

func NewVirtualClock(conf *Configuration, simulation *Simulation) *VirtualClock {

	vc := &VirtualClock{}

	vc.conf = conf
	vc.isStarted = false
	vc.stopRequested = false
	vc.simulation = simulation

	return vc
}

func (vc *VirtualClock) Start() error {

	if vc.isStarted {

		return nil
	}

	// make/remake the finished channel
	vc.finished = make(chan bool, 1)

	go func() {

		log.Println("INFO: [VIRTUAL CLOCK] virtual clock started")

		vc.virtualTime = time.Now()

		for !vc.stopRequested {

			if !SleepChannel(time.Duration(vc.conf.Clock.Interval) * time.Millisecond) {
				break
			}

			vc.virtualTime = vc.virtualTime.Add(time.Duration(vc.conf.Clock.Interval*vc.conf.Clock.Multiplier) * time.Millisecond)
			vc.simulation.Simulate(vc.virtualTime)

		}

		log.Println("INFO: [VIRTUAL CLOCK] virtual clock requested ro stop")
		vc.isStarted = false
		vc.stopRequested = false

		// set the channel so the Stop function can stop waiting
		// for loop termination
		vc.finished <- true
	}()

	return nil

}

func (vc *VirtualClock) Stop() {

	if vc.isStarted {

		log.Println("INFO: [VIRTUAL CLOCK] virtual clock stop requested")

		vc.stopRequested = true

		<-vc.finished

		log.Println("INFO: [VIRTUAL CLOCK] virtual clock stopped... waiting")
	}
}

func (vc *VirtualClock) IsStarted() bool {

	return vc.isStarted
}
