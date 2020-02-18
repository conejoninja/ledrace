package ledrace

import (
	"image/color"
	"time"

	config "github.com/conejoninja/ledrace/config/local"
	"github.com/conejoninja/ledrace/sound"
	"github.com/conejoninja/ledrace/track"

	"github.com/conejoninja/ledrace/input"
	"tinygo.org/x/drivers/buzzer"
)

type Player struct {
	Input    input.Inputer
	Speed    float32
	Position float32
	Color    color.RGBA
	Lap      uint8
}

type Info struct {
	Players     []Player
	NumPlayers  uint8
	TrackLength uint16
	Laps        uint8
	Gravity     []uint8
}

type Game struct {
	Info   Info
	Track  track.Tracker
	Sound  sound.Sounder
}

var players [4]Player
var bzr buzzer.Device
var idleTime time.Time
var activity bool

var black = color.RGBA{0, 0, 0, 0}
var red = color.RGBA{255, 0, 0, 255}
var green = color.RGBA{0, 255, 0, 255}
var yellow = color.RGBA{255, 255, 0, 255}
var blue = color.RGBA{0, 0, 255, 255}

func New(info Info, tracker track.Tracker, sounder sound.Sounder) *Game {
	return &Game{
		Info:   info,
		Track:  tracker,
		Sound:  sounder,
		Status: START,
	}
}

func (g *Game) Configure() {
	if g.GetInput(0) {
		g.Track.DrawGravity(g.Info.Gravity)
		longpress := 0
		for {
			if g.GetInput(0) {
				longpress++
				if longpress > 20 {
					break
				}
			} else {
				longpress = 0
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(2 * time.Second)
	}
}

func (g *Game) Loop() {
	g.StartRace()
	var gravity uint8
	for {
		activity = false
		for p := uint8(0); p < g.Info.NumPlayers; p++ {
			//activity = activity || getPlayerInput(p)

			gravity = g.Info.Gravity[uint16(g.Info.Players[p].Position)%g.Info.TrackLength]
			if gravity < 127 {
				g.Info.Players[p].Speed -= config.GRAVITY * float32(127-gravity)
			}
			if gravity > 127 {
				g.Info.Players[p].Speed += config.GRAVITY * float32(gravity-127)
			}

			g.Info.Players[p].Speed -= g.Info.Players[p].Speed * config.FRICTION
			g.Info.Players[p].Position += g.Info.Players[p].Speed

			if g.Info.Players[p].Position > float32(g.Info.TrackLength)*float32(g.Info.Players[p].Lap) {
				g.Info.Players[p].Speed++
			}

		}

		g.Track.Draw()

		maxPosition := g.Info.Players[0].Position
		winner := uint8(0)
		for p := uint8(0); p < g.Info.NumPlayers; p++ {
			if g.Info.Players[p].Position > maxPosition {
				maxPosition = g.Info.Players[p].Position
				winner = p
			}
		}

		if maxPosition > float32(uint16(g.Info.Laps)*g.Info.TrackLength) {
			time.Sleep(1500 * time.Millisecond)
			g.FinishRace(winner)
		}

		if activity {
			idleTime = time.Now()
		} else if time.Since(idleTime) > 30*time.Second {
			g.IdleRace()
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (g *Game) GetInput(p int) bool {
	return !g.Info.Players[0].Input.Get()
}

func (g *Game) ResetPlayers() {
	for p := uint8(0); p < g.Info.NumPlayers; p++ {
		g.Info.Players[p].Speed = 0
		g.Info.Players[p].Position = 0
		g.Info.Players[p].Input.Reset()
		g.Info.Players[p].Lap = 0
	}
}

func (g *Game) FinishRace(winner uint8) {
	g.Sound.PlayFinishFX()
	g.Track.DrawFinish(winner)
	g.StartRace()
}

func (g *Game) IdleRace() {
	for {
		g.Track.Idle()
		/*for p := uint8(0); p < PLAYERS; p++ {
			if getPlayerInput(p) {
				activity = true
				break
			}
		}
		if activity {
			break
		} */

		time.Sleep(16 * time.Millisecond)
	}
	g.StartRace()
}

func (g *Game) StartRace() {
	g.ResetPlayers()

	idleTime = time.Now()
	time.Sleep(2 * time.Second)
	g.Track.DrawStart()
}

func (p *Player) Configure(input input.Inputer, c color.RGBA) {
	p.Input = input
	p.Color = c
}
