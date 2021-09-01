package gosdk

import (
	"time"
)

func GetCurrentTS() int64 {
	return (time.Now().Unix())
}

func CreateThingHeartbeat(t *Thing, ts int64) map[string]interface{} {
	if ts == 0 {
		ts = GetCurrentTS()
	}
	data := map[string]interface{}{
		"time_stamp": ts,
		"thing_key":  t.Key,
	}

	return data
}

func CreateThingEvent(t *Thing, data map[string]interface{}, ts int64) map[string]interface{} {
	if ts == 0 {
		ts = GetCurrentTS()
	}

	eventData := map[string]interface{}{
		"time_stamp": ts,
	}

	if data != nil {
		eventData["data"] = data
	}

	if t.Key != "" {
		eventData["thing_key"] = t.Key
	}

	return eventData
}

func CreateAlert(t *Thing, alertMessage string, alertLevel int, alertData map[string]interface{}, ts int64) map[string]interface{} {
	if ts == 0 {
		ts = GetCurrentTS()
	}
	payload := map[string]interface{}{
		"alert": map[string]interface{}{
			"data":        alertData,
			"message":     alertMessage,
			"alert_level": alertLevel,
			"thing_key":   t.Key,
			"time_stamp":  ts,
		},
	}

	return payload
}

func CreateInstructionAlert(alertKey string, alertMessage string, alertLevel int, alertData map[string]interface{}, ts int64) map[string]interface{} {
	if ts == 0 {
		ts = GetCurrentTS()
	}
	payload := map[string]interface{}{
		"alert": map[string]interface{}{
			"data":        alertData,
			"message":     alertMessage,
			"alert_level": alertLevel,
			"alert_key":   alertKey,
			"time_stamp":  ts,
		},
	}

	return payload
}
