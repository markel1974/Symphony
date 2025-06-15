package sam

type Line struct {
	_in_ptr int
	command []rune
}

func (s *Line) Set(data string) {
	s.command = []rune(data)
}

func (s *Line) getChar() rune {
	s._in_ptr++
	if s._in_ptr >= len(s.command) {
		s._in_ptr = len(s.command) - 1
	}
	v := s.command[s._in_ptr]
	return v
}

func (s *Line) putBack(c rune) {
	if s._in_ptr > 0 {
		s._in_ptr--
	}
	s.command[s._in_ptr] = c
}
