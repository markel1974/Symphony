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
