package interfaces

// ProcessDescription represents the details of a system process, including its name and process ID.
type ProcessDescription struct {
	name string
	pid  int
	line string
}

func NewProcessDescription(name string, pid int, line string) *ProcessDescription {
	return &ProcessDescription{
		name: name,
		pid:  pid,
		line: line,
	}
}

// Name returns the name of the process
func (p *ProcessDescription) Name() string {
	return p.name
}

// PID returns the process ID
func (p *ProcessDescription) PID() int {
	return p.pid
}

// Line returns the command line string
func (p *ProcessDescription) Line() string {
	return p.line
}

/*
// HasPaint returns true if the ProcessDescription has a valid OnPaint assigned; otherwise, false.
func (p *ProcessDescription) HasPaint() bool {
	return p.process.GetCommand().OnPaint() != nil
}

// Paint executes the assigned OnPaint to render a task on the provided surface.
// TODO MOVE
func (p *ProcessDescription) Paint(surface ISurface) {
	p.process.Paint(surface)
}


*/
