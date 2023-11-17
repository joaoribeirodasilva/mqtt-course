# Device configuration

In this document we will explain how the device configuration files work.

## Loading a configuration file

When the device program starts there are 2 options that influence which configuration file is loaded for a device program instance.

| Option | Type | Default | Required | Description |
| ------ | ---- | ------- | -------- | ----------- |
| -c | string | /config/device[number]/config.json | No | the -c option requires the full path to the instance configuration file. |
| -d | int | [none] | Yes | The -d option requires the device instance number. |

## Configuration file objects

This section explains all the configuration files keys and values.

The device configuration file must be in JSON format. No other formats are allowed at this moment.

## Main Configuration Object

The main configuration object has the following format.

| Key | Type | Required | Description |
| --- | ---- | -------- | ----------- |
| id  | ObjectID | Yes | MongoDB object id matching the id for the device in the database. |
| account | ObjectID | No | MongoDB object id matching the user account the device belongs to. |
| clock | [VirtualClock](#virtualclock-object) | Yes | VirtualClock configuration object. |
| sensors | [Sensors](#sensors-object) | Yes | Sensors configuration object. |
| communication | [Communication](#communication-object) | Yes | Communication configuration object. |

## VirtualClock Object

The VirtualClock is used to accelerate time, so we can test our code as if much more time has passed.

Configuration values for the virtual clock are as follows:

| Key | Type | Required | Description |
| --- | ---- | -------- | ----------- |
| interval | int | yes | The interval sets the code loop interval in milliseconds. At every interval the code executes a sensor simulation. |
| multiplier | int | yes | The multiplier sets the advancement of time in the clocks, multiplying it's value with the the value in the interval. |

**Example:** If the interval is set to be 100 milliseconds and the multiplier is set to 60 at every interval the virtual clock is advanced in 100 x 60 = 6000 milliseconds which corresponds to 1 minute. So for each 100 milliseconds in real time the virtual clock advances 1 minute.

## Sensors Object

The sensors block

## Communication Object

The communications configuration object holds the various settings for the communication with the MQTT Broker and the server API.

It has the following format:

| Key | Type | Required | Description |
| --- | ---- | -------- | ----------- |
| certificates | [Certificates](#certificates)| No | Object containing the path and filenames for the certificates to use in TLS communication. |
| authentication | [Authentication](#authentication) | No | Object containing the username and password for authentication with MQTT Broker and the server API. |
| publish | [Publish](#publish-and-consume) | Yes | Object containing the information about the MQTT topic to publish messages. |
| consume | [Consume](#publish-and-consume) | Yes | Object containing the information about the MQTT topic to read from. |
| api | [Api](#api) | Yes | Object containing the information about the API communication data. |

## Certificates

The certificates configuration object holds the TLS certificates path and filenames so we can use it to connect to the MQTT Broker and the server API using TLS.

This object has the following format.

| Key | Type | Required | Description |
| --- | ---- | -------- | ----------- |
| dir | string | Yes | Directory where the certificates are located. |
| root | string | Yes | Filename of the root certificate. |
| crt | string | Yes | Filename of the public key. |
| key | string | Yes | Filename of the private key. |

## Authentication

The authentication configuration holds the username and password for authenticating with the MQTT Broker and the server's API.

The object has the following format.

| Key | Type | Required | Description |
| --- | ---- | -------- | ----------- |
| username | string | Yes | Username used to authenticate with the MQTT Broker and the API. |
| password | string | Yes | Password used to authenticate with the MQTT Broker and the API. |

## Publish and Consume

The publish and Consume configuration object holds the information about the topic and MQTT Broker to publish messages in.

It has the following format.

| Key | Type | Required | Description |
| --- | ---- | -------- | ----------- |
| host | string | Yes | MQTT Broker address to publish messages in. |
| port | int | Yes | TCP port where the MQTT Broker is listening in. |
| topic | string | Yes | Topic name to publish/consume messages in/from. |
| qos | int [0-2] | Yes | QOS level used to publish/consume the messages in/from the topic. |
| interval | int | Yes | Number of milliseconds (real time) between publish/consume operations. |

## Api

The API configuration object holds the information to the server's API connection.

It has the following format.

| Key | Type | Required | Description |
| --- | ---- | -------- | ----------- |
| host | string | Yes | API host address. |
| port | int | Yes | TCP port where the API is listening to. |
| token | string | No | Last token received from the API. Do not change this value. |
