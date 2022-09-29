# Eagle-One
The former Congress load testing tool now the LoRaWAN concentrator emulator.

## Parameters
Note that the Congress configuration is retrieved from `~/.congress` or 
environment variables. 

```
  -application-eui string
        Use existing application (-keep-application will be ignored)
  -congress-udp-port int
        Congress port (default 8000)
  -corrupt-mic int
        Percentage of packets generated that has a corrupt checksum. (default 5)
  -corrupt-payload int
        Percentage of packets generated that has a corrupt checksum.
  -devices int
        Number of devices. (default 10)
  -duplicate-message int
        Percentage of messages that will be duplicated. (default 2)
  -fancy-numerical-payload
        Generates non-insane numerical output in the form of a two bytes
  -frame-counter-errors int
        Frame counter errors (0-100) (default 5)
  -gateway-eui string
        Use existing gateway (--keep-gateway will be ignored)
  -keep-application
        Keep application when shutting down, don't remove it
  -keep-devices
        Keep devices when shutting down
  -keep-gateway
        Keep gateway when shutting down, don't delete it.
  -list-sent-messages
        Read back sent messages from Congress and list the contents
  -loglevel int
        Log level (0: Debug, 1: Info, 2: Warning: 3: Error) (default 1)
  -max-payload-size int
        ID offset for device EUI (default 222)
  -messages int
        Number of messages from each device. (default 10)
  -mode string
        Eagle One mode (interactive, batch) (default "batch")
  -transmission-delay int
        Delay (in milliseconds) between transmissions. (default 1000)
```