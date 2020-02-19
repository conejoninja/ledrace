package info

import (
	"image/color"

	"github.com/conejoninja/ledrace/input"
)

type Player struct {
	Input    input.Controller `json:"input"`
	Speed    float32          `json:"speed"`
	Position float32          `json:"position"`
	Color    color.RGBA       `json:"color"`
	Lap      uint8            `json:"lap"`
}

type Status struct {
	Players     []Player `json:"players"`
	NumPlayers  uint8    `json:"num_players"`
	TrackLength uint16   `json:"track_length"`
	Laps        uint8    `json:"laps"`
	Gravity     []uint8  `json:"gravity"`
}

func (p *Player) Configure(input input.Controller, c color.RGBA) {
	p.Input = input
	p.Color = c
}
