# Mosquitto clients

This document describes how to correctly use the mosquito command line clients in such a way that subscribers receive all the data not sent to them when they are offline.

## Flow

1. Subscribe to one or more topics using *mosquito_sub*.
2. Publish into the topic that was already subscribed using *mosquito_pub*.
3. Read from the topic using *mosquito_sub*.

## Commands

The following sections show how the commands and their options must be set in order for everything to be right.

### Environmental assumptions

* The *mosquitto broker* is located at *localhost* and listening on port *1883*.
* The relevant topic is *topictest*.
* We will use 2 subscribers one with the id *client1* and a second with the id *client2*.

### Subscribing to a topic

To subscribe *client1* and *client2* to the topic *topictest* just run the lines bellow.

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client1 -E -v 
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client2 -E -v 
```

**Parameters:**
| Short | Long | Description | Required |
| :- | :-    | :-           |  :-: |
| -h | --host | host address of the broker | No (if localhost) |
| -p | --port | port on the host where the broker is listening | No (if 1883) |
| -t | --topic | topic name to subscribe the client to | Yes |
| -q | --qos | QOS service level. The number 2 tells the broker that the subscriber would like QOS level 2 (ensures at least one delivery per subscriber). | No |
| -i | --id | client unique id. | Yes |
| -E | [none] | exit right after subscribing | No |

**Explanation:**

* We want to connect to host *localhost* at port *1883* where the broker is listening to.
* We want to subscribe to the topic *topictest* with a QOS level 2 so it stores the messages for our client if it's offline. If we don't set this value to 2 then the message delivery to our client is not guaranteed if the client is offline.
* Then we pass to the broker which client is requesting the subscription.



### Publishing to a topic

To publish messages into a topic just run the lines bellow.

```bash
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 5,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.6, \"time\": 1699948605154305446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 6,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.7, \"time\": 1699948605154306446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 7,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.8, \"time\": 1699948605154307446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 8,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.9, \"time\": 1699948605154308446}"
```

**Parameters:**
| Short | Long | Description | Required |
| :- | :-    | :-           |  :-: |
| -h | --host | host address of the broker | No (if localhost) |
| -p | --port | port on the host where the broker is listening | No (if 1883) |
| -t | --topic | topic name to subscribe the client to | Yes |
| -q | --qos | QOS service level. The number 2 tells the broker that the message is of QOS level 2 (ensures at least one delivery per subscriber). | No |
| -m | --message | Message string content | Yes |

**Explanation:**

* We want to connect to host *localhost* at port *1883* where the broker is listening to.
* We want to subscribe to the topic *topictest* with a QOS level 2 so it stores the messages for our clients if they are offline. If we don't set this value to 2 then the message delivery to our clients is not guaranteed for the clients that are offline.
* Then we send the message string (JSON in this case).

### Reading from a topic

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client1 -v -c
```

In another shell...

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client2 -v -c
```

To exit each program type CTRL+C in each correspondent shell.

**Parameters:**
| Short | Long | Description | Required |
| :- | :-    | :-           |  :-: |
| -h | --host | host address of the broker | No (if localhost) |
| -p | --port | port on the host where the broker is listening | No (if 1883) |
| -t | --topic | topic name to subscribe the client to | Yes |
| -q | --qos | QOS service level. The number 2 tells the broker that the message is of QOS level 2 (ensures at least one delivery per subscriber). | No |
| -i | --id | client unique id. | Yes |
| -v | --verbose | tells the client to print to the screen | No |
| -c | --disable-clean-session | receive all messages not yet received | No |

**Explanation:**

[TODO:]

### Unsubscribe from a topic

```bash
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client1 -E
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client2 -E
```

**Parameters:**

[TODO:]

**Explanation:**

[TODO:]

### Max message retention

A message will be stored in the queue for a client according to the server configuration.

[TODO: configuration key in mosquitto.conf]

## Exercise

[TODO:]

### Just deliver last message

[TODO:]

### Deliver all unreceived messages

[TODO:]

## Sumary

[TODO:]
