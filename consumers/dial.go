package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joaoribeirodasilva/wait_signals"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	DeviceID    string      `json:"deviceId"`
	Sensors     interface{} `json:"sensors"`
	CollectedAt time.Time   `json:"collectedAt"`
}

// TODO: Define the base message object
type MetricsModel struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	UserID      primitive.ObjectID `json:"userId" bson:"userId"`
	DeviceID    string             `json:"deviceId"`
	Consumer    string             `json:"consumer" bson:"consumer"`
	Sensors     interface{}        `json:"sensors"`
	CollectedAt time.Time          `json:"collectedAt"`
	Received    time.Time          `json:"received" bson:"received"`
}

type Dial struct {
	conf          *Configuration
	broker        *MQTTClient
	isStarted     bool
	stopRequested bool
	finished      chan bool
	db            *Database
}

// NewDial create a new Dial struct pointer
func NewDial(conf *Configuration, db *Database) *Dial {

	d := &Dial{}

	d.conf = conf
	d.isStarted = false
	d.stopRequested = false
	d.broker = NewMQTTClient(conf, d.onMessageReceived)
	d.db = db

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

		d.broker.Connect()

		for !d.stopRequested {

			var err error
			if !d.broker.IsConnected() {
				err = d.broker.Connect()
			}
			if err == nil {
				d.broker.Subscribe(false)
				d.Publish()
			}

			if sig := wait_signals.SleepWait(time.Duration(d.conf.MQTT.Interval)*time.Millisecond, syscall.SIGINT, syscall.SIGTERM); sig != nil {
				break
			}
		}

		d.broker.Disconnect()

		if d.conf.Options.debug {
			log.Println("INFO: [DIAL] MQTT dial stopping")
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
		if d.conf.Options.debug {
			log.Println("INFO: [DIAL] MQTT dial requested to stop... waiting")
		}

		d.stopRequested = true

		<-d.finished

		log.Println("INFO: [DIAL] MQTT dial stopped")
	}
}

func (d *Dial) onMessageReceived(client mqtt.Client, message mqtt.Message) {

	msgJson := make(map[string]interface{})

	bytes := message.Payload()

	if err := json.Unmarshal(bytes, &msgJson); err != nil {
		log.Printf("ERROR: [DIAL] failed to parse message bytes REASON: %v\n", err)
		return
	}

	coll := d.db.GetCollection("metrics")

	// TODO: Adap to new model

	rec := MetricsModel{
		Consumer: d.conf.Mongo.ClientID,
		Metrics:  msgJson,
		Received: time.Now(),
	}

	_, err := coll.InsertOne(context.TODO(), rec)
	if err != nil {
		log.Printf("ERROR: [DIAL] failed to save message into database REASON: %v\n", err)
	}

	if d.conf.Options.debug {
		fmt.Printf("INFO: [DIAL] received message\n ")
	}
}

func (d *Dial) Publish() error {

	return nil
}

// TODO: validate base message device ID
func (d *Dial) CheckValidDevice() bool {
	return false
}
