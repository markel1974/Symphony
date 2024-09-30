package gcr

import "fmt"

const (
	numTracks    = 35
	numSectors   = 683
	numTracksMax = numTracks * 2
	blockSize    = 256
	sectorSize   = 1 + 10 + 9 + 1 + 325 + 8 // SYNC Header Gap SYNC Data Gap (should be 5 SYNC bytes each)
	trackSize    = sectorSize * 21          // Each track in gcr has 21 sectors
	diskSize     = trackSize * numTracks
)

const (
	bamTrackIdx  = 18
	bamSectorIdx = 0
)

// Conv4 Convert 4 bytes to 5 GCR encoded bytes
func conv4to5(from [4]uint8) []uint8 {
	to := make([]uint8, 5)
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
	return to
}

type GCR struct {
	data             []uint8
	errorInfo        []uint8
	idx              int // Pointer to GCR disk under R/W head
	trackStart       int // Pointer to start of GCR disk of current track
	trackEnd         int // Pointer to end of GCR disk of current track
	currentHalfTrack int
}

func NewGCR(image []uint8) (*GCR, error) {
	gcr := &GCR{
		data:             make([]uint8, diskSize),
		errorInfo:        make([]uint8, numSectors),
		idx:              0,
		trackStart:       0,
		trackEnd:         trackSize,
		currentHalfTrack: 2,
	}
	for x := range gcr.data {
		gcr.data[x] = 0x55
	}
	for x := range gcr.errorInfo {
		gcr.errorInfo[x] = 1
	}
	imageLen := len(image)
	if imageLen < (numSectors * blockSize) {
		return nil, fmt.Errorf("invalid disk data length")
	}
	headerLen := 0
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		headerLen = 64
	}
	if (headerLen == 0) && (imageLen == (numSectors * 257)) {
		copy(gcr.errorInfo, image[numSectors*blockSize:])
	}
	bamSector := gcr.readSector(image, headerLen, bamTrackIdx, bamSectorIdx)
	if bamSector == nil {
		return nil, fmt.Errorf("nil bam sector")
	}
	bam1 := bamSector[162]
	bam2 := bamSector[163]
	for track := 1; track <= numTracks; track++ {
		for sector := 0; sector < int(_numSectors[track]); sector++ {
			if block := gcr.readSector(image, headerLen, track, sector); block != nil {
				if err := gcr.convertSector(block, bam1, bam2, track, sector); err != nil {
					return nil, err
				}
			}
		}
	}
	return gcr, nil
}

func (gcr *GCR) MoveOut() {
	if gcr.currentHalfTrack <= 2 {
		return
	}
	gcr.currentHalfTrack--
	gcr.updateTrack()
}

func (gcr *GCR) MoveIn() {
	if gcr.currentHalfTrack >= numTracksMax {
		return
	}
	gcr.currentHalfTrack++
	gcr.updateTrack()
}

func (gcr *GCR) Rotate() {
	gcr.idx++
	if gcr.idx >= gcr.trackEnd {
		gcr.idx = gcr.trackStart
	}
}

func (gcr *GCR) Read() uint8 {
	return gcr.data[gcr.idx]
}

func (gcr *GCR) Write(data uint8) {
	gcr.data[gcr.idx] = data
}

func (gcr *GCR) updateTrack() {
	halfTrack := gcr.currentHalfTrack >> 1
	gcr.idx = (halfTrack - 1) * trackSize
	gcr.trackStart = gcr.idx
	trackLength := int(_numSectors[halfTrack]) * sectorSize
	gcr.trackEnd = gcr.trackStart + trackLength
}

func (gcr *GCR) readSector(diskData []byte, headerLen int, track int, sector int) []uint8 {
	if (track < 1) || (track > numTracks) || (sector < 0) || (sector >= int(_numSectors[track])) {
		return nil
	}
	offset := (_sectorOffset[track] + sector) << 8
	if offset < 0 {
		return nil
	}
	start := offset + headerLen
	if end := start + blockSize; end >= len(diskData) {
		return nil
	}
	buffer := make([]uint8, blockSize)
	copy(buffer, diskData[start:])
	return buffer
}

func (gcr *GCR) convertSector(block []uint8, bam1 uint8, bam2 uint8, track int, sector int) error {
	if len(block) > blockSize {
		return fmt.Errorf("invalid block length")
	}
	idx := ((track - 1) * trackSize) + (sector * sectorSize)
	gcr.data[idx] = 0xff
	idx++
	// Header mark [0], Checksum [1,2,3]
	p1 := [4]uint8{0x08, uint8(sector ^ track ^ int(bam2) ^ int(bam1)), uint8(sector), uint8(track)}
	copy(gcr.data[idx:], conv4to5(p1))
	p2 := [4]uint8{bam2, bam1, 0x0f, 0x0f}
	copy(gcr.data[idx+5:], conv4to5(p2))
	idx += 10
	for x := 0; x < 9; x++ {
		gcr.data[idx+x] = 0x55
	}
	idx += 9
	// Create GCR data + SYNC
	gcr.data[idx] = 0xff
	idx++
	// Data mark
	p3 := [4]uint8{0x07, block[0], block[1], block[2]}
	checksum := p3[1]
	checksum ^= p3[2]
	checksum ^= p3[3]
	copy(gcr.data[idx:], conv4to5(p3))
	idx += 5
	for x := 3; x < 255; x += 4 {
		p4 := [4]uint8{block[x], block[x+1], block[x+2], block[x+3]}
		checksum ^= p4[0]
		checksum ^= p4[1]
		checksum ^= p4[2]
		checksum ^= p4[3]
		copy(gcr.data[idx:], conv4to5(p4))
		idx += 5
	}
	checksum ^= block[255]
	p5 := [4]uint8{block[255], checksum, 0, 0}
	copy(gcr.data[idx:], conv4to5(p5))
	idx += 5
	for x := 0; x < 8; x++ {
		gcr.data[idx+x] = 0x55
	}
	return nil
}
