package gcr

import (
	"fmt"
)

const (
	NumTracks  = 35
	NumSectors = 683
	BlockSize  = 256
	SectorSize = 1 + 10 + 9 + 1 + 325 + 8 // SYNC Header Gap SYNC Data Gap (should be 5 SYNC bytes each)
)

const (
	startTrack   = 1
	bamTrackIdx  = 18
	bamSectorIdx = 0
)

type Disk struct {
	//errorInfo        []uint8
	tracks       []*Track
	currentTrack *Track
	usable       bool
}

func NewDisk(image []uint8) (*Disk, error) {
	g := &Disk{
		//errorInfo:        make([]uint8, numSectors),
		usable: false,
	}
	if len(image) < (NumSectors * BlockSize) {
		return nil, fmt.Errorf("invalid disk data length")
	}
	headerLen := uint8(0)
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		headerLen = 64
	}
	bamTrack := NewTrack(bamTrackIdx, headerLen)
	bamSector, bErr := bamTrack.RawSector(image, bamSectorIdx)
	if bErr != nil {
		return nil, bErr
	}
	bam1 := bamSector[162]
	bam2 := bamSector[163]

	g.tracks = make([]*Track, NumTracks+1)

	for trackNum := uint8(startTrack); trackNum <= NumTracks; trackNum++ {
		track := NewTrack(trackNum, headerLen)
		if tErr := track.Load(image, bam1, bam2); tErr != nil {
			return nil, tErr
		}
		g.tracks[trackNum] = track
	}
	g.currentTrack = g.tracks[startTrack]
	g.usable = true
	return g, nil
}

func (g *Disk) Usable() bool {
	return g.usable
}

func (g *Disk) GetTracksNumber() uint8 {
	return NumTracks
}

func (g *Disk) SetHeadTrack(track uint8) {
	g.currentTrack = g.tracks[track]
}

/*
func (g *Disk) MoveOut() {
	//todo halfTrack handler
	//if g.currentHalfTrack <= 2 {
	//	return
	//}
	//g.currentHalfTrack--
	//track := g.currentHalfTrack >> 1
	//g.currentTrack = g.tracks[track]
}

func (g *Disk) MoveIn() {
	//todo halfTrack handler
	if g.currentHalfTrack >= numHalfTracks {
		return
	}
	g.currentHalfTrack++
	track := g.currentHalfTrack >> 1
	g.currentTrack = g.tracks[track]
}
*/

func (g *Disk) Rotate() {
	g.currentTrack.Advance()
}

func (g *Disk) Read() uint8 {
	return g.currentTrack.Read()
}

func (g *Disk) Write(data uint8) {
	g.currentTrack.Write(data)
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
