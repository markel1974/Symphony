package mechanics

import (
	"io"
	"os"
)

//see https://sta.c64.org/cbm1541mem.html

type Mechanics struct {
	*Core
	imageHeader   int     // Length of .d64/.x64 file header
	id1           uint8   // ID of disk
	id2           uint8   // ID of disk
	errorInfo     []uint8 // Sector error information (1 byte/sector)
	gcrData       []uint8 // Pointer to GCR encoded disk
	gcrIdx        int     // Pointer to GCR disk under R/W head
	gcrTrackStart int     // Pointer to start of GCR disk of current track
	gcrTrackEnd   int     // Pointer to end of GCR disk of current track
	diskData      []uint8
	filePath      string
	deviceNumber  uint8
	banks         IBanks
	motor         bool
}

func NewMechanics(banks IBanks, deviceNumber uint8) *Mechanics {
	j := &Mechanics{
		Core:         NewCore(),
		deviceNumber: deviceNumber,
		banks:        banks,
		motor:        false,
	}
	j.id1 = 0
	j.id2 = 0
	j.imageHeader = 0
	j.gcrData = make([]uint8, GCR_DISK_SIZE)
	j.gcrIdx = 0
	j.gcrTrackStart = 0
	j.gcrTrackEnd = j.gcrTrackStart + GCR_TRACK_SIZE
	j.currentHalfTrack = 2
	j.errorInfo = make([]uint8, NUM_SECTORS)
	return j
}

func (j *Mechanics) Setup(filePath string) {
	//filePath := prefs.GetDrivePath(int(j.deviceNumber - 8))
	if !j.HasDisk() {
		j.filePath = filePath
		j.openFile()
	} else if j.filePath != filePath {
		j.filePath = filePath
		j.closeFile()
		j.openFile()
		j.diskChanged = true
	}
}

//var _v = int64(0)

func (j *Mechanics) Emulate() {
	//The measurement of the speed of rotation of the Commodore 1541 floppy disk drive is very accurate,
	//as long as it is between 265 and 325 revolutions per minute. Outside these ranges,
	//the program may return unreliable values, but this does not matter, since the speed must be set at 300 rpm.
	/*
		if j.motor {
			now := time.Now().UnixMilli()
			v := now / 200
			if v > _v {
				_v = v
				//fmt.Println("ROTATION", v)
				j.rotateDisk()
			}
		}
	*/
}

func (j *Mechanics) RotateDisk() {
	//fmt.Println("ROTATION")
	j.gcrIdx++
	if j.gcrIdx == j.gcrTrackEnd {
		j.gcrIdx = j.gcrTrackStart
	}
}

func (j *Mechanics) SetMotor(m bool) {
	j.motor = m
}

func (j *Mechanics) HasDisk() bool {
	return j.diskData != nil
}

