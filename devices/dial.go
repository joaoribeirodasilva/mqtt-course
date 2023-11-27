package main

import (
	"encoding/json"
	"log"
	"syscall"
	"time"

	"github.com/joaoribeirodasilva/wait_signals"
)

type Dial struct {
	conf          *Configuration
	broker        *MQTTClient
	statusList    *DataList
	isStarted     bool
	stopRequested bool
	finished      chan bool
}

// NewDial create a new Dial struct pointer
func NewDial(conf *Configuration, list *DataList) *Dial {

	d := &Dial{}

	d.conf = conf
	d.isStarted = false
	d.stopRequested = false
	d.statusList = list
	d.broker = NewMQTTClient(conf)

	return d
}

// Start starts the auto MQTT communication functionality
func (d *Dial) Start() error {

	if d.isStarted {

		return nil
	}

	d.finished = make(chan bool, 1)

	go func() {

		log.Println("INFO: [DIAL] MQTT dial started")

		d.isStarted = true

		for !d.stopRequested {

			if d.statusList.IsDirty() {
				if err := d.broker.Connect(); err == nil {
					d.broker.Subscribe()
					d.Publish()
					d.broker.Disconnect()
				} else {
					log.Printf("ERROR: [DIAL] MQTT failed to connect REASON: %s", err.Error())
				}
			}

			if sig := wait_signals.SleepWait(time.Duration(d.conf.MQTT.Interval)*time.Millisecond, syscall.SIGINT, syscall.SIGTERM); sig != nil {
				break
			}
		}

		log.Println("INFO: [DIAL] MQTT dial stopping")

		if err := d.broker.Connect(); err == nil {
			d.Publish()
			if d.conf.Options.unsubscribe {
				d.broker.Unsubscribe()
			}
			d.broker.Disconnect()
		}

		d.isStarted = false
		d.stopRequested = false

		d.finished <- true
	}()

	return nil
}

// Start request to stop the auto MQTT communication functionality
func (d *Dial) Stop() {

	if d.isStarted {

		log.Println("INFO: [DIAL] MQTT dial requested to stop... waiting")

		d.stopRequested = true

		<-d.finished

		log.Println("INFO: [DIAL] MQTT dial stopped")
	}
}

// Publish publishes the DataList struct array to the MQTT Broker
// and removes the data sent from the array
func (d *Dial) Publish() error {

	// store how many messages were publish
	// to the MQTT Broker
	messageCount := 0

	// while we have items in the list
	// send them to the MQTT Broker
	for d.statusList.Len() > 0 {

		// get the list first item
		head := d.statusList.GetHead()

		// transform the list oldest item
		// into JSON bytes
		bytes, err := json.Marshal(head)
		if err != nil {
			return err
		}

		// publish the item into the MQTT Broker
		if err = d.broker.Publish(bytes); err != nil {
			return err
		}

		// remove the first item from the list
		d.statusList.Remove(1)
		messageCount++
	}

	log.Printf("INFO: [DIAL] published %d messages", messageCount)

	return nil
}
