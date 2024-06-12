//go:build waveshare_rp2040_zero
// +build waveshare_rp2040_zero

package main

import (
	"machine"

	"tinygo.org/x/drivers/buzzer"
	"tinygo.org/x/drivers/ws2812"
)

func configureHardware() {
	bzrPin := machine.D4
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	bzr = buzzer.New(bzrPin)

	players[0].button = machine.D7
	players[0].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[0].led = machine.D8
	players[0].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[0].color = red

	players[1].button = machine.D14
	players[1].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[1].led = machine.D15
	players[1].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[1].color = green

	players[2].button = machine.D26
	players[2].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[2].led = machine.D27
	players[2].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[2].color = yellow

	players[3].button = machine.D28
	players[3].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[3].led = machine.D29
	players[3].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[3].color = blue

	neo := machine.D2
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	track.ws = ws2812.New(neo)
}
