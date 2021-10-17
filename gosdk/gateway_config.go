package gosdk

import (
	"fmt"

	constant "github.com/joshsoftware/mqtt-sdk-go/common/constants"
)

type GatewayConfig struct {
	AccessKey string
	Username  string
	Password  string
	Host      string
	Port      int32
}

func InitGatewayConfig(c *GatewayConfig) *GatewayConfig {
	if c.Host == "" {
		c.Host = constant.HostAddress
	}
	if c.Port == 0 {
		c.Port = constant.MQTTPort
	}

	return c
}

func (c *GatewayConfig) Url() string {
	url := ""
	url = fmt.Sprintf("tcp://%s:%d", c.Host, c.Port)
	return url
}
