package main

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	clientId       string
	certificates   *ConfigCertificates
	authentication *ConfigAuthentication
	topic          *ConfigTopic
	mqttClient     mqtt.Client
	mqttToken      mqtt.Token
	noTls          bool
	isConnected    bool
}

func NewMQTTClient(opts *Options, id string, certificates *ConfigCertificates, authentication *ConfigAuthentication, topic *ConfigTopic) *MQTTClient {

	c := &MQTTClient{}

	c.noTls = opts.noTls
	c.certificates = certificates
	c.authentication = authentication
	c.topic = topic
	c.isConnected = false
	c.clientId = id

	return c
}

func (c *MQTTClient) Connect() error {

	if c.isConnected {

		return nil
	}

	brokerUrl := fmt.Sprintf("tcp://%s:%d", c.topic.Host, c.topic.Port)

	log.Printf("connecting to MQTT broker at %s ...", brokerUrl)

	options := mqtt.NewClientOptions()

	options.AddBroker(brokerUrl)

	options.SetClientID(c.clientId)
	options.SetDefaultPublishHandler(c.onMessagePublishedHandler)

	options.OnConnect = c.onConnectHandler
	options.OnConnectionLost = c.onConnectLostHandler

	if c.authentication.Username != "" && c.authentication.Password != "" {

		options.SetUsername(c.authentication.Username)
		options.SetPassword(c.authentication.Password)
	}

	c.mqttClient = mqtt.NewClient(options)

	if c.mqttToken = c.mqttClient.Connect(); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return fmt.Errorf("ERROR: failed to connect to MQTT Broker at %s. REASON: %s", brokerUrl, c.mqttToken.Error().Error())
	}

	log.Println(" connected")

	return nil
}

func (c *MQTTClient) Subscribe() error {

	if !c.isConnected {

		if err := c.Connect(); err != nil {

			return err
		}
	}

	log.Printf("subscribing MQTT topic %s with QOS %d ...", c.topic.Topic, c.topic.Qos)

	if c.mqttToken = c.mqttClient.Subscribe(c.topic.Topic, c.topic.Qos, c.onMessageReceived); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	log.Println(" subscribed")

	return nil
}

func (c *MQTTClient) Unsubscribe() error {

	if !c.isConnected {

		if err := c.Connect(); err != nil {

			return err
		}
	}

	log.Printf("unsubscribing from MQTT topic %s  ...", c.topic.Topic)

	if c.mqttToken = c.mqttClient.Unsubscribe(c.topic.Topic); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	log.Println(" unsubscribed")

	return nil
}

func (c *MQTTClient) Publish(data []byte) error { // change data for the list of stored metrics

	if !c.isConnected {

		if err := c.Connect(); err != nil {

			return err
		}
	}

	// Get the list of stored metrics

	log.Printf("publishing MQTT message into topic %s with QOS %d ...", c.topic.Topic, c.topic.Qos)

	if c.mqttToken = c.mqttClient.Publish(c.topic.Topic, c.topic.Qos, false, data); c.mqttToken.Wait() && c.mqttToken.Error() != nil {

		return c.mqttToken.Error()
	}

	log.Println(" published")

	return nil
}

func (c *MQTTClient) Disconnect() {

	if c.isConnected {
		c.mqttClient.Disconnect(250)
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

	c.isConnected = true
}

func (c *MQTTClient) onConnectLostHandler(client mqtt.Client, err error) {

	c.isConnected = false
}

func (c *MQTTClient) onMessageReceived(client mqtt.Client, message mqtt.Message) {

	// Received message from the broker
	// Switch for operation to perform
}
