package led

type Led interface {
	On()
	Off()
	Toggle()
	Init() bool
	Close()
}
