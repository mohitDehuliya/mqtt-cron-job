package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/joshsoftware/mqtt-sdk-go/gosdk"
)

var wg sync.WaitGroup

func getData() map[string]interface{} {
	return map[string]interface{}{
		"pulsel":       rand.Intn(60),
		"pulseh":       rand.Intn(30),
		"meter_status": 0,
		"valve_status": 2,
	}
}

func sendEvent(gw gosdk.Gateway, t *gosdk.Thing) {
	// Send Event data
	data := getData()

	eventData := gosdk.CreateThingEvent(t, data, 0)
	retVal := gw.ThingEvent(eventData)
	if retVal == nil {
		fmt.Println("Send thing event: ", eventData)
	} else {
		fmt.Println("Unable to send thing event: ", retVal)
	}

}

func sendAlert(gw gosdk.Gateway, t *gosdk.Thing, alertLevel int, alertLevelStr string) error {
	data := map[string]interface{}{"tempearture": 120}
	if err := gw.Alert(t, "This is an example "+alertLevelStr+" alert using go SDK.", alertLevel, data); err != nil {
		fmt.Println("Unable to send Alert: ", err)
	}
	return nil
}

func sendAlerts(gw gosdk.Gateway, t *gosdk.Thing) error {
	sendAlert(gw, t, 0, "INFO")
	sendAlert(gw, t, 1, "WARNING")
	sendAlert(gw, t, 2, "ERROR")
	sendAlert(gw, t, 3, "CRITICAL")
	return nil
}

func handleInstruciton(gw gosdk.Gateway, ts int64, t *gosdk.Thing, alertKey string, instruction map[string]interface{}) {
	fmt.Println("Recieved Instruction Data for thing with key-device_key ", t.Key, ": ", instruction)
	alertMsg := fmt.Sprintf("Instruction with alert key %s executed", alertKey)
	if err := gw.InstructionAck(t, alertKey, alertMsg, 0, map[string]interface{}{"execution_status": "success"}); err == nil {
		fmt.Println("Sent instruction execution ACK back to datonis")
	} else {
		fmt.Println("Could not send instruction execution ACK back to datonis")
	}
}

func main() {

	file, err := os.Open("thing_key.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	rand.Seed(time.Now().UnixNano())
	fmt.Println("Here we go..")
	config := &gosdk.GatewayConfig{
		AccessKey: "78d9fa4teebt5add59ctb86e1a286477cb147392",
		Username:  "user",
		Password:  "user1",
	}

	config = gosdk.InitGatewayConfig(config)
	gw := gosdk.CreateGateway(config)

	for _, val := range records[0] {
		wg.Add(1)
		go launch_thing_test(gw, val)
	}

	wg.Wait()
}

func launch_thing_test(gw gosdk.Gateway, test_key string) {

	defer wg.Done()

	t := gosdk.NewThing()
	t.Name = "water_meter"
	t.Key = test_key
	if err := gw.ThingRegister(t); err != nil {
		fmt.Println("Unable to Register: ", err)
	}

	// Handle instructions if using the mqtt or mqtts protocol.
	if err := gw.SetInstructionHandle(handleInstruciton); err == nil {
		fmt.Println("Started handling the instruction.")
	} else {
		fmt.Println("Could not started handling the instruction: ", err)
	}

	// Send Alerts
	sendAlerts(gw, t)
	fmt.Println("----->")
	fmt.Println(t)
	fmt.Println("----->")
	// count := 0
	// for i := 0; i < 1000; i++ {
	// if count == 0 {
	// 	// Send a heartbeat.
	// 	gw.ThingHeartbeat(t, 0)
	// }
	sendEvent(gw, t)
	time.Sleep(time.Duration(1 * time.Second))
	//count++
	//count = count % 3
	//}

	//fmt.Println("Waiting...")
	//time.Sleep(time.Duration(1 * time.Second))
	//fmt.Println("Bye Bye")

}
