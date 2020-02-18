package input

import "machine"

type Button struct {
	pin     machine.Pin
	led     machine.Pin
	pressed bool
}

func NewButton(pin machine.Pin, led machine.Pin) *Button {
	pin.Configure(machine.PinConfig{Mode: machine.PinInput})
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	return &Button{
		pin: pin,
		led: led,
	}
}

func (b *Button) Get() bool {
	pressed := b.pin.Get()
	if b.pressed && !pressed {
		b.pressed = false
		b.led.High()
		return true
	} else if !b.pressed && pressed {
		b.pressed = true
		b.led.Low()
	}
	return false
}

func (b *Button) SpeedDelta() float32 {
	pressed := b.pin.Get()
	if b.pressed && !pressed {
		b.pressed = false
		b.led.High()
		return 1
	} else if !b.pressed && pressed {
		b.pressed = true
		b.led.Low()
	}
	return 0
}

func (b *Button) Reset() {
	b.pressed = false
	b.led.Low()
}