package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Options struct {
	device      int
	subscribe   bool
	unsubscribe bool
	configFile  string
	noTls       bool
	qos         int
	help        bool
}

func main() {

	fmt.Println("MQTT Course device")

	// parse cmd line options
	opts := cmdOptions()

	// read configuration
	conf := NewConfiguration(opts)
	err := conf.Read()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Subscribe or unsubscribe functions
	// TODO:

	// Initialize the simulation list
	list := NewDataList(conf)
	list.Read()
	list.Start()

	// Initialize the simulation structure
	simulation := NewSimulation(conf, list)

	// Start the virtual clock
	clock := NewVirtualClock(conf, simulation)
	clock.Start()

	// Start the communication with MQTT Broker
	dial := NewDial(conf, list)
	dial.Start()

	// wait SIGTERM or SIGINT
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// wait for clock to stop
	clock.Stop()

	// wait for list to stop
	list.Stop()

	// wait for dial to stop
	dial.Stop()

	// unsubscribe (if set)
	// disconnect all communications

}

func cmdOptions() *Options {

	opts := &Options{}

	flag.StringVar(&opts.configFile, "-c", "", "configuration file path")
	flag.IntVar(&opts.device, "-d", 0, "device number [1-3]")
	flag.BoolVar(&opts.help, "-h", false, "help")
	flag.BoolVar(&opts.noTls, "--no-tls", false, "don't use tls certificates")
	flag.IntVar(&opts.qos, "-q", 2, "MQTT QOS level")
	flag.BoolVar(&opts.subscribe, "-s", false, "should subscribe topic on startup")
	flag.BoolVar(&opts.unsubscribe, "-u", false, "should unsubscribe topic on startup")

	flag.Parse()

	if opts.device == 0 && opts.configFile == "" && !opts.help {
		fmt.Println("ERROR: a -c or -d option is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if opts.help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	return opts
}
