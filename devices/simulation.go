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
	conf          *Configuration
	currentStatus Sensors
	statusList    *DataList
}

// New simulation creates a new simulation struct
func NewSimulation(conf *Configuration, list *DataList) *Simulation {

	s := &Simulation{}

	s.conf = conf
	s.statusList = list

	s.currentStatus.Door.openTime = time.Unix(0, 0)
	s.currentStatus.Door.closeTime = time.Unix(0, 0)
	s.currentStatus.Door.recordedTime = time.Now()
	s.currentStatus.Door.isOpen = false

	s.currentStatus.Temperature.currentValue = conf.Sensors.Temperature.Normal
	s.currentStatus.Temperature.recordedTime = time.Now()

	s.currentStatus.Humidity.currentValue = conf.Sensors.Humidity.Normal
	s.currentStatus.Humidity.recordedTime = time.Now()

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
	if !s.currentStatus.Door.isOpen {

		// door is not open so set the close time to zero
		s.currentStatus.Door.closeTime = time.Unix(0, 0)

		// check if the door opens now
		rnd := rand.Float64()
		if rnd <= s.conf.Sensors.DoorOpen.Chance {

			// if the door opens then calculate the amount of time it will remain open
			randomTimeOpen := rand.Int63n((s.conf.Sensors.DoorOpen.MaxTime - s.conf.Sensors.DoorOpen.MinTime) + s.conf.Sensors.DoorOpen.MinTime)
			s.currentStatus.Door.closeTime = virtualTime.Add(time.Duration(randomTimeOpen) * time.Millisecond)
			// set the door to open
			s.currentStatus.Door.isOpen = true
		}
	} else if s.currentStatus.Door.closeTime.Sub(virtualTime) < 0 {

		// if the door is open check if it's time to the door to close
		// if so close the door
		s.currentStatus.Door.isOpen = false
		// get the current metric time
		s.currentStatus.Door.recordedTime = virtualTime
		// set the time the metric is being closed
		s.currentStatus.Door.closeTime = virtualTime
		// set the open time to zero
		s.currentStatus.Door.openTime = time.Unix(0, 0)
	}
}

func (s *Simulation) temperature(virtualTime time.Time) {

	// calculate temperature if door is open or temperature is != normal
	if s.currentStatus.Door.isOpen {

		// if the door is open we must increase the temperature until it's maximum
		s.currentStatus.Temperature.currentValue += s.conf.Sensors.Temperature.Increase
		if s.currentStatus.Temperature.currentValue > s.conf.Sensors.Temperature.Max {

			s.currentStatus.Temperature.currentValue = s.conf.Sensors.Temperature.Max
		}

	} else if s.currentStatus.Temperature.currentValue > s.conf.Sensors.Temperature.Normal {

		// if door is closed and the current temperature is higher that the
		// normal temperature we need to decrease the temperature until it
		// reaches the normal temperature
		s.currentStatus.Temperature.currentValue -= s.conf.Sensors.Temperature.Decrease
		if s.currentStatus.Temperature.currentValue < s.conf.Sensors.Temperature.Normal {

			s.currentStatus.Temperature.currentValue = s.conf.Sensors.Temperature.Normal
		}
	}

	s.currentStatus.Temperature.recordedTime = virtualTime
}

func (s *Simulation) humidity(virtualTime time.Time) {

	// calculate humidity if door is open or humidity is != normal
	if s.currentStatus.Door.isOpen {

		// if the door is open we must increase the humidity until it's maximum
		s.currentStatus.Humidity.currentValue += s.conf.Sensors.Humidity.Increase
		if s.currentStatus.Humidity.currentValue > s.conf.Sensors.Humidity.Max {

			s.currentStatus.Humidity.currentValue = s.conf.Sensors.Humidity.Max
		}

	} else if s.currentStatus.Humidity.currentValue > s.conf.Sensors.Humidity.Normal {

		// if door is closed and the current humidity is higher that the
		// normal humidity we need to decrease the humidity until it
		// reaches the normal humidity
		s.currentStatus.Humidity.currentValue -= s.conf.Sensors.Humidity.Decrease
		if s.currentStatus.Humidity.currentValue < s.conf.Sensors.Humidity.Normal {

			s.currentStatus.Humidity.currentValue = s.conf.Sensors.Humidity.Normal
		}
	}

	s.currentStatus.Humidity.recordedTime = virtualTime
}
