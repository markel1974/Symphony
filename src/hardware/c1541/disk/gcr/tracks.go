package gcr

// TrackData represents track information including id, sectors, speed zone, transfer characteristics, and offset data.
type TrackData struct {
	id              uint8
	sectors         uint8
	microSecPerByte int
	rawKBit         float64
	offset          uint16
}

// NewTrackData initializes a new TrackData structure with the specified id and offset.
func NewTrackData(id uint8, offset uint16, sectors uint8, microSecPerByte int) *TrackData {
	return &TrackData{
		id:              id,
		offset:          offset,
		sectors:         sectors,
		microSecPerByte: microSecPerByte,
	}
}
