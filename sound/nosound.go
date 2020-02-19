package sound

type NoSound struct {
}

func NewNoSound() *NoSound {
	return &NoSound{}
}

func (n *NoSound) PlayStartFX() {
}

func (b *NoSound) PlayFinishFX() {
}
