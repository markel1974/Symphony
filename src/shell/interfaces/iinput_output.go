package interfaces

type IInputOutput interface {
	Write(data []byte) (int, error)

	Read(p []byte) (int, error)

	Type(kind KeyType, key rune)

	//TODO REMOVE
	History(verb HistoryAction, idx int)
}
