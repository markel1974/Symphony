package banks

import "os"

type RomLoader struct {
}

func NewRomLoader() *RomLoader {
	return &RomLoader{}
}

func (r *RomLoader) Load(defaultRom []byte, romName string) []byte {
	dat, err := os.ReadFile(romName)
	if err == nil {
		return dat
	}
	return defaultRom
}
