package main

import (
	"log"
	"math/rand"
	"time"
)

type SensorDoor struct {
	OpenTime  *time.Time `json:"OpenTime"`
	CloseTime *time.Time `json:"CloseTime"`
	IsOpen    bool       `json:"IsOpen"`
}

type SensorNumeric struct {
	CurrentValue float64 `json:"CurrentValue"`
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

	s.currentStatus.Door.OpenTime = nil
	s.currentStatus.Door.CloseTime = nil
	s.currentStatus.Door.IsOpen = false

	s.currentStatus.Temperature.CurrentValue = conf.Sensors.Temperature.Normal

	s.currentStatus.Humidity.CurrentValue = conf.Sensors.Humidity.Normal

	s.currentStatus.DeviceID = conf.MQTT.ClientID
	s.currentStatus.CollectedAt = time.Now().UTC()

	return s
}

// Simulate runs a sensor simulation and store the results in the DataList object
func (s *Simulation) Simulate(virtualTime time.Time) {

	// strTime := virtualTime.Format(time.RFC3339)
	// log.Printf("INFO: [SIMULATE] Time now: %s\n", strTime)

	s.door(virtualTime)
	s.temperature(virtualTime)
	s.humidity(virtualTime)

	s.currentStatus.DeviceID = s.conf.MQTT.ClientID
	s.currentStatus.CollectedAt = time.Now().UTC()

	// log.Printf("INFO: [SIMULATE] Door is open: %t\n", s.currentStatus.Door.IsOpen)
	// if s.currentStatus.Door.OpenTime != nil {
	// 	strTime := s.currentStatus.Door.OpenTime.Format(time.RFC3339)
	// 	log.Printf("INFO: [SIMULATE] Door open time: %s\n", strTime)
	// } else {
	// 	log.Println("INFO: [SIMULATE] Door open time: nil")
	// }

	// if s.currentStatus.Door.CloseTime != nil {
	// 	strTime := s.currentStatus.Door.CloseTime.Format(time.RFC3339)
	// 	log.Printf("INFO: [SIMULATE] Door close time: %s\n", strTime)
	// } else {
	// 	log.Println("INFO: [SIMULATE] Door close time: nil")
	// }

	// log.Printf("INFO: [SIMULATE] Temperature: %.2f\n", s.currentStatus.Temperature.CurrentValue)
	// log.Printf("INFO: [SIMULATE] Humidity: %.2f\n", s.currentStatus.Humidity.CurrentValue)

	s.statusList.Append(s.currentStatus)

}

func (s *Simulation) door(virtualTime time.Time) {

	// if the door is closed
	if !s.currentStatus.Door.IsOpen {

		// door is not open so set the close time to zero
		s.currentStatus.Door.CloseTime = nil

		// check if the door opens now
		rnd := rand.Float64()
		if rnd <= s.conf.Sensors.DoorOpen.Chance {

			// if the door opens then calculate the amount of time it will remain open
			randomTimeOpen := rand.Int63n(s.conf.Sensors.DoorOpen.MaxTime-s.conf.Sensors.DoorOpen.MinTime) + s.conf.Sensors.DoorOpen.MinTime

			if s.conf.Options.debug {
				log.Printf("INFO: [SIMULATE] door will be open for %d milliseconds", randomTimeOpen)
			}

			calculatedClose := virtualTime.Add(time.Duration(randomTimeOpen) * time.Millisecond)

			s.currentStatus.Door.CloseTime = &calculatedClose
			// set the door to open
			s.currentStatus.Door.IsOpen = true

		}
	} else if s.currentStatus.Door.IsOpen && s.currentStatus.Door.CloseTime.Sub(virtualTime) < 0 {

		// if the door is open and it's time to close the door

		// set the open time to nil
		s.currentStatus.Door.OpenTime = nil

		// if the door is open check if it's time to the door to close
		// if so close the door
		s.currentStatus.Door.IsOpen = false

		// set the time the metric is being closed
		s.currentStatus.Door.CloseTime = &virtualTime
		// set the open time to zero
		s.currentStatus.Door.OpenTime = nil
	}

}

func (s *Simulation) temperature(virtualTime time.Time) {

	// calculate temperature if door is open or temperature is != normal
	if s.currentStatus.Door.IsOpen {

		// if the door is open we must increase the temperature until it's maximum
		s.currentStatus.Temperature.CurrentValue += s.conf.Sensors.Temperature.Increase
		if s.currentStatus.Temperature.CurrentValue > s.conf.Sensors.Temperature.Max {

			s.currentStatus.Temperature.CurrentValue = s.conf.Sensors.Temperature.Max
		}

	} else if s.currentStatus.Temperature.CurrentValue > s.conf.Sensors.Temperature.Normal {

		// if door is closed and the current temperature is higher that the
		// normal temperature we need to decrease the temperature until it
		// reaches the normal temperature
		s.currentStatus.Temperature.CurrentValue -= s.conf.Sensors.Temperature.Decrease
		if s.currentStatus.Temperature.CurrentValue < s.conf.Sensors.Temperature.Normal {

			s.currentStatus.Temperature.CurrentValue = s.conf.Sensors.Temperature.Normal
		}
	}

}

func (s *Simulation) humidity(virtualTime time.Time) {

	// calculate humidity if door is open or humidity is != normal
	if s.currentStatus.Door.IsOpen {

		// if the door is open we must increase the humidity until it's maximum
		s.currentStatus.Humidity.CurrentValue += s.conf.Sensors.Humidity.Increase
		if s.currentStatus.Humidity.CurrentValue > s.conf.Sensors.Humidity.Max {

			s.currentStatus.Humidity.CurrentValue = s.conf.Sensors.Humidity.Max
		}

	} else if s.currentStatus.Humidity.CurrentValue > s.conf.Sensors.Humidity.Normal {

		// if door is closed and the current humidity is higher that the
		// normal humidity we need to decrease the humidity until it
		// reaches the normal humidity
		s.currentStatus.Humidity.CurrentValue -= s.conf.Sensors.Humidity.Decrease
		if s.currentStatus.Humidity.CurrentValue < s.conf.Sensors.Humidity.Normal {

			s.currentStatus.Humidity.CurrentValue = s.conf.Sensors.Humidity.Normal
		}
	}

}
