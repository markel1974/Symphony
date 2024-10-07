package gcr

/*
const sync = 0xff
const fill uint8 = 0x55

type Sector struct {
	trackNum  uint8
	sectorNum uint8
	begin     int
	data      [SectorSize]uint8
	cursor    uint
}

func NewSector(trackNum uint8, sectorNum uint8, headerLen uint8, offset uint16) *Sector {
	rOffset := (int(offset) + int(sectorNum)) << 8
	begin := rOffset + int(headerLen)
	s := &Sector{
		trackNum:  trackNum,
		sectorNum: sectorNum,
		cursor:    0,
		begin:     begin,
	}
	for i := range s.data {
		s.data[i] = fill
	}
	return s
}

func (s *Sector) Len() int {
	return len(s.data)
}

func (s *Sector) Load(image []uint8, bam1 uint8, bam2 uint8) error {
	sector, err := s.Raw(image)
	if err != nil {
		return err
	}
	s.data = s.sector2gcr(sector, bam1, bam2)
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
	//if data == sync {
	//	log.Printf("Write: Track %d -> Sector %d -> Cursor %d = %02x", s.trackNum, s.sectorNum, s.cursor, data)
	//}
	s.data[s.cursor] = data
}

func (s *Sector) Raw(diskData []byte) ([BlockSize]uint8, error) {
	var buffer [BlockSize]uint8
	end := s.begin + BlockSize
	if s.begin > len(diskData) || end > len(diskData) {
		log.Printf("invalid start/end: %d - %d", s.begin, end)
		return buffer, fmt.Errorf("invalid start/end")
	}
	copy(buffer[:], diskData[s.begin:end])
	return buffer, nil
}

func (s *Sector) sector2gcr(sector [BlockSize]uint8, bam1 uint8, bam2 uint8) [SectorSize]uint8 {
	const last = BlockSize - 1

	var ret [SectorSize]uint8
	idx := 0
	ret[idx] = sync
	idx++

	headerData := conv4to5([4]uint8{0x08, uint8(int(s.sectorNum) ^ int(s.trackNum) ^ int(bam2) ^ int(bam1)), s.sectorNum, s.trackNum})
	copy(ret[idx:], headerData[:])
	idx += len(headerData)

	bamData := conv4to5([4]uint8{bam2, bam1, 0x0f, 0x0f})
	copy(ret[idx:], bamData[:])
	idx += len(bamData)

	fillData := [9]uint8{fill, fill, fill, fill, fill, fill, fill, fill, fill}
	copy(ret[idx:], fillData[:])
	idx += len(fillData)

	ret[idx] = sync
	idx++

	dataMark := conv4to5([4]uint8{0x07, sector[0], sector[1], sector[2]})
	copy(ret[idx:], dataMark[:])
	idx += len(dataMark)

	checksum := sector[0] ^ sector[1] ^ sector[2]
	for x := 3; x < last; x += 4 {
		data := conv4to5([4]uint8{sector[x], sector[x+1], sector[x+2], sector[x+3]})
		copy(ret[idx:], data[:])
		idx += len(data)

		checksum ^= sector[x]
		checksum ^= sector[x+1]
		checksum ^= sector[x+2]
		checksum ^= sector[x+3]
	}
	checksum ^= sector[last]

	checksumData := conv4to5([4]uint8{sector[last], checksum, 0, 0})
	copy(ret[idx:], checksumData[:])
	idx += len(checksumData)

	endData := [8]uint8{fill, fill, fill, fill, fill, fill, fill, fill}
	copy(ret[idx:], endData[:])
	return ret
}

*/

/*
///Users/tinmr305/Desktop/emu/vice-emu-code-r45201-trunk-vice/src/gcr.c
func (g *GCR) gcr2sector(block []uint8, track int, sector int) []uint8 {
	var shift, i, j int
	var gcr[5] uint8
	var b uint8;
	var offsetIdx uint8
	var end uint8 = raw.data + raw.size;

	shift = p & 7;
	offsetIdx = raw->data + (p >> 3);

	b = offset[0] << shift;
	for i = 0; i < num; i++, buf += 4 {
		// get 5 bytes of gcr data
		for j = 0; j < 5; j++){
			offset++;
			if offset >= end {
				offset = raw->data;
			}
			if shift {
				gcr[j] = b | ((offset[0] << shift) >> 8);
				b = offset[0] << shift;
			} else {
				gcr[j] = b;
				b = offset[0];
			}
		}
		conv4to5(gcr, buf);
	}
}
*/
