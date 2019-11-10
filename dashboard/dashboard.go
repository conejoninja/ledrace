package dashboard

import (
	"encoding/json"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Info struct {
	Info        string   `json:"info"`
	TrackLength int      `json:"track_length"`
	NumPlayers  int      `json:"num_players"`
	Laps        int      `json:"laps"`
	Players     []Player `json:"players"`
}

type Player struct {
	ID       uint8   `json:"id"`
	Speed    float32 `json:speed"`
	Position float32 `json:position"`
	Color    string  `json:"color"`
	Loop     uint8   `json:loop"`
}

var subscriptions map[string]bool
var token mqtt.Token

func Start(mqttclient mqtt.Client) {
	c := mqttclient

	subscriptions = make(map[string]bool)

	if token = c.Subscribe("demo-track", 0, trackHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}

var trackHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var info Info
	fmt.Println("RAW MSG", msg, string(msg.Payload()))
	err := json.Unmarshal(msg.Payload(), &info)
	if err == nil {
		fmt.Println("INFO", info)
	} else {
		fmt.Println("ERROR", err)
	}
}
