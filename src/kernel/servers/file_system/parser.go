package file_system

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

// argType represents an enumeration used for command-line argument parsing.
type argType int

// argNo represents an argument type with no associated value.
// argSingle represents an argument type with a single associated value.
// argQuoted represents an argument type with a quoted associated value.
const (
	argNo argType = iota
	argSingle
	argQuoted
)

// GetEnvFn defines a function type that takes a string argument and returns a string.
type GetEnvFn func(_ string) string

// Parser is a struct that facilitates parsing and tokenizing based on custom rules and configurations.
type Parser struct {
	parseEnv      bool
	parseBacktick bool
	position      int
	dir           string
	getEnv        func(string) string
	args          []string
	buf           string
	escaped       bool
	doubleQuoted  bool
	singleQuoted  bool
	backQuote     bool
	dollarQuote   bool
	backtick      string
	pos           int
	got           argType
}

// NewParser creates and initializes a new Parser instance with options for parsing environments and backticks.
func NewParser(parseEnv bool, parseBacktick bool, dir string) *Parser {
	return &Parser{
		parseEnv:      parseEnv,
		parseBacktick: parseBacktick,
		position:      0,
		dir:           dir,
		getEnv:        func(_ string) string { return "" },
	}
}

// Parse parses the given command line string into arguments, handling quotes, backticks, and escape sequences.
func (ps *Parser) Parse(line string) ([]string, error) {
	ps.args = []string{}
	ps.buf = ""
	ps.escaped = false
	ps.doubleQuoted = false
	ps.singleQuoted = false
	ps.backQuote = false
	ps.dollarQuote = false
	ps.backtick = ""
	ps.pos = -1
	ps.got = argNo

	i := -1
loop:
	for _, r := range line {
		i++
		if ps.escaped {
			if ps.doEscape(r) {
				continue
			}
		}

		if r == '\\' {
			if ps.doBackslash(r) {
				continue
			}
		}

		if unicode.IsSpace(r) {
			if err := ps.doSpace(r); err != nil {
				return nil, err
			}
			continue
		}

		switch r {
		case '`':
			if ps.doBacktick() {
				continue
			}
		case ')':
			if ps.doRightParenthesis() {
				continue
			}
		case '(':
			ok, err := ps.doLeftParenthesis()
			if err != nil {
				return nil, err
			}
			if ok {
				continue
			}
		case '"':
			if ps.doDoubleQuotation() {
				continue
			}
		case '\'':
			if ps.doSingleQuotation() {
				continue
			}
		case ';', '&', '|', '<', '>':
			if !(ps.escaped || ps.singleQuoted || ps.doubleQuoted || ps.backQuote || ps.dollarQuote) {
				if r == '>' && len(ps.buf) > 0 {
					if c := ps.buf[0]; '0' <= c && c <= '9' {
						i -= 1
						ps.got = argNo
					}
				}
				ps.pos = i
				break loop
			}
		}

		ps.got = argSingle
		ps.buf += string(r)
		if ps.backQuote || ps.dollarQuote {
			ps.backtick += string(r)
		}
	}

	if ps.got != argNo {
		if ps.parseEnv {
			if ps.got == argSingle {
				parser := NewParser(false, false, ps.dir)
				elements, err := parser.Parse(ps.replaceEnv(ps.buf))
				if err != nil {
					return nil, err
				}
				ps.args = append(ps.args, elements...)
			} else {
				ps.args = append(ps.args, ps.replaceEnv(ps.buf))
			}
		} else {
			ps.args = append(ps.args, ps.buf)
		}
	}
	if ps.escaped || ps.singleQuoted || ps.doubleQuoted || ps.backQuote || ps.dollarQuote {
		return nil, errors.New("invalid command line string")
	}
	ps.position = ps.pos
	return ps.args, nil
}

