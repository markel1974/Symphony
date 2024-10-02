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

// conv5to4 converts GCR encoded bytes into an array of 4 decoded bytes
func conv5to4(source []uint8) [4]uint8 {
	tDest := uint32(source[0])
	tDest <<= 13
	var dest [4]uint8
	destIdx := 0
	sourceIdx := 1
	for _, i := range []uint8{5, 7, 9, 11} {
		tDest |= (uint32(source[sourceIdx])) << i
		dest[destIdx] = _gcrFromTable[(tDest>>16)&0x1f] << 4
		tDest <<= 5
		dest[destIdx] |= _gcrFromTable[(tDest>>16)&0x1f]
		tDest <<= 5
		sourceIdx++
		destIdx++
	}
	return dest
}

// conv4to5 converts 4 bytes to 5 GCR encoded bytes
func conv4to5(from [4]uint8) []uint8 {
	encode := func(f uint8) uint16 {
		f1 := f >> 4
		f2 := f & 0xf
		return (_gcrTable[f1] << 5) | _gcrTable[f2]
	}
	to := make([]uint8, 5)
	g := encode(from[0])
	to[0] = uint8(g >> 2)
	to[1] = uint8((g << 6) & 0xc0)
	g = encode(from[1])
	to[1] |= uint8((g >> 4) & 0x3f)
	to[2] = uint8((g << 4) & 0xf0)
	g = encode(from[2])
	to[2] |= uint8((g >> 6) & 0x0f)
	to[3] = uint8((g << 2) & 0xfc)
	g = encode(from[3])
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
	usable           bool
}

func NewGCR(image []uint8) (*GCR, error) {
	g := &GCR{
		data:             make([]uint8, diskSize),
		errorInfo:        make([]uint8, numSectors),
		idx:              0,
		trackStart:       0,
		trackEnd:         trackSize,
		currentHalfTrack: 2,
		usable:           false,
	}
	for x := range g.data {
		g.data[x] = 0x55
	}
	for x := range g.errorInfo {
		g.errorInfo[x] = 1
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
		copy(g.errorInfo, image[numSectors*blockSize:])
	}
	bamSector := g.readSector(image, headerLen, bamTrackIdx, bamSectorIdx)
	if bamSector == nil {
		return nil, fmt.Errorf("nil bam sector")
	}
	bam1 := bamSector[162]
	bam2 := bamSector[163]
	for track := 1; track <= numTracks; track++ {
		for sector := 0; sector < int(_numSectors[track]); sector++ {
			if block := g.readSector(image, headerLen, track, sector); block != nil {
				if err := g.sector2gcr(block, bam1, bam2, track, sector); err != nil {
					return nil, err
				}
			}
		}
	}
	g.usable = true
	return g, nil
}

func (g *GCR) Usable() bool {
	return g.usable
}

func (g *GCR) MoveOut() {
	if g.currentHalfTrack <= 2 {
		return
	}
	g.currentHalfTrack--
	g.updateTrack()
}

func (g *GCR) MoveIn() {
	if g.currentHalfTrack >= numTracksMax {
		return
	}
	g.currentHalfTrack++
	g.updateTrack()
}

func (g *GCR) Rotate() {
	g.idx++
	if g.idx >= g.trackEnd {
		g.idx = g.trackStart
	}
}

func (g *GCR) Read() uint8 {
	return g.data[g.idx]
}

func (g *GCR) Write(data uint8) {
	g.data[g.idx] = data
}

func (g *GCR) updateTrack() {
	halfTrack := g.currentHalfTrack >> 1
	g.idx = (halfTrack - 1) * trackSize
	g.trackStart = g.idx
	trackLength := int(_numSectors[halfTrack]) * sectorSize
	g.trackEnd = g.trackStart + trackLength
}

func (g *GCR) readSector(diskData []byte, headerLen int, track int, sector int) []uint8 {
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

func (g *GCR) sector2gcr(block []uint8, bam1 uint8, bam2 uint8, track int, sector int) error {
	if len(block) > blockSize {
		return fmt.Errorf("invalid block length")
	}
	idx := ((track - 1) * trackSize) + (sector * sectorSize)
	g.data[idx] = 0xff
	idx++
	// Header mark [0], Checksum [1,2,3]
	p1 := [4]uint8{0x08, uint8(sector ^ track ^ int(bam2) ^ int(bam1)), uint8(sector), uint8(track)}
	copy(g.data[idx:], conv4to5(p1))

	p2 := [4]uint8{bam2, bam1, 0x0f, 0x0f}
	copy(g.data[idx+5:], conv4to5(p2))
	idx += 10
	for x := 0; x < 9; x++ {
		g.data[idx+x] = 0x55
	}
	idx += 9
	// Create GCR data + SYNC
	g.data[idx] = 0xff
	idx++
	// Data mark
	p3 := [4]uint8{0x07, block[0], block[1], block[2]}
	checksum := p3[1]
	checksum ^= p3[2]
	checksum ^= p3[3]
	copy(g.data[idx:], conv4to5(p3))
	idx += 5
	for x := 3; x < 255; x += 4 {
		p4 := [4]uint8{block[x], block[x+1], block[x+2], block[x+3]}
		checksum ^= p4[0]
		checksum ^= p4[1]
		checksum ^= p4[2]
		checksum ^= p4[3]
		copy(g.data[idx:], conv4to5(p4))
		idx += 5
	}
	checksum ^= block[255]
	p5 := [4]uint8{block[255], checksum, 0, 0}
	copy(g.data[idx:], conv4to5(p5))
	idx += 5
	for x := 0; x < 8; x++ {
		g.data[idx+x] = 0x55
	}
	return nil
}

/*
///Users/tinmr305/Desktop/emu/vice-emu-code-r45201-trunk-vice/src/gcr.c
func (g *GCR) gcr2sector(block []uint8, track int, sector int) []uint8 {
	var shift, i, j int
	var gcr[5] uint8
	var b uint8;
	var offsetIdx uint8
	var end uint8 = raw.data + raw.size;

	shift = p & 7;
	offsetIdx = raw->data + (p >> 3);

	b = offset[0] << shift;
	for i = 0; i < num; i++, buf += 4 {
		// get 5 bytes of gcr data
		for j = 0; j < 5; j++){
			offset++;
			if offset >= end {
				offset = raw->data;
			}
			if shift {
				gcr[j] = b | ((offset[0] << shift) >> 8);
				b = offset[0] << shift;
			} else {
				gcr[j] = b;
				b = offset[0];
			}
		}
		conv4to5(gcr, buf);
	}
}
*/
