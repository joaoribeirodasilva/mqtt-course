package mqtt_client

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	conf   *ClientOptions
	client mqtt.Client
	token  mqtt.Token
}

func NewClient(conf *ClientOptions) *MQTTClient {

	c := &MQTTClient{}

	c.conf = conf

	return c
}

func (c *MQTTClient) Connect() error {

	if c.client.IsConnected() {
		return nil
	}

	c.client = mqtt.NewClient(c.conf.options)

	if c.token = c.client.Connect(); c.token.Wait() && c.token.Error() != nil {

		return fmt.Errorf("failed to connect to MQTT Broker at %s REASON: %s", c.conf.GetUrl(), c.token.Error().Error())
	}

	return nil
}

//TODO: Use topics now

func (c *MQTTClient) SubscribeAll() error {

	for _, subscriber := range c.conf.Subscribe {
		c.subscribe(subscriber.Name, subscriber.Qos)
	}
	return nil
}

func (c *MQTTClient) subscribe(topic string, qos byte) error {

	if c.token = c.client.Subscribe(topic, qos, c.conf.onPublished); c.token.Wait() && c.token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s REASON: %s", topic, c.token.Error().Error())
	}

	return nil
}

func (c *MQTTClient) UnsubscribeAll() error {

	return nil
}

func (c *MQTTClient) unsubscribe(topic string) error {

	if c.token = c.client.Unsubscribe(topic); c.token.Wait() && c.token.Error() != nil {

		return fmt.Errorf("failed to unsubscribe from topic %s REASON: %s", topic, c.token.Error().Error())
	}

	return nil
}

func (c *MQTTClient) Publish(topic string, qos byte, data []byte) error {

	if c.token = c.client.Publish(topic, qos, false, data); c.token.Wait() && c.token.Error() != nil {

		return c.token.Error()
	}

	return nil
}

func (c *MQTTClient) Disconnect() {

	if c.client.IsConnected() {
		c.client.Disconnect(250)
		log.Println("disconnected from MQTT Broker")
	}
}

func (c *MQTTClient) IsConnected() bool {

	return c.client.IsConnected()
}

// func (c *MQTTClient) onMessagePublishedHandler(client mqtt.Client, msg mqtt.Message) {

// 	// Message published to the broker
// }

// func (c *MQTTClient) onConnectHandler(client mqtt.Client) {

// 	log.Println("INFO: [MQTT CLIENT] connected to MQTT Broker")
// 	c.isConnected = true
// }

// func (c *MQTTClient) onConnectLostHandler(client mqtt.Client, err error) {

// 	c.isConnected = false
// }

// func (c *MQTTClient) onMessageReceived(client mqtt.Client, message mqtt.Message) {

// 	// TODO: callback
// 	// Received message from the broker
// 	// Switch for operation to perform
// }
