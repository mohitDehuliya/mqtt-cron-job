# mqtt-sdk-go
Go language version of the MQTT SDK for SimplySmart

Configuring the Gateway
---------------------

Modify the sample.go file as follows:

1. Add appropriate access_key, username and password in GatewayConfig struct.
2. Add Key and Name of the thing whose data you want to send to MQTT wrapper.
3. Finally add the metrics name and its value in getData function, and send it to MQTT wrapper.

Implementing the Gateway
------------------------

You can then run example as follows:

go run sample.go

