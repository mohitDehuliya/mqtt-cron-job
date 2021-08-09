package gosdk

type Gateway interface {
	GetConfig() *GatewayConfig

	// Registers the specified thing with Datonis
	ThingRegister(t *Thing) error

	// Sends a Thing Data Packet (event) to Datonis
	ThingEvent(eventData map[string]interface{}) error

	// Sends a Heart Beat message to Datonis indicating that this thing is alive
	ThingHeartbeat(t *Thing, ts int64) error

	// Sends the Alerts to Datonis
	Alert(t *Thing, message string, alertLevel int, data map[string]interface{}) error

	// Set the instruction handler to handle the instruciton received
	SetInstructionHandle(instructionHandler func(gw Gateway, ts int64, t *Thing, alertKey string, instruction map[string]interface{})) error

	// Sends the instruction acknowledgement to datonis through MQTT protocol.
	InstructionAck(t *Thing, alertKey, message string, alertLevel int, data map[string]interface{}) error
}

func CreateGateway(c *GatewayConfig) Gateway {
	gw, _ := ConnectMqtt(c)
	return gw
}
