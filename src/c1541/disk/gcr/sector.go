package gcr

import (
	"fmt"
	"log"
)

type Sector struct {
	headerLen uint8
	trackNum  uint8
	sectorNum uint8
	offset    uint16
	data      [sectorSize]uint8
	cursor    uint
}

func NewSector(trackNum uint8, sectorNum uint8, headerLen uint8, offset uint16) *Sector {
	s := &Sector{
		trackNum:  trackNum,
		sectorNum: sectorNum,
		headerLen: headerLen,
		offset:    offset,
		cursor:    0,
	}
	for i := range s.data {
		s.data[i] = 0x55
	}
	return s
}

func (s *Sector) Load(image []uint8, bam1 uint8, bam2 uint8) error {
	if sector, _ := s.Raw(image); sector != nil {
		data, err := s.sector2gcr(sector, bam1, bam2)
		if err != nil {
			return err
		}
		s.data = data
	}
	return nil
}

func (s *Sector) Reset() {
	s.cursor = 0
}

func (s *Sector) Advance() bool {
	if (s.cursor + 1) >= uint(len(s.data)) {
		return false
	}
	s.cursor++
	return true
}

func (s *Sector) Read() uint8 {
	//log.Printf("Read: Track %d -> Sector %d -> Cursor %d", s.trackNum, s.sectorNum, s.cursor)
	return s.data[s.cursor]
}

func (s *Sector) Write(data uint8) {
	//log.Printf("Write: Track %d -> Sector %d -> Cursor %d = %02x", s.trackNum, s.sectorNum, s.cursor, data)
	s.data[s.cursor] = data
}

func (s *Sector) Raw(diskData []byte) ([]uint8, error) {
	if (s.trackNum < 1) || (s.trackNum > numTracks) {
		log.Printf("Invalid track: %d", s.trackNum)
		return nil, fmt.Errorf("invalid track")
	}
	if s.sectorNum >= _numSectors[s.trackNum] {
		log.Printf("Invalid sector: %d", s.sectorNum)
		return nil, fmt.Errorf("invalid sector")
	}
	diskOffset := (int(s.offset) + int(s.sectorNum)) << 8
	start := diskOffset + int(s.headerLen)
	end := start + blockSize
	if start > len(diskData) || end > len(diskData) {
		log.Printf("Invalid offset: %d", s.offset)
		return nil, fmt.Errorf("invalid offset")
	}
	buffer := make([]uint8, blockSize)
	copy(buffer, diskData[start:end])
	return buffer, nil
}

func (s *Sector) sector2gcr(sector []uint8, bam1 uint8, bam2 uint8) ([sectorSize]uint8, error) {
	var ret [sectorSize]uint8
	if len(sector) > blockSize {
		log.Printf("Invalid block length: %d", len(sector))
		return ret, fmt.Errorf("invalid block length")
	}
	idx := 0
	ret[idx] = 0xff
	idx++

	headerData := conv4to5([4]uint8{0x08, uint8(int(s.sectorNum) ^ int(s.trackNum) ^ int(bam2) ^ int(bam1)), s.sectorNum, s.trackNum})
	copy(ret[idx:], headerData[:])
	idx += len(headerData)

	bamData := conv4to5([4]uint8{bam2, bam1, 0x0f, 0x0f})
	copy(ret[idx:], bamData[:])
	idx += len(bamData)

	fillData := [9]uint8{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}
	copy(ret[idx:], fillData[:])
	idx += len(fillData)

	ret[idx] = 0xff // Create SYNC
	idx++

	dataMark := conv4to5([4]uint8{0x07, sector[0], sector[1], sector[2]})
	copy(ret[idx:], dataMark[:])
	idx += len(dataMark)

	checksum := sector[0] ^ sector[1] ^ sector[2]
	for x := 3; x < 255; x += 4 {
		data := conv4to5([4]uint8{sector[x], sector[x+1], sector[x+2], sector[x+3]})
		copy(ret[idx:], data[:])
		idx += len(data)

		checksum ^= sector[x]
		checksum ^= sector[x+1]
		checksum ^= sector[x+2]
		checksum ^= sector[x+3]
	}
	checksum ^= sector[255]

	checksumData := conv4to5([4]uint8{sector[255], checksum, 0, 0})
	copy(ret[idx:], checksumData[:])
	idx += len(checksumData)

	endData := [8]uint8{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55}
	copy(ret[idx:], endData[:])
	return ret, nil
}
