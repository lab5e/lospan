# lospan - a LoRaWAN library

This library is based on the Congress LoRaWAN server by Telenor Digital (see the NOTICE file) and the goal is
"LoRaWAN in Box"-like features in library form. 

Some time in the future the application/device/gateway management parts and networking parts will be
split into the corresponding network server/application server. The application server doesn't do anything
except manage the encryption keys for the devices.

Bigger installations will have multiple gateways/concentrators that will report to one or more application
server processes. The application servers may or may not have clustering (with mofunk magic).