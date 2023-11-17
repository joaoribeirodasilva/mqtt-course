package main

import (
	"encoding/json"
	"time"
)

type Dial struct {
	conf          *Configuration
	broker        *MQTTClient
	statusList    *DataList
	isStarted     bool
	stopRequested bool
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

	go func() {

		d.isStarted = true

		for !d.stopRequested {

			if !d.statusList.IsDirty() {
				continue
			}

			if err := d.broker.Connect(); err != nil {
				continue
			}

			d.broker.Subscribe()
			d.Publish()
			d.broker.Disconnect()

			time.Sleep(time.Duration(d.conf.Communications.MQTT.Interval) * time.Millisecond)

		}

		if err := d.broker.Connect(); err == nil {
			d.Publish()
			if d.conf.Options.unsubscribe {
				d.broker.Unsubscribe()
			}
			d.broker.Disconnect()
		}

		d.isStarted = false
		d.stopRequested = false
	}()

	return nil
}

// Start request to stop the auto MQTT communication functionality
func (d *Dial) Stop() {

	if d.isStarted {

		done := make(chan bool)

		d.stopRequested = true

		done <- !d.isStarted

		<-done
	}
}

// Publish publishes the DataList struct array to the MQTT Broker
// and removes the data sent from the array
func (d *Dial) Publish() error {

	for d.statusList.Len() > 0 {

		head := d.statusList.GetHead()

		bytes, err := json.Marshal(head)

		if err != nil {
			return err
		}

		if err = d.broker.Publish(bytes); err != nil {
			return err
		}

		d.statusList.Remove(1)
	}

	return nil
}
