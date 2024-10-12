package gcr

import (
	"fmt"
)

const (
	blockBytesLen = 256
	dataBlockLen  = blockBytesLen + 4 // data block id (0x07) + blockBytes + checksum + 0x00 + 0x00
	// All multiple of 5
	syncLen      = 5
	headerLen    = 10
	gapBeginLen  = 9
	gapEndLen    = 8
	gcrDataLen   = (dataBlockLen / 4) * 5
	gcrSectorLen = syncLen + headerLen + gapBeginLen + syncLen + gcrDataLen + gapEndLen
	//realSectorLen = syncLen + headerLen + fillBeginLen + syncLen + dataBlockLen + fillEndLen
)

//real size: 7

type Track struct {
	trackIdx uint8
	data     []uint8
	sectors  uint8
	cursor   uint32
}

func NewTrack(trackIdx uint8, sectors uint8) *Track {
	t := &Track{
		trackIdx: trackIdx,
		sectors:  sectors,
		data:     nil,
		cursor:   0,
	}
	t.data = make([]uint8, int(t.sectors)*gcrSectorLen)
	return t
}

func (t *Track) Len() int {
	return len(t.data)
}

func (t *Track) Sectors() uint8 {
	return t.sectors
}

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

func (t *Track) Cursor() uint32 {
	return t.cursor
}

func (t *Track) Reset(cursor uint32) {
	t.cursor = cursor % uint32(len(t.data))
}

func (t *Track) Advance() {
	t.cursor++
	if t.cursor >= uint32(len(t.data)) {
		t.cursor = 0
	}
}

func (t *Track) Read() uint8 {
	return t.data[t.cursor]
}

func (t *Track) Next() uint8 {
	cursor := t.cursor + 1
	if cursor >= uint32(len(t.data)) {
		cursor = 0
	}
	return t.data[cursor]
}

func (t *Track) Write(data uint8) {
	t.data[t.cursor] = data
}
