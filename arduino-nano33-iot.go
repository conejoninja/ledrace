//go:build arduino_nano33
// +build arduino_nano33

package main

import (
	"machine"

	"tinygo.org/x/drivers/buzzer"
	"tinygo.org/x/drivers/ws2812"
)

func configureHardware() {
	bzrPin := machine.A0
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	bzr = buzzer.New(bzrPin)

	players[0].button = machine.A5
	players[0].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[0].led = machine.A6
	players[0].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[0].color = red

	players[1].button = machine.A3
	players[1].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[1].led = machine.A4
	players[1].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[1].color = green

	players[2].button = machine.A1
	players[2].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[2].led = machine.A2
	players[2].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[2].color = yellow

	players[3].button = machine.D12
	players[3].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[3].led = machine.D11
	players[3].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[3].color = blue

	neo := machine.D2
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	track.ws = ws2812.New(neo)
}
