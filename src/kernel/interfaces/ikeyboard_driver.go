package interfaces

type IKeyboardDriver interface {
	ScanKey(readBuffer []byte) (KeyType, rune, error)
}
