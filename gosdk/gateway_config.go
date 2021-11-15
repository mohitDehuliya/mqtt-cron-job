package gosdk

import "fmt"

type GatewayConfig struct {
	AccessKey string
	Username  string
	Password  string
	Host      string
	Port      int32
}

func InitGatewayConfig(c *GatewayConfig) *GatewayConfig {
	if c.Host == "" {
		c.Host = "65.0.106.100"
	}
	if c.Port == 0 {
		c.Port = 1883
	}

	return c
}

func (c *GatewayConfig) Url() string {
	url := ""
	url = fmt.Sprintf("tcp://%s:%d", c.Host, c.Port)
	return url
}
