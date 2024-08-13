package banks

import "os"

type Loader struct {
}

func NewLoader() *Loader {
	return &Loader{}
}

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

func (r *Loader) patchDriveRom(rom []byte) {
	rom[0x2ae4] = 0xea // Don't check ROM checksum
	rom[0x2ae5] = 0xea
	rom[0x2ae8] = 0xea
	rom[0x2ae9] = 0xea

	rom[0x2c9b] = 0xf2 // DOS idle loop
	rom[0x2c9c] = 0x00

	rom[0x3594] = 0x20 // Write sector
	rom[0x3595] = 0xf2
	rom[0x3596] = 0xf5
	rom[0x3597] = 0xf2
	rom[0x3598] = 0x01

	rom[0x3b0c] = 0xf2 // Format track
	rom[0x3b0d] = 0x02
}
