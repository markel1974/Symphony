package gcr

import (
	"fmt"
	"log"
)

//https://www.pagetable.com/?p=1070

type Disk struct {
	//errorInfo        []uint8
	tracks       []*Track
	currentTrack *Track
	usable       bool
}

func NewDisk(image []uint8) (*Disk, error) {
	const bamTrackIdx = 18
	const bamSectorIdx = 0

	g := &Disk{
		//errorInfo:        make([]uint8, numSectors),
		usable: false,
	}

	if uint(len(image)) < getImageSize() {
		return nil, fmt.Errorf("invalid disk data length")
	}
	headerLen := uint8(0)
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		headerLen = 64
	}
	bamSector, bErr := rawSector(image, headerLen, getTrackOffset(bamTrackIdx), bamSectorIdx)
	if bErr != nil {
		return nil, bErr
	}
	bam1 := bamSector[162]
	bam2 := bamSector[163]
	g.tracks = make([]*Track, getTrackCount()+getTrackStart())
	for trackIdx := getTrackStart(); trackIdx <= getTrackCount(); trackIdx++ {
		track := NewTrack(trackIdx, getTrackSectors(trackIdx))
		if tErr := track.Load(image, headerLen, bam1, bam2, getTrackOffset(trackIdx)); tErr != nil {
			return nil, tErr
		}
		g.tracks[trackIdx] = track
	}
	g.currentTrack = g.tracks[getTrackStart()]
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
	//Simulate all track rotation
	cursor := g.currentTrack.Cursor()
	g.currentTrack = g.tracks[track]
	g.currentTrack.Reset(cursor)
	return g.currentTrack.Len()
}

func (g *Disk) TrackLen() int {
	return g.currentTrack.Len()
}

func (g *Disk) MicroSecPerByte() uint8 {
	return getMicroSecPerByte(g.currentTrack.trackIdx)
}

func (g *Disk) TrackSectors() uint8 {
	return g.currentTrack.Sectors()
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
