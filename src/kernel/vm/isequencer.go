package vm

// ISequencer defines an interface to generate a sequence of functions for a given Virtual Machine instance.
type ISequencer interface {
	Create() []IOpExecutor
}
