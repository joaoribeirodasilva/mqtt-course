package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type TLSConf struct {
	Use  bool   `json:"use"`
	Root string `json:"root"`
	Crt  string `json:"ctr"`
	Key  string `json:"key"`
}

type MongoCompressorsConf struct {
	Snappy bool `json:"snappy"`
	Zlib   bool `json:"zlib"`
	Zstd   bool `json:"zstd"`
}

type MongoWriteConcernConf struct {
	W          int  `json:"w"`
	WTimeoutMS int  `json:"wTimeoutMS"`
	Journal    bool `json:"journal"`
}

type MongoReadPreferenceConf struct {
	ReadPreference string `json:"readPreference"`
}

type MongoConf struct {
	ClientID                 string                  `json:"clientId"`
	Uri                      string                  `json:"uri"`
	Database                 string                  `json:"database"`
	Username                 string                  `json:"username"`
	Password                 string                  `json:"password"`
	TimeoutMS                int                     `json:"timeoutMS"`
	ConnectTimeoutMS         int                     `json:"connectTimeoutMS"`
	MaxPoolSize              int                     `json:"maxPoolSize"`
	ReplicaSet               string                  `json:"replicaSet"`
	MaxIdleTimeMS            int                     `json:"maxIdleTimeMS"`
	MinPoolSize              int                     `json:"minPoolSize"`
	SocketTimeoutMS          int                     `json:"socketTimeoutMS"`
	ServerSelectionTimeoutMS int                     `json:"serverSelectionTimeoutMS"`
	HeartbeatFrequencyMS     int                     `json:"heartbeatFrequencyMS"`
	Tls                      TLSConf                 `json:"tls"`
	Compressors              MongoCompressorsConf    `json:"compressors"`
	WriteConcern             MongoWriteConcernConf   `json:"writeConcern"`
	ReadPreference           MongoReadPreferenceConf `json:"readPreference"`
	DirectConnection         bool                    `json:"directConnection"`
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
	Publish        TopicConf `json:"publish"`
	Subscribe      TopicConf `json:"subscribe"`
	Authentication AuthConf  `json:"authentication"`
	Tls            TLSConf   `json:"tls"`
}

type Configuration struct {
	Options  *Options  `json:"-"`
	ClientID string    `json:"clientId"`
	Mongo    MongoConf `json:"mongodb"`
	MQTT     MQTTConf  `json:"mqtt"`
}

const (
	defaultConfigPath = "config/config.json"
)

func NewConfiguration(opts *Options) *Configuration {

	c := &Configuration{}
	c.Options = opts

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
