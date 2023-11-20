package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	conf        *Configuration
	mqttClient  mqtt.Client
	mqttToken   mqtt.Token
	isConnected bool
}

func NewMQTTClient(conf *Configuration) *MQTTClient {

	c := &MQTTClient{}

	c.conf = conf
	c.isConnected = false

	return c
}

func (c *MQTTClient) Connect() error {

	if c.isConnected {

		return nil
	}

	brokerUrl := fmt.Sprintf("tcp://%s:%d", c.conf.Communications.MQTT.Host, c.conf.Communications.MQTT.Port)

	//log.Printf("connecting to MQTT broker at %s ...", brokerUrl)

	options := mqtt.NewClientOptions()

	options.AddBroker(brokerUrl)

	options.SetClientID(c.conf.ID)
	options.SetDefaultPublishHandler(c.onMessagePublishedHandler)

	options.OnConnect = c.onConnectHandler
	options.OnConnectionLost = c.onConnectLostHandler

	if c.conf.Communications.MQTT.Login {

		options.SetUsername(c.conf.Communications.Authentication.Username)
		options.SetPassword(c.conf.Communications.Authentication.Password)
	}

	options.SetCleanSession(c.conf.Options.subscribe)

	// Add tls code

	c.mqttClient = mqtt.NewClient(options)

	if c.mqttToken = c.mqttClient.Connect(); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return fmt.Errorf("ERROR: [MQTT CLIENT] failed to connect to MQTT Broker at %s. REASON: %s", brokerUrl, c.mqttToken.Error().Error())
	}

	return nil
}

func (c *MQTTClient) Subscribe() error {

	// if !c.isConnected {

	// 	if err := c.Connect(); err != nil {

	// 		return err
	// 	}
	// }

	//log.Printf("subscribing MQTT topic %s with QOS %d ...", c.conf.Communications.MQTT.Subscribe.Topic, c.conf.Communications.MQTT.Subscribe.Qos)

	if c.mqttToken = c.mqttClient.Subscribe(c.conf.Communications.MQTT.Subscribe.Topic, c.conf.Communications.MQTT.Subscribe.Qos, c.onMessageReceived); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	// log.Println(" subscribed")

	return nil
}

func (c *MQTTClient) Unsubscribe() error {

	// if !c.isConnected {

	// 	if err := c.Connect(); err != nil {

	// 		return err
	// 	}
	// }

	//log.Printf("unsubscribing from MQTT topic %s  ...", c.conf.Communications.MQTT.Subscribe.Topic)

	if c.mqttToken = c.mqttClient.Unsubscribe(c.conf.Communications.MQTT.Subscribe.Topic); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	//log.Println(" unsubscribed")

	return nil
}

func (c *MQTTClient) Publish(data []byte) error {

	// if !c.isConnected {

	// 	if err := c.Connect(); err != nil {

	// 		return err
	// 	}
	// }

	// Get the list of stored metrics

	//log.Printf("publishing MQTT message into topic %s with QOS %d ...", c.conf.Communications.MQTT.Publish.Topic, c.conf.Communications.MQTT.Publish.Qos)

	if c.mqttToken = c.mqttClient.Publish(c.conf.Communications.MQTT.Publish.Topic, c.conf.Communications.MQTT.Publish.Qos, false, data); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

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

	c.isConnected = false
}

func (c *MQTTClient) onMessageReceived(client mqtt.Client, message mqtt.Message) {

	// Received message from the broker
	// Switch for operation to perform
}
