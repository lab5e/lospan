# lospan - a LoRaWAN library

This library is based on the (Congress LoRaWAN server by Telenor Digital)[https://github.com/ExploratoryEngineering/congress].
(see the NOTICE file) 

The server is stripped down a bit to make a LoRaWAN server as a library. 

## Testing

Build with `make`

Launch the service itself with 

```shell
bin/congress --lora-connection-string=lora.db
``` 

then run the device emulator with 

```shell
bin/eagle-one --mode=create
```

You should see 10 devices send 10 messages each to the service.

Use `bin/lc` to interact with the gRPC API. This client is for development and testing only so expect sharp
edges.
