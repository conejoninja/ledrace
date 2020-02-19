package track

type NoTrack struct {
}

func NewNoTrack() *NoTrack {
	return &NoTrack{}
}

func (n *NoTrack) Draw() {
}

func (n *NoTrack) DrawGravity(gravity []uint8) {
}

func (n *NoTrack) Idle() {
}

func (n *NoTrack) DrawFinish(winner uint8) {
}

func (n *NoTrack) DrawStart() {
}
