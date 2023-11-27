package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	conf             *Configuration
	mqttClient       mqtt.Client
	mqttToken        mqtt.Token
	onReceiveMessage mqtt.MessageHandler
	isConnected      bool
}

func NewMQTTClient(conf *Configuration, onReceiveMessage mqtt.MessageHandler) *MQTTClient {

	c := &MQTTClient{}

	c.conf = conf
	c.isConnected = false
	c.onReceiveMessage = onReceiveMessage

	return c
}

func (c *MQTTClient) Connect() error {

	if c.isConnected {

		return nil
	}

	auth := ""
	if c.conf.MQTT.Authentication.Use {
		auth = fmt.Sprintf("%s:%s@", c.conf.MQTT.Authentication.Username, c.conf.MQTT.Authentication.Password)
	}

	brokerUrl := fmt.Sprintf("tcp://%s%s:%d", auth, c.conf.MQTT.Host, c.conf.MQTT.Port)

	//log.Printf("connecting to MQTT broker at %s ...", brokerUrl)

	options := mqtt.NewClientOptions()

	options.AddBroker(brokerUrl)

	options.SetClientID(c.conf.MQTT.ClientID)
	options.SetDefaultPublishHandler(c.onMessagePublishedHandler)

	options.OnConnect = c.onConnectHandler
	options.OnConnectionLost = c.onConnectLostHandler

	if c.conf.MQTT.Authentication.Use {

		options.SetUsername(c.conf.MQTT.Authentication.Username)
		options.SetPassword(c.conf.MQTT.Authentication.Password)
	}

	options.SetCleanSession(c.conf.Options.subscribe)

	// Add tls code

	c.mqttClient = mqtt.NewClient(options)
	if c.mqttToken = c.mqttClient.Connect(); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return fmt.Errorf("ERROR: [MQTT CLIENT] failed to connect to MQTT Broker at %s. REASON: %s", brokerUrl, c.mqttToken.Error().Error())
	}

	return nil
}

func (c *MQTTClient) Subscribe(verbose bool) error {

	if verbose {
		log.Printf("subscribing MQTT topic %s with QOS %d ...", c.conf.MQTT.Subscribe.Topic, c.conf.MQTT.Subscribe.Qos)
	}

	if c.mqttToken = c.mqttClient.Subscribe(c.conf.MQTT.Subscribe.Topic, c.conf.MQTT.Subscribe.Qos, c.onReceiveMessage); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	// log.Println(" subscribed")

	return nil
}

func (c *MQTTClient) Unsubscribe() error {

	log.Printf("unsubscribing from MQTT topic %s  ...", c.conf.MQTT.Subscribe.Topic)

	if c.mqttToken = c.mqttClient.Unsubscribe(c.conf.MQTT.Subscribe.Topic); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	//log.Println(" unsubscribed")

	return nil
}

func (c *MQTTClient) Publish(data []byte) error {

	if c.mqttToken = c.mqttClient.Publish(c.conf.MQTT.Publish.Topic, c.conf.MQTT.Publish.Qos, false, data); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	//log.Println(" published")

	return nil
}

func (c *MQTTClient) Disconnect() {

	if c.isConnected {
		c.mqttClient.Disconnect(250)
		log.Println("INFO: [MQTT CLIENT] disconnected from MQTT Broker")
		c.isConnected = false
	}
}

func (c *MQTTClient) IsConnected() bool {

	return c.isConnected
}

func (c *MQTTClient) onMessagePublishedHandler(client mqtt.Client, msg mqtt.Message) {

	// Message published to the broker
}

func (c *MQTTClient) onConnectHandler(client mqtt.Client) {

	log.Println("INFO: [MQTT CLIENT] connected to MQTT Broker")
	c.isConnected = true
}

func (c *MQTTClient) onConnectLostHandler(client mqtt.Client, err error) {

	log.Println("INFO: [MQTT CLIENT] MQTT Broker connection lost ")
	c.isConnected = false
}

// func (c *MQTTClient) onMessageReceived(client mqtt.Client, message mqtt.Message) {

// }
