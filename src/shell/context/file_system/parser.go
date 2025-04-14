package file_system

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

// argType defines a custom integer type used to represent argument constants or enumerations.
type argType int

// argNo represents an argument type with no value.
// argSingle represents an argument type with a single value.
// argQuoted represents an argument type with a quoted value.
const (
	argNo argType = iota
	argSingle
	argQuoted
)

type GetEnvFn func(_ string) string

// isEnv determines if the provided string argument follows the format of an environment variable (key=value).
func isEnv(arg string) bool {
	return len(strings.Split(arg, "=")) == 2
}

// replaceEnv replaces substrings in the input string `s` using the provided `env` function for variable substitution.
// If `env` is nil, a default implementation (`getEnv`) is used for substitution.
func replaceEnv(env GetEnvFn, s string) string {
	if env == nil {
		env = func(_ string) string {
			return ""
		}
	}

	var buf bytes.Buffer
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if r == '\\' {
			i++
			if i == len(rs) {
				break
			}
			buf.WriteRune(rs[i])
			continue
		} else if r == '$' {
			i++
			if i == len(rs) {
				buf.WriteRune(r)
				break
			}
			if rs[i] == 0x7b {
				i++
				p := i
				for ; i < len(rs); i++ {
					r = rs[i]
					if r == '\\' {
						i++
						if i == len(rs) {
							return s
						}
						continue
					}
					if r == 0x7d || (!unicode.IsLetter(r) && r != '_' && !unicode.IsDigit(r)) {
						break
					}
				}
				if r != 0x7d {
					return s
				}
				if i > p {
					buf.WriteString(env(s[p:i]))
				}
			} else {
				p := i
				for ; i < len(rs); i++ {
					r := rs[i]
					if r == '\\' {
						i++
						if i == len(rs) {
							return s
						}
						continue
					}
					if !unicode.IsLetter(r) && r != '_' && !unicode.IsDigit(r) {
						break
					}
				}
				if i > p {
					buf.WriteString(env(s[p:i]))
					i--
				} else {
					buf.WriteString(s[p:])
				}
			}
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// Parser represents a command-line arguments and environment variables parser with customizable behavior.
type Parser struct {
	parseEnv      bool
	parseBacktick bool
	position      int
	dir           string
	getEnv        func(string) string
}

// NewParser initializes and returns a new instance of the Parser struct with default settings.
func NewParser(parseEnv bool, parseBacktick bool, dir string) *Parser {
	return &Parser{
		parseEnv:      parseEnv,
		parseBacktick: parseBacktick,
		position:      0,
		dir:           dir,
	}
}

// Parse processes a command-line string, splitting it into arguments and handling quotes, backticks, and environment variables.
func (p *Parser) Parse(line string) ([]string, error) {
	var args []string
	buf := ""
	escaped := false
	doubleQuoted := false
	singleQuoted := false
	backQuote := false
	dollarQuote := false
	backtick := ""
	pos := -1
	got := argNo

	i := -1
loop:
	for _, r := range line {
		i++
		if escaped {
			if r == 't' {
				r = '\t'
			}
			if r == 'n' {
				r = '\n'
			}
			buf += string(r)
			escaped = false
			got = argSingle
			continue
		}

		if r == '\\' {
			if singleQuoted {
				buf += string(r)
			} else {
				escaped = true
			}
			continue
		}

		if unicode.IsSpace(r) {
			if singleQuoted || doubleQuoted || backQuote || dollarQuote {
				buf += string(r)
				backtick += string(r)
			} else if got != argNo {
				if p.parseEnv {
					if got == argSingle {
						parser := NewParser(false, false, p.dir)
						elements, err := parser.Parse(replaceEnv(p.getEnv, buf))
						if err != nil {
							return nil, err
						}
						args = append(args, elements...)
					} else {
						args = append(args, replaceEnv(p.getEnv, buf))
					}
				} else {
					args = append(args, buf)
				}
				buf = ""
				got = argNo
			}
			continue
		}

		switch r {
		case '`':
			if !singleQuoted && !doubleQuoted && !dollarQuote {
				if p.parseBacktick {
					if backQuote {
						buf = buf[:len(buf)-len(backtick)]
					}
					backtick = ""
					backQuote = !backQuote
					continue
				}
				backtick = ""
				backQuote = !backQuote
			}
		case ')':
			if !singleQuoted && !doubleQuoted && !backQuote {
				if p.parseBacktick {
					if dollarQuote {
						buf = buf[:len(buf)-len(backtick)-2]
					}
					backtick = ""
					dollarQuote = !dollarQuote
					continue
				}
				backtick = ""
				dollarQuote = !dollarQuote
			}
		case '(':
			if !singleQuoted && !doubleQuoted && !backQuote {
				if !dollarQuote && strings.HasSuffix(buf, "$") {
					dollarQuote = true
					buf += "("
					continue
				} else {
					return nil, errors.New("invalid command line string")
				}
			}
		case '"':
			if !singleQuoted && !dollarQuote {
				if doubleQuoted {
					got = argQuoted
				}
				doubleQuoted = !doubleQuoted
				continue
			}
		case '\'':
			if !doubleQuoted && !dollarQuote {
				if singleQuoted {
					got = argQuoted
				}
				singleQuoted = !singleQuoted
				continue
			}
		case ';', '&', '|', '<', '>':
			if !(escaped || singleQuoted || doubleQuoted || backQuote || dollarQuote) {
				if r == '>' && len(buf) > 0 {
					if c := buf[0]; '0' <= c && c <= '9' {
						i -= 1
						got = argNo
					}
				}
				pos = i
				break loop
			}
		}

		got = argSingle
		buf += string(r)
		if backQuote || dollarQuote {
			backtick += string(r)
		}
	}

	if got != argNo {
		if p.parseEnv {
			if got == argSingle {
				parser := NewParser(false, false, p.dir)
				elements, err := parser.Parse(replaceEnv(p.getEnv, buf))
				if err != nil {
					return nil, err
				}
				args = append(args, elements...)
			} else {
				args = append(args, replaceEnv(p.getEnv, buf))
			}
		} else {
			args = append(args, buf)
		}
	}
	if escaped || singleQuoted || doubleQuoted || backQuote || dollarQuote {
		return nil, errors.New("invalid command line string")
	}
	p.position = pos
	return args, nil
}

// ParseWithEnvs splits a line into environment variables and arguments, returning them separately along with any parse error.
func (p *Parser) ParseWithEnvs(line string) ([]string, []string, error) {
	_args, err := p.Parse(line)
	if err != nil {
		return nil, nil, err
	}
	var envs []string
	var args []string
	parsingEnv := true
	for _, arg := range _args {
		if parsingEnv && isEnv(arg) {
			envs = append(envs, arg)
		} else {
			if parsingEnv {
				parsingEnv = false
			}
			args = append(args, arg)
		}
	}
	return envs, args, nil
}
