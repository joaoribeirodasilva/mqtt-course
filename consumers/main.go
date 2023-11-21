package main

import (
	"flag"
	"log"
	"os"
)

type Options struct {
	subscribe   bool
	unsubscribe bool
	help        bool
}

func main() {

	// parse cmd line arguments
	opts := cmdOptions()

	// read configuration
	conf := NewConfiguration(opts)
	if err := conf.Read(); err != nil {
		panic(err)
	}

	// connect to mongodb
	db := NewDatabase(conf)
	if err := db.Connect(); err != nil {
		panic(err)
	}

	// connect to MQTT Broker
	broker := NewMQTTClient(conf)
	if err := broker.Connect(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	// initial subscribe and final unsubscribe
	if conf.Options.subscribe || conf.Options.unsubscribe {

		exit := 0
		if conf.Options.subscribe {
			if err := broker.Subscribe(); err != nil {
				log.Println(err.Error())
				exit = 1
			}
		}
		if conf.Options.subscribe && exit == 0 {
			if err := broker.Subscribe(); err != nil {
				log.Println(err.Error())
				exit = 1
			}
		}

		broker.Disconnect()
		os.Exit(exit)
	}

	// subscribe to MQTT topic

	// wait for signals

	// disconnect to MQTTBroker
	broker.Disconnect()

	// disconnect from mongodb
	db.Disconnect()

}

func cmdOptions() *Options {

	opts := &Options{}

	flag.BoolVar(&opts.subscribe, "s", false, "should subscribe topic on startup")
	flag.BoolVar(&opts.unsubscribe, "u", false, "should unsubscribe topic on startup")
	flag.BoolVar(&opts.help, "h", false, "print this help")

	flag.Parse()

	if opts.help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	return opts
}
