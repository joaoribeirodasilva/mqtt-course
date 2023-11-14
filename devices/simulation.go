package main

import "time"

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
	conf        *ConfigSensors
	door        SensorDoor
	temperature SensorNumeric
	humidity    SensorNumeric
}

func NewSimulation(conf *ConfigSensors) *Simulation {

	s := &Simulation{}

	s.conf = conf

	s.door.openTime = time.Unix(0, 0)
	s.door.closeTime = time.Unix(0, 0)
	s.door.recordedTime = time.Now()
	s.door.isOpen = false

	s.temperature.currentValue = conf.Temperature.Normal
	s.temperature.recordedTime = time.Now()

	s.humidity.currentValue = conf.Humidity.Normal
	s.humidity.recordedTime = time.Now()

	return s
}

func Simulate() {
	// TODO: think the best way to make this calculations
}
