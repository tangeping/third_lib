package timingwheel

type DelayCall interface {
	String() string
	Call()
}
