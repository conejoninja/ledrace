package track

type Tracker interface {
	Draw()
	DrawStart()
	DrawFinish(winner uint8)
	DrawGravity(gravity []uint8)
	Idle()
}
