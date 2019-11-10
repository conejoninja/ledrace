package main

import (
	"fmt"
	"time"

	"github.com/conejoninja/ledrace/dashboard"
	"github.com/eclipse/paho.mqtt.golang"

	config "../../config/local"
)

// MQTT
var token mqtt.Token
var mqttclient mqtt.Client

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MQTTServer()).SetClientID("dashboard-cli")
	if config.MQTTUser() != "" {
		opts.SetUsername(config.MQTTUser())
	}
	if config.MQTTPassword() != "" {
		opts.SetPassword(config.MQTTPassword())
	}

	fmt.Println("Connecting to MQTT...")
	mqttclient = mqtt.NewClient(opts)
	if token = mqttclient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}

	dashboard.Start(mqttclient)

	for {
		fmt.Println(time.Now(), "Still alive")
		time.Sleep(5 * time.Minute)
	}
}
