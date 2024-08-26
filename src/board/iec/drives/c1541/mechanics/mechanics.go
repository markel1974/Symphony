package mechanics

import (
	"io"
	"os"
)

//see https://sta.c64.org/cbm1541mem.html

type Mechanics struct {
	*Core
	errorInfo     []uint8 // Sector error information (1 byte/sector)
	gcrData       []uint8 // Pointer to GCRFactory encoded disk
	gcrIdx        int     // Pointer to GCRFactory disk under R/W head
	gcrTrackStart int     // Pointer to start of GCRFactory disk of current track
	gcrTrackEnd   int     // Pointer to end of GCRFactory disk of current track
	filePath      string
	deviceNumber  uint8
	banks         IBanks
	motor         bool
	factory       *GCRFactory
}

func NewMechanics(banks IBanks, deviceNumber uint8) *Mechanics {
	j := &Mechanics{
		Core:          NewCore(),
		deviceNumber:  deviceNumber,
		banks:         banks,
		motor:         false,
		gcrData:       nil,
		gcrIdx:        0,
		gcrTrackStart: 0,
		gcrTrackEnd:   0,
		errorInfo:     nil,
		factory:       NewGCRFactory(),
	}
	return j
}

func (j *Mechanics) Reset() {
	j.gcrIdx = 0
	j.gcrTrackStart = 0
	j.gcrTrackEnd = j.gcrTrackStart + GCR_TRACK_SIZE
	j.currentHalfTrack = 2
	j.writeProtected = false
	j.gcrData = nil
	j.errorInfo = nil
}

func (j *Mechanics) Setup(filePath string) {
	//filePath := prefs.GetDrivePath(int(j.deviceNumber - 8))
	if !j.HasDisk() {
		j.filePath = filePath
		_ = j.openFile(j.filePath)
	} else if j.filePath != filePath {
		j.filePath = filePath
		j.closeFile()
		_ = j.openFile(j.filePath)
		j.diskChanged = true
	}
}

func (j *Mechanics) RotateDisk() {
	j.gcrIdx++
	if j.gcrIdx == j.gcrTrackEnd {
		j.gcrIdx = j.gcrTrackStart
	}
}

func (j *Mechanics) SetMotor(m bool) {
	j.motor = m
}

func (j *Mechanics) HasDisk() bool {
	return j.gcrData != nil
}

func (j *Mechanics) WriteProtectionState() uint8 {
	r := uint8(0)
	if j.diskChanged {
		j.diskChanged = false
		if j.writeProtected {
			r = 0x10
		}
	} else {
		if !j.writeProtected {
			r = 0x10
		}
	}
	return r
}

func (j *Mechanics) SyncFound() bool {
	if j.gcrData == nil {
		return false
	}
	if j.gcrData[j.gcrIdx] == 0xff {
		return true
	}
	return false
}

func (j *Mechanics) ReadByte() uint8 {
	if j.gcrData == nil {
		return 0
	}
	data := j.gcrData[j.gcrIdx]
	return data
}

func (j *Mechanics) WriteByte(data uint8) {
	if j.gcrData == nil {
		return
	}
	j.gcrData[j.gcrIdx] = data
	//fmt.Println("WRITING ", j.gcrIdx, data)

	//track := j.gcrTrackStart
	//sector := _numSectors[j.currentHalfTrack>>1]
	//offset := j.offsetFromTrackSector(track, int(sector))
	//fmt.Println("WRITING ", j.gcrIdx, data)
	//fmt.Println(j.currentHalfTrack, sector, offset)
	//fmt.Println("------------------")
}

func (j *Mechanics) MoveHeadOut() {
	if j.currentHalfTrack == 2 {
		return
	}
	j.currentHalfTrack--
	idx := ((j.currentHalfTrack >> 1) - 1) * GCR_TRACK_SIZE
	j.gcrTrackStart = idx
	j.gcrIdx = idx
	trackLength := int(_numSectors[j.currentHalfTrack>>1]) * GCR_SECTOR_SIZE
	j.gcrTrackEnd = j.gcrTrackStart + trackLength
}

func (j *Mechanics) MoveHeadIn() {
	const maxTracks = NUM_TRACKS * 2
	if j.currentHalfTrack == maxTracks {
		return
	}
	j.currentHalfTrack++
	idx := ((j.currentHalfTrack >> 1) - 1) * GCR_TRACK_SIZE
	j.gcrTrackStart = idx
	j.gcrIdx = idx
	trackLength := int(_numSectors[j.currentHalfTrack>>1]) * GCR_SECTOR_SIZE
	j.gcrTrackEnd = j.gcrTrackStart + trackLength
}

func (j *Mechanics) openFile(filePath string) error {
	j.Reset()
	fd, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		j.writeProtected = true
		if fd, err = os.OpenFile(filePath, os.O_RDONLY, 0); err != nil {
			return err
		}
	}
	defer fd.Close()
	image, err := io.ReadAll(fd)
	if err != nil {
		return err
	}
	d, err := j.factory.Create(image)
	if err != nil {
		return err
	}
	j.gcrData = d.data
	j.errorInfo = d.errorInfo
	return nil
}

func (j *Mechanics) closeFile() {
	j.Reset()
	//TODO
}

//var _counter = 0
//var _start = int64(0)

//func (j *Mechanics) Emulate() {
//The measurement of the speed of rotation of the Commodore 1541 floppy disk drive is very accurate,
//as long as it is between 265 and 325 revolutions per minute. Outside these ranges,
//the program may return unreliable values, but this does not matter, since the speed must be set at 300 rpm.

//a rotation every 200ms (300rpm)
/*
	if j.motor {
		_counter++
		if _counter > 250000 {
			now := time.Now().UnixMilli()
			fmt.Println(now - _start)
			_start = now
			_counter = 0
			j.rotateDisk()
		}
	}
*/
//}

/*
func (j *Mechanics) WriteSector() {
	track := j.banks.Read(0x18)
	sector := j.banks.Read(0x19)
	start := uint16(j.banks.Read(0x30)) | (uint16(j.banks.Read(0x31)) << 8)
	if start <= 0x0700 {
		block := j.banks.ReadInterval(start, BLOCK_SIZE)
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
	buf := make([]uint8, BLOCK_SIZE)
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
