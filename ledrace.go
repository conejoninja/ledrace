package ledrace

import (
	"machine"
	"time"

	"image/color"

	"github.com/tinygo-org/drivers/ws2812"
)

const ACCELERATION = 0.2
const TRACKLENGHT = 240
const FRICTION = 0.015
const GRAVITY = 0.003

const PLAYERS = 2
const LAPS = 3

type Track struct {
	gravity []uint8
	leds    []color.RGBA
	ws      ws2812.Device
}

type Player struct {
	button   machine.GPIO
	led      machine.GPIO
	pressed  bool
	speed    float32
	position float32
	color    color.RGBA
}

var players [2]Player
var track Track

var black = color.RGBA{0, 0, 0, 0}
var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}
var orange = color.RGBA{255, 255, 0, 255}

func main() {
	players[0].button = machine.GPIO{machine.A3}
	players[0].button.Configure(machine.GPIOConfig{Mode: machine.GPIO_INPUT})
	players[0].led = machine.GPIO{machine.A5}
	players[0].led.Configure(machine.GPIOConfig{Mode: machine.GPIO_OUTPUT})

	players[1].button = machine.GPIO{machine.A2}
	players[1].button.Configure(machine.GPIOConfig{Mode: machine.GPIO_INPUT})
	players[1].led = machine.GPIO{machine.A4}
	players[1].led.Configure(machine.GPIOConfig{Mode: machine.GPIO_OUTPUT})

	neo := machine.GPIO{machine.A1}
	neo.Configure(machine.GPIOConfig{Mode: machine.GPIO_OUTPUT})
	track.ws = ws2812.New(neo)

	track.leds = make([]color.RGBA, TRACKLENGHT)
	track.gravity = make([]uint8, TRACKLENGHT)
	for i := 0; i < TRACKLENGHT; i++ {
		track.gravity[i] = 127
	}

	players[0].color = red
	players[1].color = green

	startRace()

	for {
		for p := uint8(0); p < PLAYERS; p++ {
			pressed := players[p].button.Get()
			if players[p].pressed && !pressed {
				players[p].pressed = false
				players[p].speed += ACCELERATION
				players[p].led.High()
			} else if !players[p].pressed && pressed {
				players[p].pressed = true
				players[p].led.Low()
			}

			gravity := track.gravity[uint16(players[p].position)%TRACKLENGHT]
			if gravity < 127 {
				players[p].speed -= GRAVITY * float32(127-gravity)
			}
			if gravity > 127 {
				players[p].speed += GRAVITY * float32(gravity-127)
			}

			players[p].speed -= players[p].speed * FRICTION
			players[p].position += players[p].speed

			if players[p].position > LAPS*TRACKLENGHT {
				finishRace(p)
				break
			}
		}

		paintTrack()
		time.Sleep(10 * time.Millisecond)
	}
}

func paintTrack() {
	for i := 0; i < TRACKLENGHT; i++ {
		track.leds[i] = black
	}
	for p := uint8(0); p < PLAYERS; p++ {
		track.leds[uint32(players[p].position)%TRACKLENGHT] = players[p].color
	}
	track.ws.WriteColors(track.leds)
}

func startRace() {
	resetPlayers()
	time.Sleep(2 * time.Second)
	for i := 0; i < TRACKLENGHT; i++ {
		track.leds[i] = black
	}
	track.ws.WriteColors(track.leds)
	time.Sleep(1 * time.Second)
	track.leds[12] = green
	track.leds[11] = green
	track.ws.WriteColors(track.leds)
	time.Sleep(1 * time.Second)
	track.leds[12] = black
	track.leds[11] = black
	track.leds[10] = orange
	track.leds[9] = orange
	track.ws.WriteColors(track.leds)
	time.Sleep(1 * time.Second)
	track.leds[10] = black
	track.leds[9] = black
	track.leds[8] = red
	track.leds[7] = red
	track.ws.WriteColors(track.leds)
	time.Sleep(1 * time.Second)
}

func finishRace(winner uint8) {
	resetPlayers()
	for i := 0; i < 6; i++ {
		for i := 0; i < TRACKLENGHT; i++ {
			track.leds[i] = players[winner].color
		}
		track.ws.WriteColors(track.leds)
		time.Sleep(300 * time.Millisecond)
		for i := 0; i < TRACKLENGHT; i++ {
			track.leds[i] = black
		}
		track.ws.WriteColors(track.leds)
		time.Sleep(300 * time.Millisecond)
	}
	startRace()
}

func resetPlayers() {
	for p := uint8(0); p < PLAYERS; p++ {
		players[p].speed = 0
		players[p].position = 0
		players[p].pressed = false
	}
}
