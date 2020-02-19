package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/conejoninja/ledrace"
	config "github.com/conejoninja/ledrace/config/local"
	"github.com/conejoninja/ledrace/info"
	"github.com/conejoninja/ledrace/input"
	"github.com/conejoninja/ledrace/sound"
	"github.com/conejoninja/ledrace/track"
)

func main() {

	time.Sleep(4*time.Second)

	players := make([]info.Player, config.PLAYERS)
	players[0].Configure(input.NewButton(machine.A5, machine.A6), color.RGBA{255, 0, 0, 255})
	players[1].Configure(input.NewButton(machine.A3, machine.A4), color.RGBA{0, 255, 0, 255})
	players[2].Configure(input.NewButton(machine.A1, machine.A2), color.RGBA{255, 255, 0, 255})
	players[3].Configure(input.NewButton(machine.D12, machine.D11), color.RGBA{0, 0, 255, 255})

	status := info.Status{
		Players:     players,
		NumPlayers:  config.PLAYERS,
		TrackLength: config.TRACKLENGTH,
		Laps:        config.LAPS,
		Gravity:     makeTrackWithGravity(12, 90, 100, 110),
		//		Gravity:     makeTrackNoGravity,
	}

	tracker := track.NewWS2812(machine.D2, &status)
	sounder := sound.NewBuzzer(machine.A0)
	game := ledrace.New(status, tracker, sounder)
	game.Configure()

	for {
		game.Loop()
	}

}

// TODO move these functions somewhere, load gravity from config
func makeTrackNoGravity() []uint8 {
	gravity := make([]uint8, config.TRACKLENGTH)
	for i := 0; i < config.TRACKLENGTH; i++ {
		gravity[i] = 127
	}
	return gravity
}

func makeTrackWithGravity(height uint8, start uint32, middle uint32, end uint32) []uint8 {
	gravity := make([]uint8, config.TRACKLENGTH)
	for i := 0; i < config.TRACKLENGTH; i++ {
		gravity[i] = 127
	}
	for i := uint32(0); i < (middle - start); i++ {
		gravity[start+i] = uint8(127 - float32(i)*float32(height)/float32(middle-start))
	}
	gravity[middle] = 127
	for i := uint32(0); i < (end - middle); i++ {
		gravity[middle+i+1] = uint8(127 + float32(height)*(1-float32(i)/float32(middle-start)))
	}
	return gravity
}
