package track

import (
	"image/color"
	"machine"
	"time"

	"github.com/conejoninja/ledrace"

	"tinygo.org/x/drivers/ws2812"
)

type WS2812 struct {
	pin  machine.Pin
	leds []color.RGBA
	ws   ws2812.Device
	aux  int
	info *ledrace.Info
}

func NewWS2812(pin machine.Pin, info ledrace.Info) *WS2812 {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return &WS2812{
		pin:  pin,
		leds: make([]color.RGBA, info.TrackLength),
		ws:   ws2812.New(pin),
	}
}

func (w *WS2812) Draw() {
	for i := uint16(0); i < w.info.TrackLength; i++ {
		w.leds[i] = color.RGBA{0, 0, 0, 0}
	}
	for p := uint8(0); p < w.info.NumPlayers; p++ {
		for l := uint8(0); l < w.info.Players[p].Lap; l++ {
			w.leds[(uint16(w.info.Players[p].Position)-uint16(l))%w.info.TrackLength] = w.info.Players[p].Color
		}
	}
	w.ws.WriteColors(w.leds)
}

func (w *WS2812) DrawGravity(gravity []uint8) {
	for i := uint16(0); i < w.info.TrackLength; i++ {
		grav := gravity[i]
		r := uint8(2)
		g := uint8(2)
		if grav < 127 {
			r = 2 * (127 - grav)
		}
		if grav > 127 {
			g = (grav - 127) * 2
		}
		w.leds[i] = color.RGBA{r, g, 0, 255}
	}
	w.ws.WriteColors(w.leds)
}

func (w *WS2812) Idle() {
	for i := uint16(0); i < w.info.TrackLength; i++ {
		w.leds[i] = getRainbowRGB(uint8((i*256)/w.info.TrackLength) + uint8(w.aux))
	}
	w.ws.WriteColors(w.leds)
	w.aux = (w.aux + 1) % 255

}

func (w *WS2812) DrawFinish(winner uint8) {
	for k := 0; k < 6; k++ {
		for i := uint16(0); i < w.info.TrackLength; i++ {
			w.leds[i] = w.info.Players[winner].Color
		}
		w.ws.WriteColors(w.leds)
		time.Sleep(300 * time.Millisecond)
		for i := uint16(0); i < w.info.TrackLength; i++ {
			w.leds[i] = color.RGBA{0, 0, 0, 0}
		}
		w.ws.WriteColors(w.leds)
		time.Sleep(300 * time.Millisecond)
	}
}

func (w *WS2812) DrawStart() {
	for i := uint16(0); i < w.info.TrackLength; i++ {
		w.leds[i] = color.RGBA{0, 0, 0, 0}
	}
	w.ws.WriteColors(w.leds)
	time.Sleep(400 * time.Millisecond)


	/*for i := 0; i < g.Info.TrackLength; i++ {
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
	time.Sleep(400 * time.Millisecond)*/

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