func (j *Mechanics) WriteProtectionState() uint8 {
	r := uint8(0)
	if j.diskChanged {
		// Disk change -> WP sensor strobe
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
	if j.gcrData[j.gcrIdx] == 0xff {
		return true
	}
	return false
}

func (j *Mechanics) ReadGCRByte() uint8 {
	data := j.gcrData[j.gcrIdx]
	return data
}

func (j *Mechanics) WriteGCRByte(data uint8) {
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

func (j *Mechanics) openFile() bool {
	for x := range j.gcrData {
		j.gcrData[x] = 0x55
	}
	j.writeProtected = false
	fd, err := os.OpenFile(j.filePath, os.O_RDWR, 0)
	if err != nil {
		j.writeProtected = true
		fd, err = os.OpenFile(j.filePath, os.O_RDONLY, 0)
	}
	if err != nil {
		return false
	}
	defer fd.Close()
	j.diskData, err = io.ReadAll(fd)
	if err != nil {
		return false
	}
	size := len(j.diskData)
	if size < NUM_SECTORS*BLOCK_SIZE {
		return false
	}
	if j.diskData[0] == 0x43 && j.diskData[1] == 0x15 && j.diskData[2] == 0x41 && j.diskData[3] == 0x64 {
		j.imageHeader = 64
	} else {
		j.imageHeader = 0
	}
	for x := range j.errorInfo {
		j.errorInfo[x] = 1
	}
	// Load sector error info from .d64 file, if present
	if j.imageHeader == 0 && size == NUM_SECTORS*257 {
		copy(j.errorInfo, j.diskData[NUM_SECTORS*BLOCK_SIZE:])
	}
	// Read BAM and get ID
	bam := j.readSector(18, 0)
	if bam == nil {
		return false
	}

	j.id1 = bam[162]
	j.id2 = bam[163]

	// Create GCR encoded disk from image
	j.disk2gcr()
	return true
}

func (j *Mechanics) closeFile() {
	j.diskData = nil
	//TODO
}

func (j *Mechanics) readSector(track int, sector int) []uint8 {
	// Convert track/sector to byte offset in file
	offset := j.offsetFromTrackSector(track, sector)
	if offset < 0 {
		return nil
	}
	start := offset + j.imageHeader
	if end := start + BLOCK_SIZE; end >= len(j.diskData) {
		return nil
	}
	buffer := make([]uint8, BLOCK_SIZE)
	copy(buffer, j.diskData[start:])
	return buffer
}

// secNumFromTs Convert track/sector to offset
func (j *Mechanics) secNumFromTs(track int, sector int) int {
	return _sectorOffset[track] + sector
}

// offsetFromTrackSector Convert track/sector to offset
func (j *Mechanics) offsetFromTrackSector(track int, sector int) int {
	if (track < 1) || (track > NUM_TRACKS) || (sector < 0) || (sector >= int(_numSectors[track])) {
		return -1
	}
	return (_sectorOffset[track] + sector) << 8
}

// gcrConv4 Convert 4 bytes to 5 GCR encoded bytes
func (j *Mechanics) gcrConv4(from []uint8, to []uint8) {
	g := (_gcrTable[from[0]>>4] << 5) | _gcrTable[from[0]&15]
	to[0] = uint8(g >> 2)
	to[1] = uint8((g << 6) & 0xc0)
	g = (_gcrTable[from[1]>>4] << 5) | _gcrTable[from[1]&15]
	to[1] |= uint8((g >> 4) & 0x3f)
	to[2] = uint8((g << 4) & 0xf0)
	g = (_gcrTable[from[2]>>4] << 5) | _gcrTable[from[2]&15]
	to[2] |= uint8((g >> 6) & 0x0f)
	to[3] = uint8((g << 2) & 0xfc)
	g = (_gcrTable[from[3]>>4] << 5) | _gcrTable[from[3]&15]
	to[3] |= uint8((g >> 8) & 0x03)
	to[4] = uint8(g)
}

func (j *Mechanics) sector2gcr(track int, sector int) {
	buf := make([]uint8, 4)
	p := (track-1)*GCR_TRACK_SIZE + sector*GCR_SECTOR_SIZE
	block := j.readSector(track, sector)
	if block == nil {
		return
	}
	// Create GCR header
	// SYNC
	j.gcrData[p] = 0xff
	p++
	// Header mark
	buf[0] = 0x08
	// Checksum
	buf[1] = uint8(sector ^ track ^ int(j.id2) ^ int(j.id1))
	buf[2] = uint8(sector)
	buf[3] = uint8(track)
	j.gcrConv4(buf, j.gcrData[p:])
	buf[0] = j.id2
	buf[1] = j.id1
	buf[2] = 0x0f
	buf[3] = 0x0f
	j.gcrConv4(buf, j.gcrData[p+5:])
	p += 10
	for x := 0; x < 9; x++ {
		j.gcrData[p+x] = 0x55
	}
	p += 9
	// Create GCR data
	// SYNC
	j.gcrData[p] = 0xff
	p++
	// Data mark
	buf[0] = 0x07
	buf[1] = block[0]
	sum := buf[1]
	buf[2] = block[1]
	sum ^= buf[2]
	buf[3] = block[2]
	sum ^= buf[3]
	j.gcrConv4(buf, j.gcrData[p:])
	p += 5
	for i := 3; i < 255; i += 4 {
		buf[0] = block[i]
		sum ^= buf[0]
		buf[1] = block[i+1]
		sum ^= buf[1]
		buf[2] = block[i+2]
		sum ^= buf[2]
		buf[3] = block[i+3]
		sum ^= buf[3]
		j.gcrConv4(buf, j.gcrData[p:])
		p += 5
	}
	buf[0] = block[255]
	sum ^= buf[0]
	// Checksum
	buf[1] = sum
	buf[2] = 0
	buf[3] = 0
	j.gcrConv4(buf, j.gcrData[p:])
	p += 5
	for x := 0; x < 8; x++ {
		j.gcrData[p+x] = 0x55
	}
}

// disk2gcr Convert all tracks and sectors
func (j *Mechanics) disk2gcr() {
	for track := 1; track <= NUM_TRACKS; track++ {
		for sector := 0; sector < int(_numSectors[track]); sector++ {
			j.sector2gcr(track, sector)
		}
	}
}

/*
func (j *Mechanics) WriteSector() {
	track := j.banks.Read(0x18)
	sector := j.banks.Read(0x19)
	start := uint16(j.banks.Read(0x30)) | (uint16(j.banks.Read(0x31)) << 8)
	if start <= 0x0700 {
		block := j.banks.ReadInterval(start, BLOCK_SIZE)
		if j.writeTrackSector(int(track), int(sector), block) {
			j.sector2gcr(int(track), int(sector))
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
		j.sector2gcr(int(track), sector)
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
	copy(j.diskData[offset+j.imageHeader:], buffer)
	_ = os.WriteFile("a", j.diskData, 0644)
	return true
}
*/
