package gcr

//http://www.unusedino.de/ec64/technical/formats/g64.html

// _numSectors  Number of sectors of each track
var _numSectors [0xff]uint8

// _sectorOffset Sector offset of start of track in .d64 file
var _sectorOffset [0xff]uint16

// _gcrTable is a lookup table used for GCR (Group Code Recording) encoding, mapping 4-bit values to their 5-bit GCR encoded values.
var _gcrTable [0xff]uint16

// _gcrFromTable is a lookup table used for converting 5-bit GCR encoded values back to their original 4-bit values.
var _gcrFromTable [0xff]uint8

func init() {
	var sectorOffset = []uint16{
		0, 0, 21, 42, 63, 84, 105, 126, 147, 168, 189, 210, 231, 252, 273, 294, 315, 336,
		357, 376, 395, 414, 433, 452, 471, 490, 508, 526, 544, 562, 580, 598, 615, 632, 649, 666,
	}
	var numSec = []uint8{
		0, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
		19, 19, 19, 19, 19, 19, 19, 18, 18, 18, 18, 18, 18, 17, 17, 17, 17, 17,
	}
	var gcrTable = []uint16{
		0x0a, 0x0b, 0x12, 0x13, 0x0e, 0x0f, 0x16, 0x17,
		0x09, 0x19, 0x1a, 0x1b, 0x0d, 0x1d, 0x1e, 0x15,
	}
	var gcrFromTable = []uint8{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 1, 0, 12, 4, 5, 0,
		0, 2, 3, 0, 15, 6, 7, 0, 9, 10, 11, 0, 13, 14, 0,
	}
	for idx, v := range gcrTable {
		_gcrTable[idx] = v
	}
	for idx, v := range gcrFromTable {
		_gcrFromTable[idx] = v
	}
	for idx, v := range sectorOffset {
		_sectorOffset[idx] = v
	}
	for idx, v := range numSec {
		_numSectors[idx] = v
	}
}

// conv5to4 converts GCR encoded bytes into an array of 4 decoded bytes
func conv5to4(src [5]uint8) [4]uint8 {
	var dst [4]uint8
	t := uint32(src[0])
	t <<= 13
	sourceIdx := 1
	for dstIdx, shift := range []uint8{5, 7, 9, 11} {
		t |= (uint32(src[sourceIdx])) << shift
		dst[dstIdx] = _gcrFromTable[(t>>16)&0x1f] << 4
		t <<= 5
		dst[dstIdx] |= _gcrFromTable[(t>>16)&0x1f]
		t <<= 5
		sourceIdx++
	}
	return dst
}

// conv4to5 converts 4 bytes to 5 GCR encoded bytes
func conv4to5(from [4]uint8) [5]uint8 {
	var to [5]uint8
	g := (_gcrTable[(from[0]>>4)&0xf] << 5) | _gcrTable[from[0]&0xf]
	to[0] = uint8(g >> 2)
	to[1] = uint8((g << 6) & 0xc0)
	g = (_gcrTable[(from[1]>>4)&0xf] << 5) | _gcrTable[from[1]&0xf]
	to[1] |= uint8((g >> 4) & 0x3f)
	to[2] = uint8((g << 4) & 0xf0)
	g = (_gcrTable[(from[2]>>4)&0xf] << 5) | _gcrTable[from[2]&0xf]
	to[2] |= uint8((g >> 6) & 0x0f)
	to[3] = uint8((g << 2) & 0xfc)
	g = (_gcrTable[(from[3]>>4)&0xf] << 5) | _gcrTable[from[3]&0xf]
	to[3] |= uint8((g >> 8) & 0x03)
	to[4] = uint8(g)
	return to
}

//func secNumFromTs(track int, sector int) int {
//	return _sectorOffset[track] + sector
//}
