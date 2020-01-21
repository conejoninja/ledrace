package ledrace

import (
	"encoding/hex"
	"image/color"
	"machine"
	"strconv"
	"time"

	"github.com/conejoninja/ledrace/telemetry"

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

type Info struct {
	TrackLength int `json:"track_length"`
	NumPlayers  int `json:"num_players"`
	Laps        int `json:"laps"`
}

type Player struct {
	button   machine.Pin
	led      machine.Pin
	pressed  bool
	ID       uint8      `json:"id"`
	Speed    float32    `json:speed"`
	Position float32    `json:position"`
	Color    color.RGBA `json:"Color"`
	Loop     uint8      `json:loop"`
}

var players [4]Player
var track Track
var bzr buzzer.Device
var info *telemetry.Telemetry
var idleTime time.Time
var activity bool

var black = color.RGBA{0, 0, 0, 0}
var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}
var yellow = color.RGBA{255, 255, 0, 255}
var blue = color.RGBA{0, 0, 255, 255}

func main() {
	bzrPin := machine.A0
	bzrPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	bzr = buzzer.New(bzrPin)

	players[0].button = machine.A5
	players[0].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[0].led = machine.A6
	players[0].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[0].Color = red

	players[1].button = machine.A3
	players[1].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[1].led = machine.A4
	players[1].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[1].Color = green

	players[2].button = machine.A1
	players[2].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[2].led = machine.A2
	players[2].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[2].Color = yellow

	players[3].button = machine.D12
	players[3].button.Configure(machine.PinConfig{Mode: machine.PinInput})
	players[3].led = machine.D11
	players[3].led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	players[3].Color = blue

	neo := machine.D2
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})
	track.ws = ws2812.New(neo)

	track.leds = make([]color.RGBA, TRACKLENGHT)
	track.gravity = make([]uint8, TRACKLENGHT)
	for i := 0; i < TRACKLENGHT; i++ {
		track.gravity[i] = 127
	}

	info = telemetry.New(&track.ws)

	setRamp(12, 90, 100, 110)

	// If at start, player 0 is pressed, enter configuration mode
	if !players[0].button.Get() {
		showGravity()
	}

	// until we get json.Marshall support
	playersStr := ""
	for p := 0; p < PLAYERS; p++ {
		playersStr += `{"id":`+strconv.Itoa(p)+`,"Color":"#` + hex.EncodeToString([]byte{players[p].Color.R, players[p].Color.G, players[p].Color.B, players[p].Color.A}) + `"}`
		if p < PLAYERS-1 {
			playersStr += ","
		}
	}
	println("going to send info mqtt")
	info.Send([]byte(`{"type":"info","track_length":` + strconv.Itoa(TRACKLENGHT) + `,"num_players":` + strconv.Itoa(PLAYERS) + `,"laps":` + strconv.Itoa(LAPS) + `,"players":[` + playersStr + `]}`))
	println("Start race")

	startRace()

	for {
		activity = false
		for p := uint8(0); p < PLAYERS; p++ {
			activity = activity || getPlayerInput(p)

			gravity := track.gravity[uint16(players[p].Position)%TRACKLENGHT]
			if gravity < 127 {
				players[p].Speed -= GRAVITY * float32(127-gravity)
			}
			if gravity > 127 {
				players[p].Speed += GRAVITY * float32(gravity-127)
			}

			players[p].Speed -= players[p].Speed * FRICTION
			players[p].Position += players[p].Speed

			if players[p].Position > TRACKLENGHT*float32(players[p].Loop) {
				players[p].Loop++
			}

		}

		paintTrack()

		// until we get json.Marshall support
		playersStr = ""
		for p := 0; p < PLAYERS; p++ {
			playersStr += `{"id":`+strconv.Itoa(p)+`,"speed":` + strconv.FormatFloat(float64(players[p].Speed), 'f', 2, 32) + `,"position":` + strconv.FormatFloat(float64(players[p].Position), 'f', 2, 32) + `,"loop":` + strconv.Itoa(int(players[p].Loop)) + `}`
			if p < PLAYERS-1 {
				playersStr += ","
			}
		}
		//go info.Send([]byte(`{"type":"status","players":[` + playersStr + `]}`))

		maxPosition := players[0].Position
		winner := uint8(0)
		for p := uint8(1); p < PLAYERS; p++ {
			if players[p].Position > maxPosition {
				maxPosition = players[p].Position
				winner = p
			}
		}
		if maxPosition > LAPS*TRACKLENGHT {
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
		players[p].Speed += ACCELERATION
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
		for l := uint8(0); l < players[p].Loop; l++ {
			track.leds[(uint32(players[p].Position)-uint32(l))%TRACKLENGHT] = players[p].Color
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
			track.leds[i] = players[winner].Color
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
		players[p].Speed = 0
		players[p].Position = 0
		players[p].pressed = false
		players[p].Loop = 0
		players[p].led.Low()
	}
}

func idleRace() {
	var k uint16
	activity = false
	for {
		for i := 0; i < TRACKLENGHT; i++ {
			track.leds[i] = getRainbowRGB(uint8((i*256)/TRACKLENGHT)+uint8(k))
		}
		track.ws.WriteColors(track.leds)
		k = (k+1)%255

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
