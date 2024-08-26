package gcr

const (
	NumTracks    = 35
	NumSectors   = 683
	NumTracksMax = NumTracks * 2
)

const BlockSize = 256

// Size of GCR encoded diskData
const (
	SectorSize = 1 + 10 + 9 + 1 + 325 + 8 // SYNC Header Gap SYNC Data Gap (should be 5 SYNC bytes each)
	TrackSize  = SectorSize * 21          // Each track in gcr_data has 21 sectors
	DiskSize   = TrackSize * NumTracks
)

// Number of sectors of each track
var _numSectors = []uint8{
	0,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	19, 19, 19, 19, 19, 19, 19,
	18, 18, 18, 18, 18, 18,
	17, 17, 17, 17, 17,
}

// Sector offset of start of track in .d64 file
var _sectorOffset = []int{
	0,
	0, 21, 42, 63, 84, 105, 126, 147, 168, 189, 210, 231, 252, 273, 294, 315, 336,
	357, 376, 395, 414, 433, 452, 471,
	490, 508, 526, 544, 562, 580,
	598, 615, 632, 649, 666,
}

var _gcrTable = []uint16{
	0x0a, 0x0b, 0x12, 0x13, 0x0e, 0x0f, 0x16, 0x17,
	0x09, 0x19, 0x1a, 0x1b, 0x0d, 0x1d, 0x1e, 0x15,
}
