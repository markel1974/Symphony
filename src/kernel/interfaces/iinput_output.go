package interfaces

type IInputOutput interface {
	IOWrite(data []byte) (int, error)

	IORead(p []byte) (int, error)

	IOType(kind KeyType, key rune)
}
