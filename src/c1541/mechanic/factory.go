package mechanic

import (
	"github.com/markel1974/c64emu/src/c1541/disk/empty"
	"github.com/markel1974/c64emu/src/c1541/disk/gcr"
)

//func GetNumSectors(d int) uint8 {
//	return _numSectors[d]
//}

type IFloppy interface {
	Read() uint8
	Write(uint8)
	MoveOut()
	MoveIn()
	Rotate()
	Usable() bool
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Create(image []byte) (IFloppy, error) {
	if image == nil {
		return empty.NewEmpty(), nil
	}
	g, err := gcr.NewGCR(image)
	return g, err
}

/*
func (j *Mechanics) WriteSector() {
	track := j.banks.Read(0x18)
	sector := j.banks.Read(0x19)
	start := uint16(j.banks.Read(0x30)) | (uint16(j.banks.Read(0x31)) << 8)
	if start <= 0x0700 {
		block := j.banks.ReadInterval(start, BlockSize)
		if j.writeTrackSector(int(track), int(sector), block) {
			j.Sector2GCR(int(track), int(sector))
		}
	}
}

func (j *Mechanics) FormatTrack() {
	track := j.banks.Read(0x51)
	// Get new ID
	bufNum := j.banks.Read(0x3d)
	j.id1 = j.banks.Read(0x12 + uint16(bufNum))
	j.id2 = j.banks.Read(0x13 + uint16(bufNum))

	// Create empty block
	buf := make([]uint8, BlockSize)
	buf[0] = 0x4b

	// Write block to all sectors on track
	for sector := 0; sector < int(_numSectors[track]); sector++ {
		j.writeTrackSector(int(track), sector, buf)
		j.Sector2GCR(int(track), sector)
	}

	// Clear error info (all sectors no error)
	if track == 35 {
		for x := range j.errorInfo {
			j.errorInfo[x] = 1
		}
		// Write error_info to disk?
	}
}
*/

/*
func (j *Mechanics) writeTrackSector(track int, sector int, buffer []uint8) bool {
	offset := j.offsetFromTrackSector(track, sector)
	// Convert track/sector to byte offset in file
	if offset < 0 {
		return false
	}
	copy(j.diskData[offset+j.headerLen:], buffer)
	_ = os.WriteFile("a", j.diskData, 0644)
	return true
}
*/
