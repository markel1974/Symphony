package mechanics

import "fmt"

type GCR struct {
	data      []uint8
	errorInfo []uint8
}

func NewGCR() *GCR {
	d := &GCR{
		data:      make([]uint8, GCR_DISK_SIZE),
		errorInfo: make([]uint8, NUM_SECTORS),
	}
	for x := range d.data {
		d.data[x] = 0x55
	}
	for x := range d.errorInfo {
		d.errorInfo[x] = 1
	}
	return d
}

type GCRFactory struct {
}

func NewGCRFactory() *GCRFactory {
	return &GCRFactory{}
}

func (gcr *GCRFactory) readSector(diskData []byte, headerLen int, track int, sector int) []uint8 {
	// Convert track/sector to byte offset in file
	offset := gcr.offsetFromTrackSector(track, sector)
	if offset < 0 {
		return nil
	}
	start := offset + headerLen
	if end := start + BLOCK_SIZE; end >= len(diskData) {
		return nil
	}
	buffer := make([]uint8, BLOCK_SIZE)
	copy(buffer, diskData[start:])
	return buffer
}

// secNumFromTs Convert track/sector to offset
func (gcr *GCRFactory) secNumFromTs(track int, sector int) int {
	return _sectorOffset[track] + sector
}

// offsetFromTrackSector Convert track/sector to offset
func (gcr *GCRFactory) offsetFromTrackSector(track int, sector int) int {
	if (track < 1) || (track > NUM_TRACKS) || (sector < 0) || (sector >= int(_numSectors[track])) {
		return -1
	}
	return (_sectorOffset[track] + sector) << 8
}

// Conv4 Convert 4 bytes to 5 GCR encoded bytes
func (gcr *GCRFactory) conv4(from []uint8, to []uint8) {
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

func (gcr *GCRFactory) convertSector(diskData []byte, gcrData []uint8, headerLen int, id1 uint8, id2 uint8, track int, sector int) {
	buf := make([]uint8, 4)
	p := (track-1)*GCR_TRACK_SIZE + sector*GCR_SECTOR_SIZE
	block := gcr.readSector(diskData, headerLen, track, sector)
	if block == nil {
		return
	}
	// Create GCRFactory header
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

// Create GCR disk from image disk
func (gcr *GCRFactory) Create(image []byte) (*GCR, error) {
	d := NewGCR()
	diskDataLen := len(image)
	if diskDataLen < NUM_SECTORS*BLOCK_SIZE {
		return nil, fmt.Errorf("invalid disk data length")
	}
	headerLen := 0
	if image[0] == 0x43 && image[1] == 0x15 && image[2] == 0x41 && image[3] == 0x64 {
		headerLen = 64
	}
	// Load sector error info
	if headerLen == 0 && diskDataLen == NUM_SECTORS*257 {
		copy(d.errorInfo, image[NUM_SECTORS*BLOCK_SIZE:])
	}
	// Read BAM
	bam := gcr.readSector(image, headerLen, 18, 0)
	if bam == nil {
		return nil, fmt.Errorf("nil bam")
	}
	id1 := bam[162]
	id2 := bam[163]
	// Create GCR encoded disk from image
	for track := 1; track <= NUM_TRACKS; track++ {
		for sector := 0; sector < int(_numSectors[track]); sector++ {
			gcr.convertSector(image, d.data, headerLen, id1, id2, track, sector)
		}
	}
	return d, nil
}
