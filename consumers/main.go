package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/joaoribeirodasilva/wait_signals"
)

type Options struct {
	consumer    int
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

	// initial subscribe and final unsubscribe
	if conf.Options.subscribe || conf.Options.unsubscribe {

		// connect to MQTT Broker
		broker := NewMQTTClient(conf, nil)
		if err := broker.Connect(); err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}

		exit := 0

		if conf.Options.subscribe {
			// subscribes to the MQTT topic in the configuration
			if err := broker.Subscribe(true); err != nil {
				log.Println(err.Error())
				exit = 1
			}
		}

		if conf.Options.unsubscribe && exit == 0 {
			// unsubscribes from the MQTT topic in the configuration
			if err := broker.Unsubscribe(); err != nil {
				log.Println(err.Error())
				exit = 1
			}
		}

		broker.Disconnect()
		os.Exit(exit)
	}

	// created a new dial
	dial := NewDial(conf, db)

	// starts the dial loop
	dial.Start()

	// wait for signals
	wait_signals.Wait(syscall.SIGINT, syscall.SIGTERM)

	// disconnect to MQTTBroker
	dial.Stop()

	// disconnect from mongodb
	db.Disconnect()

}

func cmdOptions() *Options {

	opts := &Options{}

	flag.IntVar(&opts.consumer, "c", 1, "consumer number")
	flag.BoolVar(&opts.subscribe, "s", false, "should subscribe topic on startup")
	flag.BoolVar(&opts.unsubscribe, "u", false, "should unsubscribe topic on startup")
	flag.BoolVar(&opts.help, "h", false, "print this help")

	flag.Parse()

	if opts.help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if opts.consumer == 0 {
		fmt.Println("ERROR: a -c with the consumer number is required [1-2]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	return opts
}
