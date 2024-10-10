package gcr

import (
	"fmt"
	"log"
)

//https://www.pagetable.com/?p=1070

type Disk struct {
	//errorInfo        []uint8
	halfTracks       []*Track
	currentHalfTrack *Track
	usable           bool
}

func NewDisk(image []uint8) (*Disk, error) {
	const bamTrackIdx = 18
	const bamSectorIdx = 0
	if uint(len(image)) < getImageSize() {
		return nil, fmt.Errorf("invalid disk data length")
	}
	hLen := uint8(0)
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		hLen = 64
	}
	bamSector, bErr := rawSector(image, hLen, getTrackOffset(bamTrackIdx), bamSectorIdx)
	if bErr != nil {
		return nil, bErr
	}
	g := &Disk{
		//errorInfo:        make([]uint8, numSectors),
		usable: false,
	}
	bam1 := bamSector[162]
	bam2 := bamSector[163]
	startTrack := getTrackStart()
	halfTracks := getTrackCount() * 2
	g.halfTracks = make([]*Track, halfTracks+startTrack)
	for trackIdx, halfTrackIdx := startTrack, startTrack; halfTrackIdx <= halfTracks; trackIdx, halfTrackIdx = trackIdx+1, halfTrackIdx+2 {
		track := NewTrack(trackIdx, getTrackSectors(trackIdx))
		if tErr := track.Load(image, hLen, bam1, bam2, getTrackOffset(trackIdx)); tErr != nil {
			return nil, tErr
		}
		g.halfTracks[halfTrackIdx] = track
		g.halfTracks[halfTrackIdx+1] = track
	}
	g.currentHalfTrack = g.halfTracks[startTrack]
	g.usable = true
	return g, nil
}

func (g *Disk) Usable() bool {
	return g.usable
}

func (g *Disk) SetHeadHalfTrack(halfTrack uint8) int {
	if halfTrack >= uint8(len(g.halfTracks)) {
		log.Printf("invalid half track: %d", halfTrack)
		return -1
	}
	//Simulate all track rotation
	cursor := g.currentHalfTrack.Cursor()
	g.currentHalfTrack = g.halfTracks[halfTrack]
	g.currentHalfTrack.Reset(cursor)
	return g.currentHalfTrack.Len()
}

func (g *Disk) TrackLen() int {
	return g.currentHalfTrack.Len()
}

func (g *Disk) MicroSecPerByte() uint8 {
	return getMicroSecPerByte(g.currentHalfTrack.trackIdx)
}

func (g *Disk) TrackSectors() uint8 {
	return g.currentHalfTrack.Sectors()
}

func (g *Disk) Rotate() {
	g.currentHalfTrack.Advance()
}

func (g *Disk) Read() uint8 {
	return g.currentHalfTrack.Read()
}

func (g *Disk) Next() uint8 {
	return g.currentHalfTrack.Next()
}

func (g *Disk) Write(data uint8) {
	g.currentHalfTrack.Write(data)
}
