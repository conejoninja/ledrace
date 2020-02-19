package sound

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/buzzer"
)

type Buzzer struct {
	pin machine.Pin
	bzr buzzer.Device
}

func NewBuzzer(pin machine.Pin) *Buzzer {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return &Buzzer{
		pin: pin,
		bzr: buzzer.New(pin),
	}
}

func (b *Buzzer) PlayStartFX() {
	b.bzr.Tone(buzzer.E4, buzzer.Quarter)
	time.Sleep(400 * time.Millisecond)
	b.bzr.Tone(buzzer.E4, buzzer.Quarter)
	time.Sleep(400 * time.Millisecond)
	b.bzr.Tone(buzzer.B4, buzzer.Quarter)
	time.Sleep(400 * time.Millisecond)
}

func (b *Buzzer) PlayFinishFX() {
	b.bzr.Tone(buzzer.C4, 0.25)
	time.Sleep(100 * time.Millisecond)
	b.bzr.Tone(buzzer.C4, 0.25)
	time.Sleep(100 * time.Millisecond)
	b.bzr.Tone(buzzer.C4, 0.25)
	time.Sleep(100 * time.Millisecond)
	b.bzr.Tone(buzzer.C4, 0.25)
	time.Sleep(100 * time.Millisecond)
	b.bzr.Tone(buzzer.G3, 0.5)
	time.Sleep(200 * time.Millisecond)
	b.bzr.Tone(buzzer.A3, 0.5)
	time.Sleep(200 * time.Millisecond)
	b.bzr.Tone(buzzer.C4, 0.25)
	time.Sleep(100 * time.Millisecond)
	b.bzr.Tone(buzzer.A3, 0.25)
	time.Sleep(100 * time.Millisecond)
	b.bzr.Tone(buzzer.C4, 0.5)
	time.Sleep(100 * time.Millisecond)
}
