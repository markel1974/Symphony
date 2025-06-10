package gcr

import (
	"fmt"
	"log"
	"math/rand"
)

//https://www.pagetable.com/?p=1070

// Disk represents a storage medium comprising multiple tracks that can store and manage data.
// It maintains the current track and indicates whether the disk is usable.
type Disk struct {
	//errorInfo        []uint8
	tracks       []*Track
	currentTrack *Track
	usable       bool
	wp           bool
}

// NewDisk initializes a new Disk instance with tracks and sets it as usable with the start track as the current track.
func NewDisk(wp bool) *Disk {
	startTrack := getTrackStart()
	tracks := getTrackCount()
	g := &Disk{
		wp:     wp,
		tracks: make([]*Track, tracks+startTrack),
		usable: true,
		//errorInfo:        make([]uint8, numSectors),
	}
	for trackIdx := startTrack; trackIdx <= tracks; trackIdx++ {
		g.tracks[trackIdx] = NewTrack(trackIdx, getTrackSectors(trackIdx), false)
	}
	g.currentTrack = g.tracks[startTrack]
	return g
}

func (e *Disk) WriteProtected() bool {
	return e.wp
}

// Load initializes the Disk by reading from the provided image data and loading each track with its sectors and metadata.
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

// Usable returns the usability status of the disk, indicating whether it is functional or not.
func (g *Disk) Usable() bool {
	return g.usable
}

// SetHeadHalfTrack sets the disk head to the specified half-track position and returns the length of the target track.
// If the half-track is invalid, it logs an error and returns -1. The function also retains the cursor position on the track.
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

// TrackLen returns the length of the current track's data in bytes.
func (g *Disk) TrackLen() int {
	return g.currentTrack.Len()
}

// MicroSecPerByte returns the number of microseconds required to process a single byte for the current track.
func (g *Disk) MicroSecPerByte() int {
	return getMicroSecPerByte(g.currentTrack.trackIdx)
}

// TrackSectors returns the number of sectors in the current track of the disk.
func (g *Disk) TrackSectors() uint8 {
	return g.currentTrack.Sectors()
}

// Rotate advances the cursor of the current track to the next position, simulating the rotation of the disk.
func (g *Disk) Rotate() {
	g.currentTrack.Advance()
}

// Read returns the current byte at the cursor position of the current track.
func (g *Disk) Read() uint8 {
	return g.currentTrack.Read()
}

// Next retrieves the value at the next position of the current track's cursor without advancing it.
func (g *Disk) Next() uint8 {
	return g.currentTrack.Next()
}

// Write writes the specified byte of data to the current position in the current track of the disk.
func (g *Disk) Write(data uint8) {
	g.currentTrack.Write(data)
}
