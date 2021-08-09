package gosdk

type ThingsPayload struct {
	Data      map[string]interface{} `json:"data"`
	ThingKey  string                 `json:"thing_key"`
	AccessKey string                 `json:"access_key"`
	TimeStamp uint64                 `json:"time_stamp"`
}
