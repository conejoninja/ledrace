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
	enabled bool
}

var black = color.RGBA{0, 0, 0, 0}
var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}
var orange = color.RGBA{255, 255, 0, 255}
var white = color.RGBA{255, 255, 255, 255}
var blue = color.RGBA{0, 0, 255, 255}

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
	t.enabled = true
	t.payload = []byte("boot sequence")
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

	go t.sendLoop()
	time.Sleep(4 * time.Second)

	for i := 0; i < 10; i++ {
		ledsColor[i] = green
	}
	t.leds.WriteColors(ledsColor)

	return &t
}

func NewDisabled(leds *ws2812.Device) *Telemetry {
	var t Telemetry
	t.leds = leds
	t.payload = []byte("boot sequence")
	ledsColor := make([]color.RGBA, 10)
	for i := 0; i < 10; i++ {
		ledsColor[i] = blue
	}
	t.leds.WriteColors(ledsColor)
	time.Sleep(3000 * time.Millisecond)

	return &t
}

func (t *Telemetry) Enabled() bool {
	return t.enabled
}

func (t *Telemetry) sendLoop() {
	retries := uint8(0)
	//ledsColor := make([]color.RGBA, 10)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.MQTTServer()).SetClientID(config.DeviceName())
	if config.MQTTUser() != "" {
		opts.SetUsername(config.MQTTUser())
	}
	if config.MQTTPassword() != "" {
		opts.SetPassword(config.MQTTPassword())
	}

	println("Connecting to MQTT...")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		println(token.Error().Error(), "NOT CONNECTED TO MQTT :(")
	}
	var token mqtt.Token

	for {
		if t.enabled {
			if retries == 0 {
				println("Publishing MQTT message...", string(t.payload))
				token = client.Publish(config.TrackChannel(), 0, false, t.payload)
				token.Wait()
			}
			if retries > 0 || token.Error() != nil {
				if retries < 10 {
					token = client.Connect()
					if token.Wait() && token.Error() != nil {
						retries++
						println("NOT CONNECTED TO MQTT (sendLoop)")
					} else {
						retries = 0
					}
				} else {
					t.enabled = false
				}
			}
			t.payload = []byte("none")
			time.Sleep(800 * time.Millisecond)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}

func (t *Telemetry) Send(payload []byte) {
	t.payload = payload
}

func (t *Telemetry) Enable() {
	t.enabled = true
	time.Sleep(1 * time.Second)
}

func (t *Telemetry) Disable() {
	t.enabled = false
	time.Sleep(1 * time.Second)
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
