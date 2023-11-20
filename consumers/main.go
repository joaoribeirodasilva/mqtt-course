package main

func main() {

	// read configuration
	conf := NewConfiguration()
	if err := conf.Read(); err != nil {
		panic(err)
	}

	// connect to mongodb
	// connect to MQTT Broker
	// subscribe to MQTT topic

	// wait for signals

	// disconnect to MQTTBroker
	// disconnect from mongodb

}
