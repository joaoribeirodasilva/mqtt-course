# Mosquitto clients subscribe and publish

This document describes how to correctly use the mosquito command line clients in such a way that subscribers receive all the data not sent to them when they are offline.

The documentation for the subjects in this document can be found at:

* [mosquitto_sub man page](https://mosquitto.org/man/mosquitto_sub-1.html)
* [mosquitto_pub man page](https://mosquitto.org/man/mosquitto_pub-1.html)
* [mosquitto.conf man page](https://mosquitto.org/man/mosquitto-conf-5.html)

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
| -c | --disable-clean-session | don't delete old unreceived messages | No |

**Explanation:**

* We want to connect to host *localhost* at port *1883* where the broker is listening to.
* We want to read messages from topic *topictest* with a QOS level 2.
* We inform the id of the client making the request.
* We require the program to be verbose
* We inform the broker to do not delete unreceived messages that may be in the queue.

### Unsubscribe from a topic

```bash
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client1 -E
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client2 -E
```

**Parameters:**

| Short | Long | Description | Required |
| :- | :-    | :-           |  :-: |
| -h | --host | host address of the broker | No (if localhost) |
| -p | --port | port on the host where the broker is listening | No (if 1883) |
| -U | --unsubscribe | topic name to unsubscribe the client from | Yes |
| -t | --topic | topic name he client is subscribed to | Yes |
| -i | --id | client unique id. | Yes |
| -E | [none] | exit right after unsubscribed. | No |

**Explanation:**

* We want to connect to host *localhost* at port *1883* where the broker is listening to.
* We want to unsubscribe *client1* and *client2* from the topic *topictest*.
* We want to terminate the program right after unsubscribe.

### Max message retention

A message will be stored in the queue for a client according to the server configuration.

There are some configuration values in the *mosquito.conf* file that can affect the number of messages stored in the queue for clients, they are:

| Key | Description |
| --- | ----------- |
| max_queued_bytes | Sets the maximum number of bytes stored in te queue for each client. |
| max_queued_messages | Sets the maximum number of messages stored in te queue for each client. |

**NOTE:** If both *max_queued_bytes* and *max_queued_messages* are set, the limit is the lower of them. After any of this limits is reached the broker silently will not receive any more messages from the publisher (producers) util the limit drops.

## Exercise

Bellow you can find a little exercise to understand how the MQTT broker behaves.

### Just deliver last message

In some cases there is no need to receive all messages sent to the brokers. This is when is acceptable to loose some data.
In this case we don't need to subscribe nor publish with a QOS 1 or 2, but instead with QOS 0 (default).

To see this behavior in action just subscribe to a topic with the line bellow.

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -i client1 -E -v
```

We are not subscribing with QOS 2 anymore.

The let's publish some messages into the topic.

```bash
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 5,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.6, \"time\": 1699948605154305446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 6,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.7, \"time\": 1699948605154306446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 7,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.8, \"time\": 1699948605154307446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 8,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.9, \"time\": 1699948605154308446}"
```

Now let's read from the topic.

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -i client1 -v
```

And although we publish the messages using the *-q 2* parameter nothing comes out of our command.

Now let's unsubscribe both clients from the topic so we reset their subscriptions.

```bash
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client1 -E
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client2 -E
```

### Deliver all unreceived messages

In order to receive old messages we need to comply with the following steps:

1. The subscription must be made using the *-q 2* parameter for the client subscribing the topic.
2. The publication of the messages in the topic must also be made with the *-q 2* parameter.
3. The read of the message must be made with the *-q 2* and the *-c* parameter.

**NOTE**: If any read is made without the *-q 2* or the *-c* the old messages for the client will be deleted from the topic.

In order to receive any unreceived messages we must do the following.

Subscribe to the topic:

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client1 -E -v
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client2 -E -v
```

Publish into the topic:

```bash
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 5,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.6, \"time\": 1699948605154305446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 6,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.7, \"time\": 1699948605154306446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 7,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.8, \"time\": 1699948605154307446}"
mosquitto_pub -h localhost -p 1883 -t topictest -q 2 -m "{\"id\": 8,\"device\": \"45f3-3467-2c67-348f\",\"sensor\": \"temp\", \"value\": 15.9, \"time\": 1699948605154308446}"
```

Read from the topic:

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client1 -c -v
```

In another shel...

```bash
mosquitto_sub -h localhost -p 1883 -t topictest -q 2 -i client2 -c -v
```

Now let's unsubscribe both clients from the topic so we reset their subscriptions.

```bash
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client1 -E
mosquitto_sub -h localhost -p 1883 -U topictest -t topictest -i client2 -E
```

## Conclusion

* By default subscribers can only receive messages from a topic if they are connected.
* In order for a subscriber to get the unreceived messages from a topic it has to subscribe to the topic using QOS level level 1 or 2 and the messages need to be published with QOS 1 or 2 and the read again by the client using QOS 1 or 2 and the *disable clean session* option.
