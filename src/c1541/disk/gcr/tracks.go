package gcr

type TrackData struct {
	id               uint8
	sectors          uint8
	speedZone        uint8
	bytesPerMicroSec uint8
	rawKBit          float64
	offset           uint16
}

func NewTrackData(id uint8, offset uint16) *TrackData {
	return &TrackData{
		id:     id,
		offset: offset,
	}
}

func (t *TrackData) Update(sectors uint8, speedZone uint8, bytesPerMicroSec uint8, rawKBit float64) {
	t.sectors = sectors
	t.speedZone = speedZone
	t.bytesPerMicroSec = bytesPerMicroSec
	t.rawKBit = rawKBit
}

var _tracks []*TrackData
var _totalSectors uint

func init() {
	_totalSectors = 0
	var currentOffset uint16
	for x := uint8(0); x <= 35; x++ {
		track := NewTrackData(x, currentOffset)
		if x == 0 {
			track.Update(0, 0, 0, 0)
		} else if x >= 1 && x <= 17 {
			track.Update(21, 3, 26, 60.0)
		} else if x >= 18 && x <= 24 {
			track.Update(19, 2, 28, 55.8)
		} else if x >= 25 && x <= 30 {
			track.Update(18, 1, 30, 52.1)
		} else if x >= 31 && x <= 35 {
			track.Update(17, 0, 32, 48.8)
		}
		currentOffset += uint16(track.sectors)
		_totalSectors += uint(track.sectors)
		_tracks = append(_tracks, track)
	}
}

func getImageSize() uint {
	return _totalSectors * blockSize
}

func getTrackStart() uint8 {
	return 1
}

func getTrackCount() uint8 {
	return uint8(len(_tracks))
}

func getTrackSectors(idx uint8) uint8 {
	if idx >= uint8(len(_tracks)) {
		return 0
	}
	return _tracks[idx].sectors
}

func getTrackOffset(idx uint8) uint16 {
	if idx >= uint8(len(_tracks)) {
		return 0
	}
	return _tracks[idx].offset
}
