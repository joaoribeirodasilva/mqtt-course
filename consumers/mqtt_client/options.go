package mqtt_client

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TopicOptions struct {
	Name   string `json:"name"`
	Qos    byte   `json:"qos"`
	Retain bool   `json:"retain"`
}

type AuthenticationOptions struct {
	Use      bool   `json:"use"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type OnConnectionAttemptFunc func(broker *url.URL, tlsCfg *tls.Config) *tls.Config
type OnOpenConnectionFunc func(uri *url.URL, options ClientOptions) (net.Conn, error)
type OnConnectFunc func(client mqtt.Client)
type OnConnectionLostFunc func(client mqtt.Client, err error)
type OnReconnectFunc func(client mqtt.Client, opts *mqtt.ClientOptions)
type OnReceiveFunc func(client mqtt.Client, message mqtt.Message)
type OnPublishedFunc func(client mqtt.Client, message mqtt.Message)

type ClientOptions struct {
	ClientID            string                        `json:"clientId"`
	Host                string                        `json:"host"`
	Port                int                           `json:"port"`
	Publish             []TopicOptions                `json:"publish"`
	Subscribe           []TopicOptions                `json:"subscribe"`
	Authentication      AuthenticationOptions         `json:"authentication"`
	UseTls              bool                          `json:"useTls"`
	Tls                 *tls.Config                   `json:"tls"`
	CleanSession        bool                          `json:"cleanSession"`
	AutoAckDisabled     bool                          `json:"disableAutoAcknowledge"`
	AutoReconnect       bool                          `json:"autoReconnect"`
	PersistPath         string                        `json:"persistPath"`
	url                 string                        `json:"-"`
	options             *mqtt.ClientOptions           `json:"-"`
	onConnectionAttempt mqtt.ConnectionAttemptHandler `json:"-"`
	onOpenConnection    mqtt.OpenConnectionFunc       `json:"-"`
	onConnect           mqtt.OnConnectHandler         `json:"-"`
	onConnectionLost    mqtt.ConnectionLostHandler    `json:"-"`
	onReconnect         mqtt.ReconnectHandler         `json:"-"`
	onReceive           mqtt.MessageHandler           `json:"-"`
	onPublished         mqtt.MessageHandler           `json:"-"`
}

const (
	defaultPort     = 1883
	defaultHost     = "localhost"
	defaultProtocol = "tcp"
)

func NewClientOptions() *ClientOptions {

	co := &ClientOptions{}

	return co
}

func (co *ClientOptions) Factory() (*mqtt.ClientOptions, error) {

	co.options = mqtt.NewClientOptions()
	co.options.SetClientID(co.ClientID)

	// create the TLS configuration
	proto := defaultProtocol
	if co.UseTls {
		proto = "ssl"
		co.options.SetTLSConfig(co.Tls)
	}

	if co.Port == 0 {
		co.Port = defaultPort
	}

	if co.Host == "" {
		co.Host = defaultHost
	}

	co.url = fmt.Sprintf("%s://%s:%d", proto, co.Host, co.Port)
	co.options.AddBroker(co.url)

	co.options.OnConnectAttempt = co.onConnectionAttempt
	co.options.CustomOpenConnectionFn = co.onOpenConnection
	co.options.OnConnect = co.onConnect
	co.options.OnConnectionLost = co.onConnectionLost
	co.options.OnReconnecting = co.onReconnect
	co.options.DefaultPublishHandler = co.onPublished

	if co.Authentication.Use {
		if co.Authentication.Username == "" {
			return nil, fmt.Errorf("no username set")
		}
		if co.Authentication.Password == "" {
			return nil, fmt.Errorf("no password set")
		}
		co.options.SetUsername(co.Authentication.Username)
		co.options.SetPassword(co.Authentication.Password)

	}

	co.options.SetCleanSession(co.CleanSession)

	return co.options, nil
}

func (co *ClientOptions) GetUrl() string {
	return co.url
}

func (co *ClientOptions) GetOptions() *mqtt.ClientOptions {

	return co.options
}

func (co *ClientOptions) SetOnConnectionAttempt(f mqtt.ConnectionAttemptHandler) {
	co.onConnectionAttempt = f
}

func (co *ClientOptions) SetOnOpenConnection(f mqtt.OpenConnectionFunc) {
	co.onOpenConnection = f
}

func (co *ClientOptions) SetOnConnect(f mqtt.OnConnectHandler) {
	co.onConnect = f
}

func (co *ClientOptions) SetOnConnectionLost(f mqtt.ConnectionLostHandler) {
	co.onConnectionLost = f
}

func (co *ClientOptions) SetOnReconnect(f mqtt.ReconnectHandler) {
	co.onReconnect = f
}

func (co *ClientOptions) SetOnReceive(f mqtt.MessageHandler) {
	co.onReceive = f
}

func (co *ClientOptions) SetOnPublished(f mqtt.MessageHandler) {
	co.onPublished = f
}
