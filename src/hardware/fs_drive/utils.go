package fs_drive

import "unicode"

// ParseFileName parses the provided file name and extracts relevant information including mode, type, and record length.
// name: the input file name to parse.
// convertCharset: a boolean flag to determine if the character set should be converted.
// Returns the cleaned filename as a string, the mode as an int, the file type as an int, and the record length as an int.
func ParseFileName(name string, convertCharset bool) (string, int, int, int) {
	//TODO IMPLEMENT
	mode := FMODE_READ
	kind := FTYPE_PRG
	//dest uint8* , dest_len int& , mode int& , kind int, rec_len int& ,
	return "", mode, kind, 0
}

func FillName(n string) []byte {
	name := make([]byte, 16)
	for idx := range name {
		name[idx] = ' '
		if idx < len(n) {
			name[idx] = n[idx]
		}
	}
	return name
}

func CleanFileName(name string) []uint8 {
	const maxName = 16
	var vName []rune
	for _, k := range name {
		if k >= 'a' && k <= 'z' {
			vName = append(vName, unicode.ToUpper(k))
		} else if k >= 'A' && k <= 'Z' {
			vName = append(vName, k)
		} else if k >= '0' && k <= '9' {
			vName = append(vName, k)
		} else if k == '.' {
			vName = append(vName, k)
		} else {
			vName = append(vName, ' ')
		}
	}
	name = string(vName)
	//if ext := path.Ext(name); len(ext) > 0 {
	//	name = name[:len(name)-len(ext)]

	//if len(name) > maxName {
	//	name = name[:maxName]
	//}
	//return []uint8(name)
	//}
	if len(name) > maxName {
		name = name[:maxName]
	}
	return []uint8(name)
}
