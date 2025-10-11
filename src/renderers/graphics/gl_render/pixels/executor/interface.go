package executor

// BeginEnder represents an interface with methods to mark the beginning and ending of an operation or process.
type BeginEnder interface {
	Begin()
	End()
}
