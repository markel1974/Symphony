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
	headerLen := uint8(0)
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		headerLen = 64
	}
	if (headerLen == 0) && (imageLen == (numSectors * 257)) {
		copy(g.errorInfo, image[numSectors*blockSize:])
	}
	bamSector, _ := g.readSector(image, headerLen, bamTrackIdx, bamSectorIdx)
	if bamSector == nil {
		return nil, fmt.Errorf("invalid track/sector")
	}
	bam1 := bamSector[162]
	bam2 := bamSector[163]
	for trackIdx := uint8(1); trackIdx <= numTracks; trackIdx++ {
		for sectorIdx := uint8(0); sectorIdx < _numSectors[trackIdx]; sectorIdx++ {
			sector, _ := g.readSector(image, headerLen, trackIdx, sectorIdx)
			if sector != nil {
				gcrSector, err := g.sector2gcr(sector, bam1, bam2, trackIdx, sectorIdx)
				if err != nil {
					return nil, err
				}
				idx := (int(trackIdx-1) * trackSize) + (int(sectorIdx) * sectorSize)
				copy(g.data[idx:], gcrSector)
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
	//if g.data[g.idx] != 0xff {
	//	fmt.Printf("Read %d -> %02x\n", g.idx, g.data[g.idx])
	//}
	return g.data[g.idx]
}

func (g *GCR) Write(data uint8) {
	//fmt.Printf("Write %d -> %02x\n", g.idx, data)
	g.data[g.idx] = data
}

func (g *GCR) updateTrack() {
	halfTrack := g.currentHalfTrack >> 1
	g.idx = (halfTrack - 1) * trackSize
	g.trackStart = g.idx
	trackLength := int(_numSectors[halfTrack]) * sectorSize
	g.trackEnd = g.trackStart + trackLength
}

func (g *GCR) readSector(diskData []byte, headerLen uint8, track uint8, sector uint8) ([]uint8, error) {
	if (track < 1) || (track > numTracks) || (sector >= _numSectors[track]) {
		return nil, fmt.Errorf("invalid track/sector")
	}
	offset := (int(_sectorOffset[track]) + int(sector)) << 8
	start := offset + int(headerLen)
	end := start + blockSize
	if start > len(diskData) || end > len(diskData) {
		return nil, fmt.Errorf("invalid track/sector")
	}
	buffer := make([]uint8, blockSize)
	copy(buffer, diskData[start:end])
	return buffer, nil
}

func (g *GCR) sector2gcr(sector []uint8, bam1 uint8, bam2 uint8, trackIdx uint8, sectorIdx uint8) ([]uint8, error) {
	if len(sector) > blockSize {
		return nil, fmt.Errorf("invalid block length")
	}
	idx := 0
	ret := make([]uint8, sectorSize)
	ret[idx] = 0xff
	idx++

	headerData := conv4to5([4]uint8{0x08, uint8(int(sectorIdx) ^ int(trackIdx) ^ int(bam2) ^ int(bam1)), sectorIdx, trackIdx})
	copy(ret[idx:], headerData[:])
	idx += len(headerData)

	bamData := conv4to5([4]uint8{bam2, bam1, 0x0f, 0x0f})
	copy(ret[idx:], bamData[:])
	idx += len(bamData)

	fillData := [9]uint8{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}
	copy(ret[idx:], fillData[:])
	idx += len(fillData)

	ret[idx] = 0xff // Create SYNC
	idx++

	dataMark := conv4to5([4]uint8{0x07, sector[0], sector[1], sector[2]})
	copy(ret[idx:], dataMark[:])
	idx += len(dataMark)

	checksum := sector[0] ^ sector[1] ^ sector[2]
	for x := 3; x < 255; x += 4 {
		data := conv4to5([4]uint8{sector[x], sector[x+1], sector[x+2], sector[x+3]})
		copy(ret[idx:], data[:])
		idx += len(data)

		checksum ^= sector[x]
		checksum ^= sector[x+1]
		checksum ^= sector[x+2]
		checksum ^= sector[x+3]
	}
	checksum ^= sector[255]

	checksumData := conv4to5([4]uint8{sector[255], checksum, 0, 0})
	copy(ret[idx:], checksumData[:])
	idx += len(checksumData)

	endData := [8]uint8{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}
	copy(ret[idx:], endData[:])
	return ret, nil
}

/*
func (g *GCR) sector2gcr_old(block []uint8, bam1 uint8, bam2 uint8, track int, sector int) error {
	if len(block) > blockSize {
		return fmt.Errorf("invalid block length")
	}
	idx := ((track - 1) * trackSize) + (sector * sectorSize)
	g.data[idx] = 0xff
	idx++
	for z, v := range conv4to5([4]uint8{0x08, uint8(sector ^ track ^ int(bam2) ^ int(bam1)), uint8(sector), uint8(track)}) {
		g.data[idx+z] = v
	}
	idx += 5
	for z, v := range conv4to5([4]uint8{bam2, bam1, 0x0f, 0x0f}) {
		g.data[idx+z] = v
	}
	idx += 5
	for x := 0; x < 9; x++ {
		g.data[idx+x] = 0x55
	}
	idx += 9
	// Create GCR data + SYNC
	g.data[idx] = 0xff
	idx++
	// Data mark
	for z, v := range conv4to5([4]uint8{0x07, block[0], block[1], block[2]}) {
		g.data[idx+z] = v
	}
	checksum := block[0]
	checksum ^= block[1]
	checksum ^= block[2]
	idx += 5
	for x := 3; x < 255; x += 4 {
		b0 := block[x]
		b1 := block[x+1]
		b2 := block[x+2]
		b3 := block[x+3]
		for z, v := range conv4to5([4]uint8{b0, b1, b2, b3}) {
			g.data[idx+z] = v
		}
		checksum ^= b0
		checksum ^= b1
		checksum ^= b2
		checksum ^= b3
		idx += 5
	}
	checksum ^= block[255]
	for z, v := range conv4to5([4]uint8{block[255], checksum, 0, 0}) {
		g.data[idx+z] = v
	}
	idx += 5
	for x := 0; x < 8; x++ {
		g.data[idx+x] = 0x55
	}
	return nil
}
*/

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
