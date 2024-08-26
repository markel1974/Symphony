package gcr

import "fmt"

//func GetNumSectors(d int) uint8 {
//	return _numSectors[d]
//}

func GetTrackLen(d int) int {
	trackLength := int(_numSectors[d]) * SectorSize
	return trackLength
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (gcr *Factory) readSector(diskData []byte, headerLen int, track int, sector int) []uint8 {
	offset := gcr.offsetFromTrackSector(track, sector)
	if offset < 0 {
		return nil
	}
	start := offset + headerLen
	if end := start + BlockSize; end >= len(diskData) {
		return nil
	}
	buffer := make([]uint8, BlockSize)
	copy(buffer, diskData[start:])
	return buffer
}

func (gcr *Factory) secNumFromTs(track int, sector int) int {
	return _sectorOffset[track] + sector
}

func (gcr *Factory) offsetFromTrackSector(track int, sector int) int {
	if (track < 1) || (track > NumTracks) || (sector < 0) || (sector >= int(_numSectors[track])) {
		return -1
	}
	return (_sectorOffset[track] + sector) << 8
}

// Conv4 Convert 4 bytes to 5 GCR encoded bytes
func (gcr *Factory) conv4(from []uint8, to []uint8) {
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

func (gcr *Factory) convertSector(diskData []byte, gcrData []uint8, headerLen int, id1 uint8, id2 uint8, track int, sector int) {
	buf := make([]uint8, 4)
	p := (track-1)*TrackSize + sector*SectorSize
	block := gcr.readSector(diskData, headerLen, track, sector)
	if block == nil {
		return
	}
	// Create GCR header
	// SYNC
	gcrData[p] = 0xff
	p++
	// Header mark
	buf[0] = 0x08
	// Checksum
	buf[1] = uint8(sector ^ track ^ int(id2) ^ int(id1))
	buf[2] = uint8(sector)
	buf[3] = uint8(track)
	gcr.conv4(buf, gcrData[p:])
	buf[0] = id2
	buf[1] = id1
	buf[2] = 0x0f
	buf[3] = 0x0f
	gcr.conv4(buf, gcrData[p+5:])
	p += 10
	for x := 0; x < 9; x++ {
		gcrData[p+x] = 0x55
	}
	p += 9
	// Create GCR data
	// SYNC
	gcrData[p] = 0xff
	p++
	// Data mark
	buf[0] = 0x07
	buf[1] = block[0]
	sum := buf[1]
	buf[2] = block[1]
	sum ^= buf[2]
	buf[3] = block[2]
	sum ^= buf[3]
	gcr.conv4(buf, gcrData[p:])
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
		gcr.conv4(buf, gcrData[p:])
		p += 5
	}
	buf[0] = block[255]
	sum ^= buf[0]
	// Checksum
	buf[1] = sum
	buf[2] = 0
	buf[3] = 0
	gcr.conv4(buf, gcrData[p:])
	p += 5
	for x := 0; x < 8; x++ {
		gcrData[p+x] = 0x55
	}
}

func (gcr *Factory) Create(image []byte) (*GCR, error) {
	d := NewGCR()
	diskDataLen := len(image)
	if diskDataLen < NumSectors*BlockSize {
		return nil, fmt.Errorf("invalid disk data length")
	}
	headerLen := 0
	if image[0] == 0x43 && image[1] == 0x15 && image[2] == 0x41 && image[3] == 0x64 {
		headerLen = 64
	}
	// Load sector error info
	if headerLen == 0 && diskDataLen == NumSectors*257 {
		copy(d.errorInfo, image[NumSectors*BlockSize:])
	}
	// Read BAM
	bam := gcr.readSector(image, headerLen, 18, 0)
	if bam == nil {
		return nil, fmt.Errorf("nil bam")
	}
	id1 := bam[162]
	id2 := bam[163]
	// Create GCR encoded disk from image
	for track := 1; track <= NumTracks; track++ {
		for sector := 0; sector < int(_numSectors[track]); sector++ {
			gcr.convertSector(image, d.data, headerLen, id1, id2, track, sector)
		}
	}
	return d, nil
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
