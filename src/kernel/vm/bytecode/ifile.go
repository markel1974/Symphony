package bytecode

type IFile interface {
	Position(p int) (*FilePos, error)

	Size() int

	Base() int
}
