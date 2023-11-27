package configuration

import (
	"encoding/json"
	"fmt"
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

type ServerConf struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	JwtKey  string `json:"jwtKey"`
}

type Configuration struct {
	Mongo  MongoConf  `json:"mongodb"`
	Server ServerConf `json:"server"`
}

const (
	defaultConfigurationfile = "./config/config.json"
)

func NewConfiguration() *Configuration {

	c := &Configuration{}

	return c
}

func (c *Configuration) Read() error {

	jsonBytes, err := os.ReadFile(defaultConfigurationfile)
	if err != nil {
		return fmt.Errorf("ERROR: [CONFIGURATION] failed to read config file %s", defaultConfigurationfile)
	}

	if err := json.Unmarshal(jsonBytes, c); err != nil {
		return fmt.Errorf("ERROR: [CONFIGURATION] decoding configuration file %s", defaultConfigurationfile)
	}

	return nil

}
