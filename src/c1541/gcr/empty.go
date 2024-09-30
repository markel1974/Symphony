package gcr

type Empty struct{}

func NewEmpty() *Empty {
	return &Empty{}
}

func (e *Empty) Read() uint8 {
	return 0
}

func (e *Empty) Write(_ uint8) {
}

func (e *Empty) MoveOut() {
}

func (e *Empty) MoveIn() {
}

func (e *Empty) Rotate() {
}
