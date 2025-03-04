package gcr

import (
	"fmt"
	"log"
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

// rawSector extracts a specific sector from a disk image based on track offset, header length, and sector index.
// Returns a sector-sized buffer and an error if the extraction fails due to bounds or other inconsistencies.
func rawSector(disk []uint8, headerLen uint8, trackOffset uint16, sectorIdx uint8) ([blockBytesLen]uint8, error) {
	var buffer [blockBytesLen]uint8
	rOffset := (int(trackOffset) + int(sectorIdx)) << 8
	begin := rOffset + int(headerLen)
	end := begin + blockBytesLen
	if begin > len(disk) || end > len(disk) {
		log.Printf("invalid start/end: %d - %d", begin, end)
		return buffer, fmt.Errorf("sector index out of range")
	}
	copy(buffer[:], disk[begin:end])
	return buffer, nil
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
	t.data = make([]uint8, int(t.sectors)*gcrSectorLen)
	if !overlap {
		for i := range t.data {
			t.data[i] = gapByte
		}
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

// Load reads sectors from a disk, converts them to GCR format, and stores the data in the track's buffer.
func (t *Track) Load(disk []uint8, headerLen, bam1 uint8, bam2 uint8, trackOffset uint16) error {
	for sectorIdx := uint8(0); sectorIdx < t.sectors; sectorIdx++ {
		sector, err := rawSector(disk, headerLen, trackOffset, sectorIdx)
		if err != nil {
			return fmt.Errorf("failed to read sector %d: %w", sectorIdx, err)
		}
		gcr := sector2gcr(sector, bam1, bam2, t.trackIdx, sectorIdx)
		start := uint(sectorIdx) * uint(len(gcr))
		end := start + uint(len(gcr))
		if start >= uint(len(disk)) || end > uint(len(disk)) {
			return fmt.Errorf("sector index out of range")
		}
		copy(t.data[start:end], gcr[:])
	}
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
