package gcr

import (
	"bytes"
)

//real size: 7

// Track represents a data structure for managing a disk's track with its associated data and state.
type Track struct {
	trackIdx   uint8
	overlap    bool
	data       []uint8
	sectors    uint8
	cursor     uint32
	writeCount int
}

// NewTrack initializes a new Track with the specified track index, sectors, and overlap settings, allocating its data buffer.
func NewTrack(trackIdx uint8, sectors uint8, overlap bool) *Track {
	t := &Track{
		trackIdx:   trackIdx,
		sectors:    sectors,
		overlap:    overlap,
		data:       nil,
		cursor:     0,
		writeCount: 0,
	}
	return t
}

// Index returns the index of the track as a uint8.
func (t *Track) Index() uint8 {
	return t.trackIdx
}

// Overlap returns true if the track overlaps with another track, otherwise false.
func (t *Track) Overlap() bool {
	return t.overlap
}

// Len returns the number of bytes stored in the track's data slice.
func (t *Track) Len() int {
	return len(t.data)
}

// Sectors returns the number of sectors in the track.
func (t *Track) Sectors() uint8 {
	return t.sectors
}

func (t *Track) Load(diskImage []byte, headerLen uint8, id1 uint8, id2 uint8, trackOffset uint16) error {
	const maxTailGap = 380 //400  timeout limit
	//const maxTailGap = 1000000000
	numSectors, cyclesPerByte := getTrackInfo(t.trackIdx)
	if numSectors == 0 {
		return nil
	}
	totalTrackBytes := int(rotationTimeCycles / cyclesPerByte)
	totalUsedBySectors := int(numSectors) * gcrSectorLen
	tailGap := []byte(nil)
	sectorGap := []byte(nil)
	if totalEmptySpace := totalTrackBytes - totalUsedBySectors; totalEmptySpace > 0 {
		tailGapSize := totalEmptySpace
		if tailGapSize > maxTailGap {
			tailGapSize = maxTailGap
			spaceToDistribute := totalEmptySpace - tailGapSize
			if spaceToDistribute > 0 {
				sectorGapSize := spaceToDistribute / int(numSectors)
				sectorGap = createGap(sectorGapSize)
			}
		}
		tailGap = createGap(tailGapSize)
	}
	var trackBuffer bytes.Buffer
	trackBuffer.Grow(totalTrackBytes)
	sectorTailGapIdx, interleaveTable := getInterleaveTable(t.trackIdx, numSectors)

	for _, sectorIdx := range interleaveTable {
		sectorData, err := rawSector(diskImage, headerLen, trackOffset, sectorIdx)
		if err != nil {
			sectorData = [blockBytesLen]uint8{}
		}
		gcrBlock := sector2gcr(sectorData, id1, id2, t.trackIdx, sectorIdx)
		trackBuffer.Write(gcrBlock[:])
		if len(sectorGap) > 0 {
			trackBuffer.Write(sectorGap)
		}
		if sectorIdx == sectorTailGapIdx && len(tailGap) > 0 {
			trackBuffer.Write(tailGap)
		}
	}
	finalAdjustment := totalTrackBytes - trackBuffer.Len()
	if finalAdjustment > 0 {
		extraGap := createGap(finalAdjustment)
		trackBuffer.Write(extraGap)
	}
	t.data = trackBuffer.Bytes()
	return nil
}

// Cursor returns the current cursor position of the track as an unsigned 32-bit integer.
func (t *Track) Cursor() uint32 {
	return t.cursor
}

// Enter sets the cursor to the specified position modulo the track data length and resets the write count to zero.
func (t *Track) Enter(cursor uint32) {
	t.cursor = cursor % uint32(len(t.data))
	t.writeCount = 0
	//fmt.Println("ENTERING TRACK", t.Index(), len(t.data))
}

// Leave resets the write count of the track to zero, indicating no pending writes on the track.
func (t *Track) Leave() {
	//if t.writeCount > 0 {
	//	fmt.Println("LEAVING TRACK", t.Index(), len(t.data), t.writeCount)
	//}
	t.writeCount = 0
}

// Advance increments the cursor position within the track and wraps around to the start if the end is reached.
func (t *Track) Advance() {
	t.cursor++
	if t.cursor >= uint32(len(t.data)) {
		t.cursor = 0
	}
}

// Read retrieves the byte at the current cursor position in the track's data buffer.
func (t *Track) Read() uint8 {
	return t.data[t.cursor]
}

// Next retrieves the value at the position immediately following the current cursor without advancing the cursor.
func (t *Track) Next() uint8 {
	cursor := t.cursor + 1
	if cursor >= uint32(len(t.data)) {
		cursor = 0
	}
	return t.data[cursor]
}

// Write writes the specified byte of data to the current cursor position in the track.
func (t *Track) Write(data uint8) {
	t.data[t.cursor] = data
	t.writeCount++
	//if t.Index() == 18 {
	//	fmt.Printf("Write %d -> %02x\n", t.Cursor(), data)
	//}
	//TODO
	//fmt.Printf("%d - %d Write %02x\n", t.Index(), t.Cursor(), data)
}
