package gcr

import "fmt"

const (
	bamTrack  = 18
	bamSector = 0
)

type GCR struct {
	data             []uint8
	errorInfo        []uint8
	idx              int // Pointer to GCR disk under R/W head
	trackStart       int // Pointer to start of GCR disk of current track
	trackEnd         int // Pointer to end of GCR disk of current track
	currentHalfTrack int
}

func NewGCR() *GCR {
	d := &GCR{
		data:             make([]uint8, DiskSize),
		errorInfo:        make([]uint8, NumSectors),
		idx:              0,
		trackStart:       0,
		trackEnd:         TrackSize,
		currentHalfTrack: 2,
	}
	for x := range d.data {
		d.data[x] = 0x55
	}
	for x := range d.errorInfo {
		d.errorInfo[x] = 1
	}
	return d
}

func (gcr *GCR) Setup(image []uint8) error {
	imageLen := len(image)
	if imageLen < (NumSectors * BlockSize) {
		return fmt.Errorf("invalid disk data length")
	}
	headerLen := 0
	if (image[0] == 0x43) && (image[1] == 0x1) && (image[2] == 0x41) && (image[3] == 0x64) {
		headerLen = 64
	}
	if (headerLen == 0) && (imageLen == (NumSectors * 257)) {
		copy(gcr.errorInfo, image[NumSectors*BlockSize:])
	}
	bam := gcr.readSector(image, headerLen, bamTrack, bamSector)
	if bam == nil {
		return fmt.Errorf("nil bam")
	}
	id1 := bam[162]
	id2 := bam[163]
	// Create GCR encoded disk from image
	for track := 1; track <= NumTracks; track++ {
		for sector := 0; sector < int(_numSectors[track]); sector++ {
			if block := gcr.readSector(image, headerLen, track, sector); block != nil {
				if err := gcr.convertSector(block, id1, id2, track, sector); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (gcr *GCR) MoveOut() {
	if gcr.currentHalfTrack <= 2 {
		return
	}
	gcr.currentHalfTrack--
	gcr.updateTrack()
}

func (gcr *GCR) MoveIn() {
	if gcr.currentHalfTrack >= NumTracksMax {
		return
	}
	gcr.currentHalfTrack++
	gcr.updateTrack()
}

func (gcr *GCR) Rotate() {
	gcr.idx++
	if gcr.idx == gcr.trackEnd {
		gcr.idx = gcr.trackStart
	}
}

func (gcr *GCR) Read() uint8 {
	return gcr.data[gcr.idx]
}

func (gcr *GCR) Write(data uint8) {
	gcr.data[gcr.idx] = data
}

//func (gcr *GCR) GetErrorInfo() []uint8 {
//	return gcr.errorInfo
//}

func (gcr *GCR) updateTrack() {
	gcr.idx = ((gcr.currentHalfTrack >> 1) - 1) * TrackSize
	gcr.trackStart = gcr.idx
	trackLength := getTrackLen(gcr.currentHalfTrack >> 1)
	gcr.trackEnd = gcr.trackStart + trackLength
}

func (gcr *GCR) readSector(diskData []byte, headerLen int, track int, sector int) []uint8 {
	offset := offsetFromTrackSector(track, sector)
	if offset < 0 {
		return nil
	}
	start := offset + headerLen
	if end := start + BlockSize; end >= len(diskData) {
		return nil
	}
	buffer := make([]uint8, BlockSize)
	copy(buffer, diskData[start:])
	return buffer
}

func (gcr *GCR) convertSector(block []uint8, id1 uint8, id2 uint8, track int, sector int) error {
	if len(block) > BlockSize {
		return fmt.Errorf("invalid block length")
	}
	// Create GCR header
	// SYNC
	buf := make([]uint8, 4)
	p := ((track - 1) * TrackSize) + (sector * SectorSize)
	gcr.data[p] = 0xff
	p++
	// Header mark
	buf[0] = 0x08
	// Checksum
	buf[1] = uint8(sector ^ track ^ int(id2) ^ int(id1))
	buf[2] = uint8(sector)
	buf[3] = uint8(track)
	conv4(buf, gcr.data[p:])
	buf[0] = id2
	buf[1] = id1
	buf[2] = 0x0f
	buf[3] = 0x0f
	conv4(buf, gcr.data[p+5:])
	p += 10
	for x := 0; x < 9; x++ {
		gcr.data[p+x] = 0x55
	}
	p += 9
	// Create GCR data
	// SYNC
	gcr.data[p] = 0xff
	p++
	// Data mark
	buf[0] = 0x07
	buf[1] = block[0]
	sum := buf[1]
	buf[2] = block[1]
	sum ^= buf[2]
	buf[3] = block[2]
	sum ^= buf[3]
	conv4(buf, gcr.data[p:])
	p += 5
	for i := 3; i < 255; i += 4 {
		buf[0] = block[i]
		sum ^= buf[0]
		buf[1] = block[i+1]
		sum ^= buf[1]
		buf[2] = block[i+2]
		sum ^= buf[2]
		buf[3] = block[i+3]
		sum ^= buf[3]
		conv4(buf, gcr.data[p:])
		p += 5
	}
	buf[0] = block[255]
	sum ^= buf[0]
	// Checksum
	buf[1] = sum
	buf[2] = 0
	buf[3] = 0
	conv4(buf, gcr.data[p:])
	p += 5
	for x := 0; x < 8; x++ {
		gcr.data[p+x] = 0x55
	}
	return nil
}
