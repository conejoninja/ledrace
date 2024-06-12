package main

import (
	"image/color"
	"machine"
	"strconv"
	"time"

	"tinygo.org/x/drivers/buzzer"
	"tinygo.org/x/drivers/ws2812"
)

const ACCELERATION = 0.2
const TRACKLENGHT = 300
const FRICTION = 0.015
const GRAVITY = 0.003

const PLAYERS = 4
const LAPS = 3

type Track struct {
	gravity []uint8
	leds    []color.RGBA
	ws      ws2812.Device
}

type Player struct {
	button   machine.Pin
	led      machine.Pin
	pressed  bool
	speed    float32
	position float32
	color    color.RGBA
	loop     uint8
}

var players [4]Player
var track Track
var bzr buzzer.Device
var idleTime time.Time
var activity bool

var black = color.RGBA{0, 0, 0, 0}
var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}
var yellow = color.RGBA{255, 255, 0, 255}
var blue = color.RGBA{0, 0, 255, 255}

func main() {
	configureHardware()

	track.leds = make([]color.RGBA, TRACKLENGHT)
	track.gravity = make([]uint8, TRACKLENGHT)
	for i := 0; i < TRACKLENGHT; i++ {
		track.gravity[i] = 127
	}

	setRamp(12, 90, 100, 110)

	// If at start, player 0 is pressed, enter configuration mode
	if !players[0].button.Get() {
		showGravity()
	}

	startRace()

	for {
		activity = false
		for p := uint8(0); p < PLAYERS; p++ {
			activity = activity || getPlayerInput(p)

			gravity := track.gravity[uint16(players[p].position)%TRACKLENGHT]
			if gravity < 127 {
				players[p].speed -= GRAVITY * float32(127-gravity)
			}
			if gravity > 127 {
				players[p].speed += GRAVITY * float32(gravity-127)
			}

			players[p].speed -= players[p].speed * FRICTION
			players[p].position += players[p].speed

			if players[p].position > TRACKLENGHT*float32(players[p].loop) {
				players[p].loop++
			}

		}

		paintTrack()

		maxPosition := players[0].position
		winner := uint8(0)
		playerStr := ""
		for p := uint8(0); p < PLAYERS; p++ {
			if players[p].position > maxPosition {
				maxPosition = players[p].position
				winner = p
			}

			playerStr += strconv.Itoa(int(1000*players[p].speed)) + "," + strconv.Itoa(int(players[p].position)) + "," + strconv.Itoa(int(players[p].loop))
			if p != PLAYERS-1 {
				playerStr += "|"
			}
		}

		if maxPosition > LAPS*TRACKLENGHT {
			time.Sleep(1500 * time.Millisecond)
			finishRace(winner)
		}

		if activity {
			idleTime = time.Now()
		} else if time.Since(idleTime) > 30*time.Second {
			idleRace()
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func getPlayerInput(p uint8) bool {
	pressed := players[p].button.Get()
	if players[p].pressed && !pressed {
		players[p].pressed = false
		players[p].speed += ACCELERATION
		players[p].led.High()
		return true
	} else if !players[p].pressed && pressed {
		players[p].pressed = true
		players[p].led.Low()
	}
	return false
}

func paintTrack() {
	for i := 0; i < TRACKLENGHT; i++ {
		track.leds[i] = black
	}
	for p := uint8(0); p < PLAYERS; p++ {
		for l := uint8(0); l < players[p].loop; l++ {
			track.leds[(uint32(players[p].position)-uint32(l))%TRACKLENGHT] = players[p].color
		}
	}
	track.ws.WriteColors(track.leds)
}

func startRace() {
	resetPlayers()

	idleTime = time.Now()
	time.Sleep(2 * time.Second)
	for i := 0; i < TRACKLENGHT; i++ {
		track.leds[i] = black
	}
	track.ws.WriteColors(track.leds)
	time.Sleep(1 * time.Second)
	track.leds[12] = red
	track.leds[11] = red
	track.ws.WriteColors(track.leds)
	bzr.Tone(buzzer.E4, buzzer.Quarter)
	time.Sleep(400 * time.Millisecond)
	track.leds[12] = black
	track.leds[11] = black
	track.leds[10] = yellow
	track.leds[9] = yellow
	track.ws.WriteColors(track.leds)
	bzr.Tone(buzzer.E4, buzzer.Quarter)
	time.Sleep(400 * time.Millisecond)
	track.leds[10] = black
	track.leds[9] = black
	track.leds[8] = green
	track.leds[7] = green
	track.ws.WriteColors(track.leds)
	bzr.Tone(buzzer.B4, buzzer.Quarter)
	time.Sleep(400 * time.Millisecond)
}

func finishRace(winner uint8) {
	resetPlayers()

	for k := 0; k < 6; k++ {
		for i := 0; i < TRACKLENGHT; i++ {
			track.leds[i] = players[winner].color
		}
		track.ws.WriteColors(track.leds)
		time.Sleep(300 * time.Millisecond)
		if k == 1 {
			// winning melody
			players[winner].led.High()
			bzr.Tone(buzzer.C4, 0.25)
			time.Sleep(100 * time.Millisecond)
			players[winner].led.Low()
			bzr.Tone(buzzer.C4, 0.25)
			time.Sleep(100 * time.Millisecond)
			players[winner].led.Low()
			bzr.Tone(buzzer.C4, 0.25)
			time.Sleep(100 * time.Millisecond)
			players[winner].led.High()
			bzr.Tone(buzzer.C4, 0.25)
			time.Sleep(100 * time.Millisecond)
			players[winner].led.Low()
			bzr.Tone(buzzer.G3, 0.5)
			time.Sleep(200 * time.Millisecond)
			players[winner].led.High()
			bzr.Tone(buzzer.A3, 0.5)
			time.Sleep(200 * time.Millisecond)
			players[winner].led.Low()
			bzr.Tone(buzzer.C4, 0.25)
			time.Sleep(100 * time.Millisecond)
			players[winner].led.High()
			bzr.Tone(buzzer.A3, 0.25)
			time.Sleep(100 * time.Millisecond)
			players[winner].led.Low()
			bzr.Tone(buzzer.C4, 0.5)
			time.Sleep(100 * time.Millisecond)
			players[winner].led.High()
		}
		for i := 0; i < TRACKLENGHT; i++ {
			track.leds[i] = black
		}
		track.ws.WriteColors(track.leds)
		time.Sleep(300 * time.Millisecond)
	}
	players[winner].led.Low()
	startRace()
}

func resetPlayers() {
	for p := uint8(0); p < PLAYERS; p++ {
		players[p].speed = 0
		players[p].position = 0
		players[p].pressed = false
		players[p].loop = 0
		players[p].led.Low()
	}
}

func idleRace() {
	var k uint16
	activity = false
	for {
		for i := 0; i < TRACKLENGHT; i++ {
			track.leds[i] = getRainbowRGB(uint8((i*256)/TRACKLENGHT) + uint8(k))
		}
		track.ws.WriteColors(track.leds)
		k = (k + 1) % 255

		for p := uint8(0); p < PLAYERS; p++ {
			if getPlayerInput(p) {
				activity = true
				break
			}
		}
		if activity {
			break
		}

		time.Sleep(16 * time.Millisecond)
	}
	startRace()
}

func getRainbowRGB(i uint8) color.RGBA {
	if i < 85 {
		return color.RGBA{i * 3, 255 - i*3, 0, 255}
	} else if i < 170 {
		i -= 85
		return color.RGBA{255 - i*3, 0, i * 3, 255}
	}
	i -= 170
	return color.RGBA{0, i * 3, 255 - i*3, 255}
}

func setRamp(height uint8, start uint32, middle uint32, end uint32) {
	for i := uint32(0); i < (middle - start); i++ {
		track.gravity[start+i] = uint8(127 - float32(i)*float32(height)/float32(middle-start))
	}
	track.gravity[middle] = 127
	for i := uint32(0); i < (end - middle); i++ {
		track.gravity[middle+i+1] = uint8(127 + float32(height)*(1-float32(i)/float32(middle-start)))
	}
}

func showGravity() {
	for i := 0; i < TRACKLENGHT; i++ {
		grav := track.gravity[i]
		r := uint8(2)
		g := uint8(2)
		if grav < 127 {
			r = 2 * (127 - grav)
		}
		if grav > 127 {
			g = (grav - 127) * 2
		}
		track.leds[i] = color.RGBA{r, g, 0, 255}
	}
	track.ws.WriteColors(track.leds)
	longpress := 0
	onoff := false
	for {
		if !players[0].button.Get() {
			longpress++
			if longpress > 20 {
				break
			}
		} else {
			longpress = 0
		}
		if onoff {
			players[0].led.Low()
		} else {
			players[0].led.High()
		}
		onoff = !onoff
		time.Sleep(100 * time.Millisecond)
	}
	players[0].led.Low()
	time.Sleep(2 * time.Second)
}
