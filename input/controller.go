package input

type Controller interface {
	Get() bool
	Reset()
}
