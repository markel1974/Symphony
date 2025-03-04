package roms

import (
	"os"
)

// Loader is a type responsible for loading ROM data from files or built-in ROM options.
type Loader struct {
}

// NewLoader creates and returns a new instance of the Loader struct.
func NewLoader() *Loader {
	return &Loader{}
}

// Load attempts to load the ROM data from a file if a valid name is provided; otherwise, it returns embedded ROM data.
// If useJiffy is true, it returns the _jiffyRom; if not, it defaults to returning the _builtinRom.
func (r *Loader) Load(useJiffy bool, romName string) []byte {
	if len(romName) > 0 {
		dat, err := os.ReadFile(romName)
		if err == nil {
			return dat
		}
	}
	if useJiffy {
		return _jiffyRom
	}
	//r.patchDriveRom(_jiffyRom)
	return _builtinRom
}
