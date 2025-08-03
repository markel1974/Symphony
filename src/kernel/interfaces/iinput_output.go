package interfaces

type IInputOutput interface {
	IOWrite(data []byte) (int, error)

	//IOType(kind KeyType, key rune)

	//IORead(p []byte) (int, error)
}
