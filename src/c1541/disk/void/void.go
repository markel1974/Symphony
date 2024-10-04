package void

type Void struct{}

func NewDisk() *Void {
	return &Void{}
}

func (e *Void) Read() uint8 {
	return 0
}

func (e *Void) Write(_ uint8) {
}

func (e *Void) MoveOut() {
}

func (e *Void) MoveIn() {
}

func (e *Void) Rotate() {
}

func (e *Void) Usable() bool {
	return false
}
