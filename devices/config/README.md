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

## Sensors Object

## Communication Object
