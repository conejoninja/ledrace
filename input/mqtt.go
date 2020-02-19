package input

import (
	"tinygo.org/x/drivers/net/mqtt"
)

type MQTTController struct {
	client  *mqtt.Client
	channel string
	pressed uint8
}

func NewMQTTController(client *mqtt.Client, channel string) *MQTTController {
	return &MQTTController{
		client:  client,
		channel: channel,
	}
}

func (c *MQTTController) Configure() {
	mc := *c.client
	token := mc.Subscribe(c.channel, 0, c.subHandler)
	token.Wait()
	if token.Error() != nil {
		println(token.Error().Error())
	}
}

func (c *MQTTController) subHandler(client mqtt.Client, msg mqtt.Message) {
	// Process msg.Payload() maybe?

	// increment pressed counter, and decrease it on Get() when it's read by the race's logic
	c.pressed++
}

func (c *MQTTController) Get() bool {
	if c.pressed > 0 {
		c.pressed--
		return true
	}
	return false
}

func (c *MQTTController) Reset() {
	c.pressed = 0
}
