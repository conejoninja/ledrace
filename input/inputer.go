package input

type Inputer interface {
	Get() bool
	SpeedDelta() float32
	Reset()
}
