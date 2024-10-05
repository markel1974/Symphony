package gcr

import (
	"fmt"
	"log"
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

func (g *Disk) SetHeadTrack(track uint8) int {
	if track >= uint8(len(g.tracks)) {
		log.Printf("invalid track number: %d", track)
		return -1
	}
	g.currentTrack = g.tracks[track]
	return g.currentTrack.Len()
}

func (g *Disk) HeadTrackLen() int {
	return g.currentTrack.Len()
}

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
