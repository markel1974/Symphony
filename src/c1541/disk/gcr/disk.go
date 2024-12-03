package gcr

import (
	"fmt"
	"log"
	"math/rand"
)

//https://www.pagetable.com/?p=1070

type Disk struct {
	//errorInfo        []uint8
	tracks       []*Track
	currentTrack *Track
	usable       bool
}

func NewDisk() *Disk {
	startTrack := getTrackStart()
	tracks := getTrackCount()
	g := &Disk{
		//errorInfo:        make([]uint8, numSectors),
		tracks: make([]*Track, tracks+startTrack),
		usable: true,
	}
	for trackIdx := startTrack; trackIdx <= tracks; trackIdx++ {
		g.tracks[trackIdx] = NewTrack(trackIdx, getTrackSectors(trackIdx), false)
	}
	g.currentTrack = g.tracks[startTrack]
	return g
}

func (g *Disk) Load(image []byte) error {
	const bamTrackIdx = 18
	const bamSectorIdx = 0
	if uint(len(image)) < getImageSize() {
		return fmt.Errorf("invalid disk data length")
	}
	hLen := uint8(0)
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		hLen = 64
	}
	id1 := uint8(rand.Intn(0xff))
	id2 := uint8(rand.Intn(0xff))
	if bam, err := rawSector(image, hLen, getTrackOffset(bamTrackIdx), bamSectorIdx); err == nil {
		id1 = bam[162]
		id2 = bam[163]
	}
	for _, track := range g.tracks {
		if track == nil {
			continue
		}
		if err := track.Load(image, hLen, id1, id2, getTrackOffset(track.Index())); err != nil {
			return err
		}
	}
	g.currentTrack = g.tracks[getTrackStart()]
	return nil
}

func (g *Disk) Usable() bool {
	return g.usable
}

func (g *Disk) SetHeadHalfTrack(halfTrack uint8) int {
	track := halfTrack >> 1
	if track >= uint8(len(g.tracks)) {
		log.Printf("invalid half track: %d", halfTrack)
		return -1
	}
	//Simulate all track rotation
	cursor := g.currentTrack.Cursor()
	g.currentTrack.Leave()

	g.currentTrack = g.tracks[track]
	g.currentTrack.Enter(cursor)
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

func (g *Disk) Next() uint8 {
	return g.currentTrack.Next()
}

func (g *Disk) Write(data uint8) {
	g.currentTrack.Write(data)
}
