package telemetry

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/espat"
	"tinygo.org/x/drivers/net/mqtt"
	"tinygo.org/x/drivers/ws2812"

	config "../config/local"
)

type Telemetry struct {
	c       mqtt.Client
	leds    *ws2812.Device
	payload []byte
}

var black = color.RGBA{0, 0, 0, 0}
var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}
var orange = color.RGBA{255, 255, 0, 255}
var white = color.RGBA{255, 255, 255, 255}

var (
	uart = machine.UART1
	tx   = machine.PA22
	rx   = machine.PA23

	console = machine.UART0

	adaptor *espat.Device
)

func New(leds *ws2812.Device) *Telemetry {
	var t Telemetry
	t.leds = leds
	time.Sleep(3000 * time.Millisecond)

	uart.Configure(machine.UARTConfig{TX: tx, RX: rx})

	// Init esp8266/esp32
	adaptor = espat.New(uart)
	adaptor.Configure()

	ledsColor := make([]color.RGBA, 10)
	for i := 0; i < 10; i++ {
		ledsColor[i] = orange
	}
	t.leds.WriteColors(ledsColor)

	// first check if connected
	if connectToESP() {
		println("Connected to wifi adaptor.")
		for i := 0; i < 10; i++ {
			ledsColor[i] = white
		}
		t.leds.WriteColors(ledsColor)
		adaptor.Echo(false)

		connectToAP()
	} else {
		for i := 0; i < 10; i++ {
			ledsColor[i] = red
		}
		t.leds.WriteColors(ledsColor)
		println("Unable to connect to wifi adaptor.")
		return &t
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MQTTServer()).SetClientID(config.DeviceName())
	if config.MQTTUser() != "" {
		opts.SetUsername(config.MQTTUser())
	}
	if config.MQTTPassword() != "" {
		opts.SetPassword(config.MQTTPassword())
	}

	println("Connecting to MQTT...")
	t.c = mqtt.NewClient(opts)
	if token := t.c.Connect(); token.Wait() && token.Error() != nil {
		for i := 0; i < 10; i++ {
			ledsColor[i] = red
		}
		t.leds.WriteColors(ledsColor)
		println(token.Error().Error())
	}

	for i := 0; i < 10; i++ {
		ledsColor[i] = green
	}
	t.leds.WriteColors(ledsColor)



	return &t
}

func (t *Telemetry) sendLoop() {
	for{
		if t.payload!=nil {
			token := t.c.Publish(config.TrackChannel(), 0, false, t.payload)
			token.Wait()
			if token.Error() == nil {
				t.payload = nil
			}
		}
		time.Sleep(1000*time.Millisecond)
	}
}

func (t *Telemetry) Send(payload []byte) {
	t.payload = payload
}

func (t *Telemetry) Disconnect() {
	t.c.Disconnect(100)
}

// connect to ESP8266/ESP32
func connectToESP() bool {
	for i := 0; i < 5; i++ {
		println("Connecting to wifi adaptor...")
		if adaptor.Connected() {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

// connect to access point
func connectToAP() {
	println("Connecting to wifi network...")

	adaptor.SetWifiMode(espat.WifiModeClient)
	adaptor.ConnectToAP(config.SSID(), config.Password(), 10)

	println("Connected.")
	println(adaptor.GetClientIP())
}
