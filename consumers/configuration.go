package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type MongoConf struct {
	Uri      string `json:"uri"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TopicConf struct {
	Topic string `json:"topic"`
	Qos   byte   `json:"qos"`
}

type AuthConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CertificatesConf struct {
	Root string `json:"root"`
	Crt  string `json:"ctr"`
	Key  string `json:"key"`
}

type MQTTConf struct {
	Host           string
	Port           int
	Login          bool
	Tls            bool
	Publish        TopicConf
	Subscribe      TopicConf
	Authentication AuthConf
	Certificates   CertificatesConf
}

type Configuration struct {
	Mongo MongoConf
	MQTT  MQTTConf
}

const (
	defaultConfigPath = "config/config.json"
)

func NewConfiguration() *Configuration {

	c := &Configuration{}

	return c
}

func (c *Configuration) Read() error {

	log.Println("INFO: [CONFIGURATION] reading configuration")

	file, err := os.ReadFile(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("ERROR: failed to read configuration file: %s REASON: %s", defaultConfigPath, err.Error())
	}

	err = json.Unmarshal([]byte(file), c)
	if err != nil {
		return fmt.Errorf("ERROR: failed to parse configuration file: %s REASON: %s", defaultConfigPath, err.Error())
	}

	return nil

}
