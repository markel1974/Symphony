package references

// IC64RomSocket is an interface that provides methods to load various ROM sections, including Kernal, Basic, and Char ROMs.
// LoadKernal loads the Kernal ROM bytes and returns the data as a slice of bytes.
// LoadBasic loads the Basic ROM bytes and returns the data as a slice of bytes.
// LoadChar loads the Character ROM bytes and returns the data as a slice of bytes.
type IC64RomSocket interface {
	LoadKernal() []byte

	LoadBasic() []byte

	LoadChar() []byte
}
