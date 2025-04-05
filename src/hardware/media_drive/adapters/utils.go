package adapters

import (
	"path"
	"strconv"
	"strings"
)

// _asciiToPetsciiTable is a lookup table for converting ASCII characters to PETSCII equivalents.
var _asciiToPetsciiTable [256]byte

// maxName defines the maximum allowable length for a name, typically used for filename constraints or similar purposes.
const maxName = 16

// init initializes the ASCII to PETSCII translation table with default mappings and specific character conversions.
func init() {
	const replacement byte = '?'
	for i := range _asciiToPetsciiTable {
		_asciiToPetsciiTable[i] = replacement
	}
	const cx = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_"
	for _, c := range []byte(cx) {
		_asciiToPetsciiTable[c] = c
	}
	for c := byte('a'); c <= 'z'; c++ {
		_asciiToPetsciiTable[c] = c - 32
	}
	_asciiToPetsciiTable['\n'] = 0x0D
}

func ParseFileName(nameWithParams string) (string, int, int, int) {
	mode := FModeRead
	filetype := FTypePrg
	recLen := 0
	parts := strings.SplitN(nameWithParams, ",", 2)
	filename := strings.TrimSpace(parts[0])
	if len(filename) > maxName {
		filename = filename[:maxName]
	}

	if len(parts) > 1 {
		params := strings.Split(parts[1], ",")
		parsingRecLen := false
		for _, p := range params {
			param := strings.TrimSpace(strings.ToUpper(p))
			if parsingRecLen {
				rl, err := strconv.Atoi(param)
				if err == nil && rl > 0 {
					recLen = rl
				} else {
					// Errore: parametro non numerico dopo L, o valore non valido TODO ERRORE
				}
				parsingRecLen = false
				continue
			}

			if len(param) == 1 {
				switch param[0] {
				case 'D':
					filetype = FTypeDel
				case 'S':
					filetype = FTypeSeq
				case 'P':
					filetype = FTypePrg
				case 'U':
					filetype = FTypeUsr
				case 'L':
					filetype = FTypeRel
					parsingRecLen = true
				case 'R':
					mode = FModeRead
				case 'W':
					mode = FModeWrite
				case 'A':
					mode = FModeAppend
				case 'M':
					mode = FModeM
				}
			}
		}
		if parsingRecLen {
			// Errore: L senza lunghezza record
		}
	}
	return filename, mode, filetype, recLen
}

// CreateFileNameFilled generates a byte slice with a file name padded to a fixed length, filling with the specified byte.
func CreateFileNameFilled(in string, fill uint8) []byte {
	v := CreateFileName(path.Base(in))
	name := make([]byte, maxName)
	for idx := range name {
		if idx < len(v) {
			name[idx] = v[idx]
		} else {
			name[idx] = fill
		}
	}
	return name
}

// CreateFileName converts a string to a PETSCII-encoded filename as a slice of bytes, truncated to a maximum length.
func CreateFileName(name string) []uint8 {
	var vName []uint8
	for _, k := range name {
		vName = append(vName, _asciiToPetsciiTable[k&0xff])
		if len(vName) >= maxName {
			break
		}
	}
	return vName
}
