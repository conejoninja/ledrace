package ledrace

import (
	"time"

	"github.com/conejoninja/ledrace/info"

	config "github.com/conejoninja/ledrace/config/local"
	"github.com/conejoninja/ledrace/sound"
	"github.com/conejoninja/ledrace/track"
)

type Game struct {
	Status   info.Status
	Track    track.Tracker
	Sound    sound.Sounder
	IdleTime time.Time
}

func New(status info.Status, tracker track.Tracker, sounder sound.Sounder) *Game {
	return &Game{
		Status: status,
		Track:  tracker,
		Sound:  sounder,
	}
}

func (g *Game) Configure() {
	if g.GetInput(0) {
		g.Track.DrawGravity(g.Status.Gravity)
		presses := 0
		for {
			if g.GetInput(0) {
				presses++
				if presses > 20 {
					break
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(2 * time.Second)
	}
}

func (g *Game) Loop() {
	g.StartRace()
	var gravity uint8
	var activity bool
	for {
		activity = false
		for p := uint8(0); p < g.Status.NumPlayers; p++ {
			if g.GetInput(p) {
				g.Status.Players[p].Speed += config.ACCELERATION
				activity = true
			}

			gravity = g.Status.Gravity[uint16(g.Status.Players[p].Position)%g.Status.TrackLength]
			if gravity < 127 {
				g.Status.Players[p].Speed -= config.GRAVITY * float32(127-gravity)
			}
			if gravity > 127 {
				g.Status.Players[p].Speed += config.GRAVITY * float32(gravity-127)
			}

				g.Status.Players[p].Speed -= g.Status.Players[p].Speed * config.FRICTION
			g.Status.Players[p].Position += g.Status.Players[p].Speed

			if g.Status.Players[p].Position > float32(g.Status.TrackLength)*float32(g.Status.Players[p].Lap) {
				g.Status.Players[p].Lap++
			}
		}

		g.Track.Draw()

		maxPosition := g.Status.Players[0].Position
		winner := uint8(0)
		for p := uint8(0); p < g.Status.NumPlayers; p++ {
			if g.Status.Players[p].Position > maxPosition {
				maxPosition = g.Status.Players[p].Position
				winner = p
			}
		}

		if maxPosition > float32(uint16(g.Status.Laps)*g.Status.TrackLength) {
			g.FinishRace(winner)
		}

		if activity {
			g.IdleTime = time.Now()
		} else if time.Since(g.IdleTime) > 30*time.Second {
			g.IdleRace()
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (g *Game) GetInput(p uint8) bool {
	return g.Status.Players[p].Input.Get()
}

func (g *Game) ResetPlayers() {
	for p := uint8(0); p < g.Status.NumPlayers; p++ {
		g.Status.Players[p].Speed = 0
		g.Status.Players[p].Position = 0
		g.Status.Players[p].Input.Reset()
		g.Status.Players[p].Lap = 0
	}
}

func (g *Game) FinishRace(winner uint8) {
	g.Track.DrawFinish(winner)
	g.Sound.PlayFinishFX()
	g.StartRace()
}

func (g *Game) IdleRace() {
	keepIdleing := true
	for keepIdleing {
		g.Track.Idle()
		for p := uint8(0); p < g.Status.NumPlayers; p++ {
			if g.Status.Players[p].Input.Get() {
				keepIdleing = false
				break
			}
		}
		time.Sleep(16 * time.Millisecond)
	}
	g.StartRace()
}

func (g *Game) StartRace() {
	g.ResetPlayers()

	g.IdleTime = time.Now()
	time.Sleep(2 * time.Second)
	g.Track.DrawStart()
	g.Sound.PlayStartFX()
}
