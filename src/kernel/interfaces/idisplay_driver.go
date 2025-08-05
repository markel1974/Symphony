package interfaces

// IDisplayDriver defines methods for handling terminal video display functionalities including text color and cursor control.
type IDisplayDriver interface {
	ITerminal
	Write(p []byte) (n int, err error)
}
