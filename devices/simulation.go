package main

import (
	"math/rand"
	"time"
)

type SensorDoor struct {
	openTime     time.Time
	closeTime    time.Time
	isOpen       bool
	recordedTime time.Time
}

type SensorNumeric struct {
	currentValue float64
	recordedTime time.Time
}

type Simulation struct {
	conf          *ConfigSensors
	currentStatus Sensors
	statusList    *DataList
}

// New simulation creates a new simulation struct
func NewSimulation(conf *ConfigSensors, list *DataList) *Simulation {

	s := &Simulation{}

	s.conf = conf
	s.statusList = list

	s.currentStatus.door.openTime = time.Unix(0, 0)
	s.currentStatus.door.closeTime = time.Unix(0, 0)
	s.currentStatus.door.recordedTime = time.Now()
	s.currentStatus.door.isOpen = false

	s.currentStatus.temperature.currentValue = conf.Temperature.Normal
	s.currentStatus.temperature.recordedTime = time.Now()

	s.currentStatus.humidity.currentValue = conf.Humidity.Normal
	s.currentStatus.humidity.recordedTime = time.Now()

	return s
}

// Simulate runs a sensor simulation and store the results in the DataList object
func (s *Simulation) Simulate(virtualTime time.Time) {

	s.door(virtualTime)
	s.temperature(virtualTime)
	s.humidity(virtualTime)

	s.statusList.Append(s.currentStatus)

}

func (s *Simulation) door(virtualTime time.Time) {

	// if the door is closed
	if !s.currentStatus.door.isOpen {

		// door is not open so set the close time to zero
		s.currentStatus.door.closeTime = time.Unix(0, 0)

		// check if the door opens now
		rnd := rand.Float64()
		if rnd <= s.conf.DoorOpen.Chance {

			// if the door opens then calculate the amount of time it will remain open
			randomTimeOpen := rand.Int63n((s.conf.DoorOpen.MaxTime - s.conf.DoorOpen.MinTime) + s.conf.DoorOpen.MinTime)
			s.currentStatus.door.closeTime = virtualTime.Add(time.Duration(randomTimeOpen) * time.Millisecond)
			// set the door to open
			s.currentStatus.door.isOpen = true
		}
	} else if s.currentStatus.door.closeTime.Sub(virtualTime) < 0 {

		// if the door is open check if it's time to the door to close
		// if so close the door
		s.currentStatus.door.isOpen = false
		// get the current metric time
		s.currentStatus.door.recordedTime = virtualTime
		// set the time the metric is being closed
		s.currentStatus.door.closeTime = virtualTime
		// set the open time to zero
		s.currentStatus.door.openTime = time.Unix(0, 0)
	}
}

func (s *Simulation) temperature(virtualTime time.Time) {

	// calculate temperature if door is open or temperature is != normal
	if s.currentStatus.door.isOpen {

		// if the door is open we must increase the temperature until it's maximum
		s.currentStatus.temperature.currentValue += s.conf.Temperature.Increase
		if s.currentStatus.temperature.currentValue > s.conf.Temperature.Max {

			s.currentStatus.temperature.currentValue = s.conf.Temperature.Max
		}

	} else if s.currentStatus.temperature.currentValue > s.conf.Temperature.Normal {

		// if door is closed and the current temperature is higher that the
		// normal temperature we need to decrease the temperature until it
		// reaches the normal temperature
		s.currentStatus.temperature.currentValue -= s.conf.Temperature.Decrease
		if s.currentStatus.temperature.currentValue < s.conf.Temperature.Normal {

			s.currentStatus.temperature.currentValue = s.conf.Temperature.Normal
		}
	}

	s.currentStatus.temperature.recordedTime = virtualTime
}

func (s *Simulation) humidity(virtualTime time.Time) {

	// calculate humidity if door is open or humidity is != normal
	if s.currentStatus.door.isOpen {

		// if the door is open we must increase the humidity until it's maximum
		s.currentStatus.humidity.currentValue += s.conf.Humidity.Increase
		if s.currentStatus.humidity.currentValue > s.conf.Humidity.Max {

			s.currentStatus.humidity.currentValue = s.conf.Humidity.Max
		}

	} else if s.currentStatus.humidity.currentValue > s.conf.Humidity.Normal {

		// if door is closed and the current humidity is higher that the
		// normal humidity we need to decrease the humidity until it
		// reaches the normal humidity
		s.currentStatus.humidity.currentValue -= s.conf.Humidity.Decrease
		if s.currentStatus.humidity.currentValue < s.conf.Humidity.Normal {

			s.currentStatus.humidity.currentValue = s.conf.Humidity.Normal
		}
	}

	s.currentStatus.humidity.recordedTime = virtualTime
}