// ParseWithEnvs parses the input line, separating environment variables and arguments into two distinct slices.
func (ps *Parser) ParseWithEnvs(line string, env GetEnvFn) ([]string, []string, error) {
	ps.getEnv = env
	_args, err := ps.Parse(line)
	if err != nil {
		return nil, nil, err
	}
	var envs []string
	var args []string
	parsingEnv := true
	for _, arg := range _args {
		if parsingEnv && ps.isEnv(arg) {
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

// doEscape processes escape sequences for specific runes, updates the buffer, and resets the escaped state.
func (ps *Parser) doEscape(r rune) bool {
	if r == 't' {
		r = '\t'
	}
	if r == 'n' {
		r = '\n'
	}
	ps.buf += string(r)
	ps.escaped = false
	ps.got = argSingle
	return true
}

// doBackslash handles backslash characters during parsing, supporting single-quoted context or marking as escaped.
func (ps *Parser) doBackslash(r rune) bool {
	if ps.singleQuoted {
		ps.buf += string(r)
	} else {
		ps.escaped = true
	}
	return true
}

// doSpace processes a space character and updates the parser state based on the current quoting and parsing context.
func (ps *Parser) doSpace(r rune) error {
	if ps.singleQuoted || ps.doubleQuoted || ps.backQuote || ps.dollarQuote {
		ps.buf += string(r)
		ps.backtick += string(r)
	} else if ps.got != argNo {
		if ps.parseEnv {
			if ps.got == argSingle {
				parser := NewParser(false, false, ps.dir)
				elements, err := parser.Parse(ps.replaceEnv(ps.buf))
				if err != nil {
					return err
				}
				ps.args = append(ps.args, elements...)
			} else {
				ps.args = append(ps.args, ps.replaceEnv(ps.buf))
			}
		} else {
			ps.args = append(ps.args, ps.buf)
		}
		ps.buf = ""
		ps.got = argNo
	}
	return nil
}

// doBacktick toggles the state of backtick quotation in the parser and updates buffer if necessary.
func (ps *Parser) doBacktick() bool {
	if !ps.singleQuoted && !ps.doubleQuoted && !ps.dollarQuote {
		if ps.parseBacktick {
			if ps.backQuote {
				ps.buf = ps.buf[:len(ps.buf)-len(ps.backtick)]
			}
			ps.backtick = ""
			ps.backQuote = !ps.backQuote
			return true
		}
		ps.backtick = ""
		ps.backQuote = !ps.backQuote
	}
	return false
}

// doRightParenthesis processes a right parenthesis in the parser, toggling dollarQuote or clearing backtick as necessary.
func (ps *Parser) doRightParenthesis() bool {
	if !ps.singleQuoted && !ps.doubleQuoted && !ps.backQuote {
		if ps.parseBacktick {
			if ps.dollarQuote {
				ps.buf = ps.buf[:len(ps.buf)-len(ps.backtick)-2]
			}
			ps.backtick = ""
			ps.dollarQuote = !ps.dollarQuote
			return true
		}
		ps.backtick = ""
		ps.dollarQuote = !ps.dollarQuote
	}
	return false
}

// doLeftParenthesis handles the processing of a left parenthesis '(' in the parser logic, updating state as needed.
// Returns true if the character is successfully processed, or an error if parsing is invalid.
func (ps *Parser) doLeftParenthesis() (bool, error) {
	if !ps.singleQuoted && !ps.doubleQuoted && !ps.backQuote {
		if !ps.dollarQuote && strings.HasSuffix(ps.buf, "$") {
			ps.dollarQuote = true
			ps.buf += "("
			return true, nil
		} else {
			return false, errors.New("invalid command line string")
		}
	}
	return false, nil
}

// doDoubleQuotation toggles the doubleQuoted state if singleQuoted and dollarQuote are not active and updates the got state.
func (ps *Parser) doDoubleQuotation() bool {
	if !ps.singleQuoted && !ps.dollarQuote {
		if ps.doubleQuoted {
			ps.got = argQuoted
		}
		ps.doubleQuoted = !ps.doubleQuoted
		return true
	}
	return false
}

// doSingleQuotation toggles the state of single quotation in the parser. Returns true if the state was successfully toggled.
func (ps *Parser) doSingleQuotation() bool {
	if !ps.doubleQuoted && !ps.dollarQuote {
		if ps.singleQuoted {
			ps.got = argQuoted
		}
		ps.singleQuoted = !ps.singleQuoted
		return true
	}
	return false
}

// doOthers processes characters that do not fall into specific parsing cases and updates the parser's position or state.
func (ps *Parser) doOthers(r rune, i int) bool {
	if !(ps.escaped || ps.singleQuoted || ps.doubleQuoted || ps.backQuote || ps.dollarQuote) {
		if r == '>' && len(ps.buf) > 0 {
			if c := ps.buf[0]; '0' <= c && c <= '9' {
				i -= 1
				ps.got = argNo
			}
		}
		ps.pos = i
		return false
		//break loop
	}
	return true
}

// isEnv determines if the given string argument represents an environment variable in the format "key=value".
func (ps *Parser) isEnv(arg string) bool {
	return len(strings.Split(arg, "=")) == 2
}

// replaceEnv replaces environment variable references in the input string `s` using the provided `GetEnvFn` function.
// Escaped characters and variables enclosed within `${}` or preceded by `$` are processed correctly.
// Undefined variables result in an empty string unless handled otherwise in the custom `GetEnvFn`.
func (ps *Parser) replaceEnv(s string) string {
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
					buf.WriteString(ps.getEnv(s[p:i]))
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
					buf.WriteString(ps.getEnv(s[p:i]))
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
