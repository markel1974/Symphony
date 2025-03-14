package fsdrive

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
