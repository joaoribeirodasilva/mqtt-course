package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
	MaxTime int64   `json:"maxTime"`
	MinTime int64   `json:"minTime"`
}

type ConfigSensorNumeric struct {
	Normal   float64 `json:"normal"`
	Increase float64 `json:"increase"`
	Decrease float64 `json:"decrease"`
	Max      float64 `json:"max"`
}

type ConfigSensors struct {
	DoorOpen    ConfigSensorDoor    `json:"doorOpen"`
	Temperature ConfigSensorNumeric `json:"temperature"`
	Humidity    ConfigSensorNumeric `json:"humidity"`
}

type Data struct {
	Path         string `json:"path"`
	SaveInterval int64  `json:"saveInterval"`
	MaxMessages  uint32 `json:"maxMessages"`
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

type ConfigMQTT struct {
	Host      string      `json:"host"`
	Port      int         `json:"port"`
	Interval  int64       `json:"interval"`
	Subscribe ConfigTopic `json:"subscribe"`
	Publish   ConfigTopic `json:"publish"`
	Login     bool        `json:"login"`
	Tls       bool        `json:"tls"`
}

type ConfigTopic struct {
	Publish bool   `json:"publish"`
	Topic   string `json:"topic"`
	Qos     byte   `json:"qos"`
}

type ConfigApi struct {
	Host  string `json:"certificates"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}

type ConfigCommunications struct {
	Certificates   ConfigCertificates   `json:"certificates"`
	Authentication ConfigAuthentication `json:"authentication"`
	MQTT           ConfigMQTT           `json:"mqtt"`
	Api            ConfigApi            `json:"api"`
}

type Configuration struct {
	device         int                  `json:"-"`
	ID             string               `json:"id"`
	Account        string               `json:"account"`
	Clock          ConfigClock          `json:"clock"`
	Sensors        ConfigSensors        `json:"sensors"`
	Data           Data                 `json:"data"`
	Communications ConfigCommunications `json:"communications"`
	Options        *Options             `json:"-"`
	configPath     string
}

const (
	defaultConfigPath = "config/device:num/config.json"
)

func NewConfiguration(opts *Options) *Configuration {

	conf := &Configuration{}

	conf.Options = opts
	conf.device = opts.device

	strDevice := fmt.Sprintf("%d", conf.device)

	opts.configFile = strings.Replace(defaultConfigPath, ":num", strDevice, 1)
	if opts.configFile != "" {
		conf.configPath = opts.configFile
	}

	return conf
}

func (conf *Configuration) Read() error {

	log.Println("INFO: [CONFIGURATION] reading client configuration")

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

	log.Println("INFO: [CONFIGURATION] writing client configuration")

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
