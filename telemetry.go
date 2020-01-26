// This is a sensor station that uses a ESP8266 or ESP32 running on the device UART1.
// It creates a UDP connection you can use to get info to/from your computer via the microcontroller.
//
// In other words:
// Your computer <--> UART0 <--> MCU <--> UART1 <--> ESP8266
//
package main

import (
	"image/color"
	"machine"
	"time"

	config "./config/local"
	"tinygo.org/x/drivers/espat"
	"tinygo.org/x/drivers/net"
)

type Telemetry struct {
	enabled bool
	payload []byte
	conn    *net.UDPSerialConn
}

var telemetry Telemetry

var (
	uart = machine.UART1
	tx   = machine.PA22
	rx   = machine.PA23

	adaptor *espat.Device
)

var orange = color.RGBA{255, 255, 0, 255}
var white = color.RGBA{255, 255, 255, 255}

func InitTelemetry() {
	ledsColor := make([]color.RGBA, 10)
	for i := 0; i < 10; i++ {
		ledsColor[i] = orange
	}
	track.ws.WriteColors(ledsColor)
	time.Sleep(3 * time.Second)
	uart.Configure(machine.UARTConfig{TX: tx, RX: rx})

	// Init esp8266/esp32
	adaptor = espat.New(uart)
	adaptor.Configure()

	// first check if connected
	if connectToESP() {
		println("Connected to wifi adaptor.")
		for i := 0; i < 10; i++ {
			ledsColor[i] = white
		}
		track.ws.WriteColors(ledsColor)
		adaptor.Echo(false)

		connectToAP()
	} else {
		println("")
		for i := 0; i < 10; i++ {
			ledsColor[i] = red
		}
		track.ws.WriteColors(ledsColor)
		println("Unable to connect to wifi adaptor.")
		return
	}

	for i := 0; i < 10; i++ {
		ledsColor[i] = green
	}
	track.ws.WriteColors(ledsColor)

	// now make TCP connection
	ip := net.ParseIP(config.ServerIP())
	raddr := &net.UDPAddr{IP: ip, Port: config.ServerpPort()}
	laddr := &net.UDPAddr{Port: config.ServerpPort()}

	println("Dialing UDP connection...")
	var err error
	telemetry.conn, err = net.DialUDP("udp", laddr, raddr)
	if err != nil {
		println(err.Error())
		for i := 0; i < 10; i++ {
			ledsColor[i] = red
		}
		return
	}

	telemetry.enabled = true
	telemetry.payload = []byte("boot")
	go telemetry.sendLoop()

}

func (t *Telemetry) sendLoop() {
	for {
		if t.enabled {
			println("Sending data...")
			telemetry.conn.Write(t.payload)
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
	println("Connecting to wifi network '" + config.SSID() + "'")

	adaptor.SetWifiMode(espat.WifiModeClient)
	adaptor.ConnectToAP(config.SSID(), config.Password(), 10)

	println("Connected.")
	ip, err := adaptor.GetClientIP()
	if err != nil {
		println(err.Error())
	}

	println(ip)
}
