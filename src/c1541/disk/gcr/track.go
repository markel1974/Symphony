package gcr

import "fmt"

type SectorItem struct {
	sector *Sector
	next   *SectorItem
}

type Track struct {
	trackNum  uint8
	offset    uint16
	headerLen uint8
	sectors   []*SectorItem
	current   *SectorItem
}

func NewTrack(trackNum uint8, headerLen uint8) *Track {
	t := &Track{
		trackNum:  trackNum,
		headerLen: headerLen,
		offset:    _sectorOffset[trackNum],
		sectors:   make([]*SectorItem, _numSectors[trackNum]),
	}
	for sectorNum := range t.sectors {
		item := &SectorItem{sector: NewSector(t.trackNum, uint8(sectorNum), t.headerLen, t.offset), next: nil}
		t.sectors[sectorNum] = item
		if sectorNum > 0 {
			t.sectors[sectorNum-1].next = item
		}
	}
	t.sectors[len(t.sectors)-1].next = t.sectors[0]
	t.current = t.sectors[0]
	return t
}

func (t *Track) RawSector(disk []uint8, sectorNum uint8) ([]uint8, error) {
	if sectorNum >= uint8(len(t.sectors)) {
		return nil, fmt.Errorf("sector index out of range")
	}
	return t.current.sector.Raw(disk)
}

func (t *Track) Load(disk []uint8, bam1 uint8, bam2 uint8) error {
	for _, item := range t.sectors {
		if err := item.sector.Load(disk, bam1, bam2); err != nil {
			return err
		}
	}
	return nil
}

func (t *Track) Advance() {
	if !t.current.sector.Advance() {
		t.current = t.current.next
		t.current.sector.Reset()
	}
}

func (t *Track) Read() uint8 {
	return t.current.sector.Read()
}

func (t *Track) Write(data uint8) {
	t.current.sector.Write(data)
}
