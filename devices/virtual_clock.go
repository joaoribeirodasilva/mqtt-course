package main

import "time"

type VirtualClock struct {
	conf          *Configuration
	virtualTime   time.Time
	isStarted     bool
	stopRequested bool
	simulation    *Simulation
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

	go func() {

		vc.virtualTime = time.Now()

		for !vc.stopRequested {

			time.Sleep(time.Duration(vc.conf.Clock.Interval))
			vc.virtualTime = vc.virtualTime.Add(time.Duration(vc.conf.Clock.Interval*vc.conf.Clock.Multiplier) * time.Millisecond)

			vc.simulation.Simulate(vc.virtualTime)

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
