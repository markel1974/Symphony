package gcr

import (
	"fmt"
	"log"
)

//https://www.pagetable.com/?p=1070

// _numSectorsPerTrack  Number of sectors of each track
var _numSectorsPerTrack = []uint8{
	0, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	19, 19, 19, 19, 19, 19, 19, 18, 18, 18, 18, 18, 18, 17, 17, 17, 17, 17,
}

var _sectorOffsetPerTrack []uint16
var _totalSectors int

func init() {
	_totalSectors = 0
	_sectorOffsetPerTrack = make([]uint16, len(_numSectorsPerTrack))
	var totalOffset uint16
	for x, v := range _numSectorsPerTrack {
		_totalSectors += int(v)
		_sectorOffsetPerTrack[x] = totalOffset
		totalOffset += uint16(v)
	}
}

type Disk struct {
	//errorInfo        []uint8
	tracks       []*Track
	currentTrack *Track
	usable       bool
}

func NewDisk(image []uint8) (*Disk, error) {
	const numTracks = 35
	const startTrack = 1
	const bamTrackIdx = 18
	const bamSectorIdx = 0

	g := &Disk{
		//errorInfo:        make([]uint8, numSectors),
		usable: false,
	}

	if len(image) < (_totalSectors * blockSize) {
		return nil, fmt.Errorf("invalid disk data length")
	}

	headerLen := uint8(0)
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		headerLen = 64
	}
	//bamTrack := NewTrack(bamTrackIdx, headerLen, _numSectorsPerTrack[bamTrackIdx], _sectorOffsetPerTrack[bamTrackIdx])
	bamSector, bErr := rawSector(image, headerLen, _sectorOffsetPerTrack[bamTrackIdx], bamSectorIdx)
	if bErr != nil {
		return nil, bErr
	}
	bam1 := bamSector[162]
	bam2 := bamSector[163]

	g.tracks = make([]*Track, numTracks+1)
	for trackIdx := uint8(startTrack); trackIdx <= numTracks; trackIdx++ {
		track := NewTrack(trackIdx, _numSectorsPerTrack[trackIdx])
		if tErr := track.Load(image, headerLen, bam1, bam2, _sectorOffsetPerTrack[trackIdx]); tErr != nil {
			return nil, tErr
		}
		g.tracks[trackIdx] = track
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
	//Simulate all track rotation
	cursor := g.currentTrack.Cursor()
	g.currentTrack = g.tracks[track]
	g.currentTrack.Reset(cursor)
	return g.currentTrack.Len()
}

func (g *Disk) TrackLen() int {
	return g.currentTrack.Len()
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
