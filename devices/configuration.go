package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// ConfigClock holds the information to set the virtual clock.
// The virtual clock is a way to accelerate time as if time is passing
// faster than it is.
// Ex: If interval is 100 real milliseconds and the multiplier is 60 then
// the virtual clock is advanced by 100*60=6000 milliseconds (one minute) for every real
// 100 milliseconds.
type ConfigClock struct {

	// The interval contains how many virtual milliseconds are in a real
	// millisecond and sets the loop sleep.
	Interval uint64 `json:"interval"`

	// The multiplier contains the value that needs to be multiplied by
	// the interval after an interval sleep and added to the date time
	// of the virtual clock.
	Multiplier uint64 `json:"multiplier"`
}

type ConfigSensorDoor struct {
	Chance  float64 `json:"chance"`
	MaxTime uint64  `json:"maxTime"`
	MinTime uint64  `json:"minTime"`
}

type ConfigSensorNumeric struct {
	Normal   float64 `json:"normal"`
	Increase float64 `json:"increase"`
	Decrease float64 `json:"decrease"`
}

type ConfigSensors struct {
	DoorOpen    ConfigSensorDoor    `json:"doorOpen"`
	Temperature ConfigSensorNumeric `json:"temperature"`
	Humidity    ConfigSensorNumeric `json:"humidity"`
}

type ConfigCertificates struct {
	Dir  string `json:"dir"`
	Root string `json:"root"`
	Crt  string `json:"crt"`
	Key  string `json:"key"`
}

type ConfigAuthentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConfigTopic struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Topic    string `json:"topic"`
	Interval uint64 `json:"interval"`
	Qos      byte   `json:"qos"`
}

type ConfigApi struct {
	Host  string `json:"certificates"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}

type ConfigCommunications struct {
	Certificates   ConfigCertificates   `json:"certificates"`
	Authentication ConfigAuthentication `json:"authentication"`
	Publish        ConfigTopic          `json:"publish"`
	Consume        ConfigTopic          `json:"consume"`
	Api            ConfigApi            `json:"api"`
}

type Configuration struct {
	device         int                  `json:"-"`
	ID             string               `json:"id"`
	Account        string               `json:"account"`
	Clock          ConfigClock          `json:"clock"`
	Sensors        ConfigSensors        `json:"sensors"`
	Communications ConfigCommunications `json:"communications"`
	configPath     string
}

const (
	defaultConfigPath = "config/device:num/config.json"
)

func NewConfiguration(opts *Options) *Configuration {

	conf := &Configuration{}

	conf.device = opts.device

	opts.configFile = defaultConfigPath
	if opts.configFile != "" {
		conf.configPath = opts.configFile
	}

	return conf
}

func (conf *Configuration) Read() error {

	file, err := os.ReadFile(conf.configPath)
	if err != nil {
		return fmt.Errorf("ERROR: failed to read configuration file: %s REASON: %s", conf.configPath, err.Error())
	}

	err = json.Unmarshal([]byte(file), conf)
	if err != nil {
		return fmt.Errorf("ERROR: failed to parse configuration file: %s REASON: %s", conf.configPath, err.Error())
	}

	return nil
}

func (conf *Configuration) Write() error {

	data, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("ERROR: failed to create JSON for configuration file REASON: %s", err.Error())
	}

	err = os.WriteFile(conf.configPath, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("ERROR: failed to write configuration file: %s REASON: %s", conf.configPath, err.Error())
	}

	return nil
}
