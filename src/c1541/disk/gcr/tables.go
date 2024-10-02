package gcr

//http://www.unusedino.de/ec64/technical/formats/g64.html

// _numSectors  Number of sectors of each track
var _numSectors = []uint8{
	0,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	19, 19, 19, 19, 19, 19, 19,
	18, 18, 18, 18, 18, 18,
	17, 17, 17, 17, 17,
}

// _sectorOffset Sector offset of start of track in .d64 file
var _sectorOffset = []int{
	0,
	0, 21, 42, 63, 84, 105, 126, 147, 168, 189, 210, 231, 252, 273, 294, 315, 336,
	357, 376, 395, 414, 433, 452, 471,
	490, 508, 526, 544, 562, 580,
	598, 615, 632, 649, 666,
}

// _gcrTable is a lookup table used for GCR (Group Code Recording) encoding, mapping 4-bit values to their 5-bit GCR encoded values.
var _gcrTable = []uint16{
	0x0a, 0x0b, 0x12, 0x13, 0x0e, 0x0f, 0x16, 0x17,
	0x09, 0x19, 0x1a, 0x1b, 0x0d, 0x1d, 0x1e, 0x15,
}

var _gcrFromTable = []uint8{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 8, 0, 1, 0, 12, 4, 5,
	0, 0, 2, 3, 0, 15, 6, 7,
	0, 9, 10, 11, 0, 13, 14, 0,
}

//func secNumFromTs(track int, sector int) int {
//	return _sectorOffset[track] + sector
//}
