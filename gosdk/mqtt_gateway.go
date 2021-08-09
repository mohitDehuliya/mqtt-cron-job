package gosdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	QoS0 = iota
	QoS1
	QoS2
)

var MqttClientId string

type Instruction struct {
	ThingKey  string                 `json:"thing_key,omitempty"`
	AlertKey  string                 `json:"alert_key,omitempty"`
	Inst      map[string]interface{} `json:"instruction,omitempty"`
	Timestamp int64                  `json:"time_stamp,omitempty"`
	AccessKey string                 `json:"access_key,omitempty"`
}

type MqttGateway struct {
	*GatewayConfig

	//Queue sync.WaitGroup // to wait on messages.
	Acks               chan string
	Instructions       chan Instruction
	InstructionHandler func(gw Gateway, ts int64, t *Thing, alertKey string, instruction map[string]interface{})
	client             MQTT.Client
	//RegisteredKeys       []string
	//RegisteredDeviceKeys []string
	//SubscribedTopics     []string
}

func NewClient(config *GatewayConfig) MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(config.Url())
	MqttClientId = config.Username + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	opts.SetClientID(MqttClientId)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)

	return MQTT.NewClient(opts)
}

func ConnectMqtt(config *GatewayConfig) (*MqttGateway, error) {
	client := NewClient(config)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Printf("Gateway connection: %v\n", client.IsConnected())
	gw := MqttGateway{
		GatewayConfig: config,
		Acks:          make(chan string),
		Instructions:  make(chan Instruction),
		client:        client,
	}

	go instructionWorker("InstructionWorker-Thread", &gw)

	var callback MQTT.MessageHandler = func(c MQTT.Client, msg MQTT.Message) {
		// Ack it
		gw.Acks <- string(msg.Payload())
	}

	http_ack := fmt.Sprintf("SimplySmart/%s/httpAck", MqttClientId)
	token = client.Subscribe(http_ack, 1, callback)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &gw, nil
}

func instructionWorker(threadName string, gw *MqttGateway) {
	fmt.Println("Started: ", threadName)
	for {
		ins := <-gw.Instructions
		if gw.InstructionHandler != nil {
			fmt.Println("Recieved new Instruction: ", ins)
			t := NewThing()
			t.Key = ins.ThingKey
			gw.InstructionHandler(gw, ins.Timestamp, t, ins.AlertKey, ins.Inst)
		} else {
			fmt.Println("Received instruction: ", ins, ". But, no handler set for executing it. It will be ignored")
		}
	}
}

func (gw *MqttGateway) IsConnected() bool {
	return gw.client.IsConnected()
}

func (gw *MqttGateway) GetConfig() *GatewayConfig {
	return gw.GatewayConfig
}

func (gw *MqttGateway) SubscribeForInstructions(t *Thing) error {
	var callback MQTT.MessageHandler = func(c MQTT.Client, msg MQTT.Message) {
		// Ack it
		gw.Acks <- string(msg.Payload())
	}

	alert_ack := fmt.Sprintf("SimplySmart/%s/thing/%s/alertsAck", gw.AccessKey, t.Key)
	token := gw.client.Subscribe(alert_ack, 1, callback)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := fmt.Sprintf("SimplySmart/%s/thing/%s/executeInstruction", gw.AccessKey, t.Key)
	fmt.Println("Sending -", topic)
	token = gw.client.Subscribe(topic, QoS2, func(c MQTT.Client, msg MQTT.Message) {
		rawPayload := msg.Payload()
		ins := Instruction{}
		buf := bytes.NewReader(rawPayload)
		dec := json.NewDecoder(buf)
		dec.Decode(&ins)
		gw.Instructions <- ins
	})

	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (gw *MqttGateway) ThingRegister(t *Thing) error {
	return gw.SubscribeForInstructions(t)
}

func (gw *MqttGateway) ThingEvent(eventData map[string]interface{}) error {
	topic := fmt.Sprintf("SimplySmart/%s/event", MqttClientId)
	eventData["access_key"] = gw.AccessKey
	return gw.sync_send_message(topic, eventData, QoS2)
}

func (gw *MqttGateway) SetInstructionHandle(instructionHandler func(gw Gateway, ts int64, t *Thing, alertKey string, instruction map[string]interface{})) error {
	gw.InstructionHandler = instructionHandler
	return nil
}

func (gw *MqttGateway) ThingHeartbeat(t *Thing, ts int64) error {
	eventData := CreateThingHeartbeat(t, ts)
	topic := fmt.Sprintf("SimplySmart/%s/heartbeat", MqttClientId)
	return gw.publish(topic, eventData, 0)
}

func (gw *MqttGateway) InstructionAck(t *Thing, alertKey, message string, alertLevel int, data map[string]interface{}) error {
	topic := fmt.Sprintf("SimplySmart/%s/thing/%s/alerts", gw.AccessKey, t.Key)
	payload := CreateInstructionAlert(alertKey, message, alertLevel, data, 0)
	payload["alert"].(map[string]interface{})["access_key"] = gw.AccessKey
	return gw.publish(topic, payload, 1)
}

func (gw *MqttGateway) Alert(t *Thing, alertMessage string, alertLevel int, alertData map[string]interface{}) error {
	topic := fmt.Sprintf("SimplySmart/%s/thing/%s/alerts", gw.AccessKey, t.Key)
	payload := CreateAlert(t, alertMessage, alertLevel, alertData, 0)
	payload["alert"].(map[string]interface{})["access_key"] = gw.AccessKey
	return gw.sync_send_message(topic, payload, 1)
}

func (gw *MqttGateway) sync_send_message(topic string, payload map[string]interface{}, qos byte) error {
	fmt.Printf("Sending - %s\n", topic)
	message, _ := json.Marshal(payload)
	if token := gw.client.Publish(topic, qos, false, message); token.Wait() && token.Error() != nil {
		// No point in panic here. Log this message and move on!
		fmt.Println("Error: ", token.Error())
	}

	/* Loop over the channel waiting for data */
	for value := range gw.Acks {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(value), &m)
		if err == nil && m["access_key"] == gw.AccessKey {
			var c, _ = m["http_code"].(float64)
			code := int(c)
			if code != 200 {
				return fmt.Errorf("there was an error with the request. Response code: %d", code)
			}
			return nil
		} else {
			// Ignore.put the data back in the channel
			fmt.Println("ignore: ", value)
			// gw.Acks <- value
		}
	}
	return nil
}

func (gw *MqttGateway) publish(topic string, payload map[string]interface{}, qos byte) error {
	fmt.Printf("Sending - %s\n", topic)
	message, _ := json.Marshal(payload)
	if token := gw.client.Publish(topic, qos, false, message); token.Wait() && token.Error() != nil {
		// No point in panic here. Log this message and move on!
		fmt.Println("Error: ", token.Error())
	}

	return nil
}
