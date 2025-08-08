package interfaces

import "time"

// ProcessDescription represents the details of a system process, including its name and process ID.
type ProcessDescription struct {
	name     string
	user     string
	pid      int
	line     string
	starTime time.Time
}

func NewProcessDescription(name string, user string, pid int, line string, time time.Time) *ProcessDescription {
	return &ProcessDescription{
		name:     name,
		user:     user,
		pid:      pid,
		line:     line,
		starTime: time,
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

// User returns the name of the user associated with the process.
func (p *ProcessDescription) User() string {
	return p.user
}

func (p *ProcessDescription) Time() time.Time {
	return p.starTime
}
