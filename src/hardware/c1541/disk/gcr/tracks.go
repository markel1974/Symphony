package gcr

// TrackData represents track information including id, sectors, speed zone, transfer characteristics, and offset data.
type TrackData struct {
	id              uint8
	sectors         uint8
	speedZone       uint8
	microSecPerByte uint8
	rawKBit         float64
	offset          uint16
}

// NewTrackData initializes a new TrackData structure with the specified id and offset.
func NewTrackData(id uint8, offset uint16) *TrackData {
	return &TrackData{
		id:     id,
		offset: offset,
	}
}

// Update modifies the TrackData fields with new provided values for sectors, speedZone, microSecPerByte, and rawKBit.
func (t *TrackData) Update(sectors uint8, speedZone uint8, microSecPerByte uint8, rawKBit float64) {
	t.sectors = sectors
	t.speedZone = speedZone
	t.microSecPerByte = microSecPerByte
	t.rawKBit = rawKBit
}

// _tracks stores a list of TrackData pointers representing track metadata, including sectors, speed zones, and offsets.
var _tracks []*TrackData

// _totalSectors holds the total count of sectors across all tracks, computed during the initialization process.
var _totalSectors uint

// init initializes the track data and updates track properties based on predefined ranges and configurations.
func init() {
	_totalSectors = 0
	var currentOffset uint16
	for x := uint8(0); x <= 35; x++ {
		track := NewTrackData(x, currentOffset)
		if x == 0 {
			track.Update(0, 0, 0, 0)
		} else if x >= 1 && x <= 17 {
			track.Update(21, 3, 20 /*19*/ /*26*/, 60.0)
		} else if x >= 18 && x <= 24 {
			track.Update(19, 2, 20 /*28*/, 55.8)
		} else if x >= 25 && x <= 30 {
			track.Update(18, 1, 20 /*30*/, 52.1)
		} else if x >= 31 && x <= 35 {
			track.Update(17, 0, 20 /*32*/, 48.8)
		}
		currentOffset += uint16(track.sectors)
		_totalSectors += uint(track.sectors)
		_tracks = append(_tracks, track)
	}
}

// getImageSize calculates the total size of the image in bytes based on the number of sectors and bytes per sector.
func getImageSize() uint {
	return _totalSectors * blockBytesLen
}

// getTrackStart returns the index of the first usable track on the disk.
func getTrackStart() uint8 {
	return 1
}

// getTrackCount calculates and returns the total number of tracks in the disk as a uint8 value.
func getTrackCount() uint8 {
	return uint8(len(_tracks))
}

// getTrackSectors returns the number of sectors for the specified track index. If the index is invalid, it returns 0.
func getTrackSectors(idx uint8) uint8 {
	if idx >= uint8(len(_tracks)) {
		return 0
	}
	return _tracks[idx].sectors
}

// getTrackOffset returns the offset of a track within the disk image based on the track index.
// If the index is out of range, it returns 0.
func getTrackOffset(idx uint8) uint16 {
	if idx >= uint8(len(_tracks)) {
		return 0
	}
	return _tracks[idx].offset
}

// getMicroSecPerByte returns the number of microseconds required to process a single byte for the given track index.
// If the provided index exceeds the available tracks, it returns 0.
func getMicroSecPerByte(idx uint8) uint8 {
	if idx >= uint8(len(_tracks)) {
		return 0
	}
	return _tracks[idx].microSecPerByte
}
