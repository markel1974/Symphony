package gcr

import (
	"bytes"
	"fmt"
	"math/rand"
)

//real size: 7

// Track represents a data structure for managing a disk's track with its associated data and state.
type Track struct {
	trackIdx uint8
	valid    bool
	readable bool
	overlap  bool
	data     []uint8
	sectors  uint8
	cursor   uint32
	rng      *rand.Rand
}

// NewTrack initializes a new Track with the specified track index, sectors, and overlap settings, allocating its data buffer.
func NewTrack(trackIdx uint8, valid bool, sectors uint8, readable bool, overlap bool) *Track {
	t := &Track{
		trackIdx: trackIdx,
		valid:    valid,
		sectors:  sectors,
		readable: readable,
		overlap:  overlap,
		data:     nil,
		cursor:   0,
		rng:      rand.New(rand.NewSource(int64(trackIdx))),
	}
	return t
}

func (t *Track) Valid() bool {
	return t.valid
}

func (t *Track) Readable() bool {
	return t.readable
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

// Load initializes the track data by processing sectors from the provided disk image and allocating gaps where necessary.
// diskImage represents the entire disk in bytes; headerLen defines the size of the header per sector.
// id1 and id2 specify the unique disk IDs; trackOffset determines the starting byte for this track.
// It calculates interleave, creates GCR-encoded sectors, and adds inter-sector and tail gaps.
// Returns an error if sector extraction or processing fails.
func (t *Track) Load(diskImage []byte, headerLen uint8, id1 uint8, id2 uint8, trackOffset uint16) error {
	numSectors, _, totalTrackBytes := getTrackInfo(t.trackIdx)
	if numSectors == 0 {
		return nil
	}
	totalUsedBySectors := int(numSectors) * gcrSectorLen
	tailGap := []byte(nil)
	sectorGap := []byte(nil)
	if totalEmptySpace := totalTrackBytes - totalUsedBySectors; totalEmptySpace > 0 {
		tailGapSize := totalEmptySpace
		if tailGapSize > gcrSectorLen {
			fmt.Printf("track %d exceeding tail gap limit (%d), resizing to %d\n", t.trackIdx, tailGapSize, gcrSectorLen)
			tailGapSize = gcrSectorLen
			spaceToDistribute := totalEmptySpace - tailGapSize
			if spaceToDistribute > 0 {
				sectorGapSize := spaceToDistribute / int(numSectors)
				sectorGap = createGap(sectorGapSize)
				//fmt.Printf("track %d -> sector gap %d\n", t.trackIdx, len(sectorGap))
			}
		}
		tailGap = createGap(tailGapSize)
		//fmt.Printf("track %d -> tail gap %d\n", t.trackIdx, len(sectorGap))
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

// ApplyNoise introduces random bit-level corruption to the track data at predefined probabilities for data and gap bytes.
func (t *Track) ApplyNoise() {
	const dataCorruptionChance float32 = 0.03 // 3% probabilità per i byte di dati/header
	const gapCorruptionChance float32 = 0.10  // 10% probabilità per i byte di gap
	const syncMark byte = 0xFF
	const gapFill byte = 0x55
	for i, b := range t.data {
		if b == syncMark {
			continue
		}
		var chance float32
		if b == gapFill {
			chance = gapCorruptionChance
		} else {
			chance = dataCorruptionChance
		}
		if t.rng.Float32() < chance {
			bitToFlip := uint8(1 << t.rng.Intn(8))
			t.data[i] = b ^ bitToFlip
		}
	}
}

// Cursor returns the current cursor position of the track as an unsigned 32-bit integer.
func (t *Track) Cursor() uint32 {
	return t.cursor
}

// Enter sets the cursor to the specified position modulo the track data length and resets the write count to zero.
func (t *Track) Enter(cursor uint32) {
	t.cursor = cursor % uint32(len(t.data))
}

// Leave resets the write count of the track to zero, indicating no pending writes on the track.
func (t *Track) Leave() {
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
}
