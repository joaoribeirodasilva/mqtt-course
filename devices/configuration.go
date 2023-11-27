package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type ClockConf struct {
	Interval uint64 `json:"interval"`
	Multiplier uint64 `json:"multiplier"`
}

type SensorDoorConf struct {
	Chance  float64 `json:"chance"`
	MaxTime int64   `json:"maxTime"`
	MinTime int64   `json:"minTime"`
}

type SensorNumericConf struct {
	Normal   float64 `json:"normal"`
	Increase float64 `json:"increase"`
	Decrease float64 `json:"decrease"`
	Max      float64 `json:"max"`
}

type SensorsConf struct {
	DoorOpen    SensorDoorConf    `json:"doorOpen"`
	Temperature SensorNumericConf `json:"temperature"`
	Humidity    SensorNumericConf `json:"humidity"`
}

type DataConf struct {
	Path         string `json:"path"`
	SaveInterval int64  `json:"saveInterval"`
	MaxMessages  uint32 `json:"maxMessages"`
}

type TLSConf struct {
	Use  bool   `json:"use"`
	Root string `json:"root"`
	Crt  string `json:"ctr"`
	Key  string `json:"key"`
}

type TopicConf struct {
	Topic string `json:"topic"`
	Qos   byte   `json:"qos"`
}


type AuthConf struct {
	Use      bool   `json:"use"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type MQTTConf struct {
	ClientID       string    `json:"clientId"`
	Host           string    `json:"host"`
	Port           int       `json:"port"`
	Interval int64 `json:"interval"`
	Publish        TopicConf `json:"publish"`
	Subscribe      TopicConf `json:"subscribe"`
	Authentication AuthConf  `json:"authentication"`
	Tls            TLSConf   `json:"tls"`
}

type ApiConf struct {
	Host  string `json:"certificates"`
	Port  int    `json:"port"`
	Token string `json:"token"`
}

type Configuration struct {
	device         int                  `json:"-"`
	ID             string               `json:"id"`
	Account        string               `json:"account"`
	Clock          ClockConf          `json:"clock"`
	Sensors        SensorsConf        `json:"sensors"`
	Data           DataConf                 `json:"data"`
	MQTT 		   MQTTConf  			`json:"mqtt"`
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
