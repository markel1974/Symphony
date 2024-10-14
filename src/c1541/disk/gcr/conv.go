package gcr

//http://www.unusedino.de/ec64/technical/formats/g64.html
//http://www.baltissen.org/newhtm/1541c.htm

const (
	blockBytesLen     = 256
	dataBlockLen      = blockBytesLen + 4 // data block id (0x07) + blockBytes + checksum + 0x00 + 0x00
	syncLen           = 5
	headerLen         = 10 // GCR([ID $08] [Checksum] [Sector Number] [Track Number]) GCR([ID Char #2] [ID Char #1] [$0F] [$0F])
	gapHeaderLen      = 9
	gapInterSectorLen = 8
	gcrDataLen        = (dataBlockLen / 4) * 5
	gcrSectorLen      = syncLen + headerLen + gapHeaderLen + syncLen + gcrDataLen + gapInterSectorLen
)

const (
	syncMarker       = 0xff
	gapByte    uint8 = 0x55
)

// _gcrTable is a lookup table used for GCR (Group Code Recording) encoding, mapping 4-bit values to their 5-bit GCR encoded values.
var _gcrTable = []uint16{
	0x0a, 0x0b, 0x12, 0x13, 0x0e, 0x0f, 0x16, 0x17,
	0x09, 0x19, 0x1a, 0x1b, 0x0d, 0x1d, 0x1e, 0x15,
}

// _gcrFromTable is a lookup table used for converting 5-bit GCR encoded values back to their original 4-bit values.
var _gcrFromTable = []uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 1, 0, 12, 4, 5, 0,
	0, 2, 3, 0, 15, 6, 7, 0, 9, 10, 11, 0, 13, 14, 0,
}

var _syncBlock [syncLen]uint8
var _gapHeaderBlock [gapHeaderLen]uint8
var _gapInterSectorBlock [gapInterSectorLen]uint8

func init() {
	for i := range _syncBlock {
		_syncBlock[i] = syncMarker
	}
	for i := range _gapHeaderBlock {
		_gapHeaderBlock[i] = gapByte
	}
	for i := range _gapInterSectorBlock {
		_gapInterSectorBlock[i] = gapByte
	}
}

func toGCR(from uint8) uint16 {
	to := (_gcrTable[(from>>4)&0xf] << 5) | _gcrTable[from&0xf]
	return to
}

func fromGCR(from uint32) uint8 {
	g := _gcrFromTable[(from>>16)&0x1f]
	return g
}

// conv4to5 converts 4 bytes to 5 GCR encoded bytes
func conv4to5(from [4]uint8) [5]uint8 {
	var to [5]uint8
	g := toGCR(from[0])
	to[0] = uint8(g >> 2)
	to[1] = uint8((g << 6) & 0xc0)
	g = toGCR(from[1])
	to[1] |= uint8((g >> 4) & 0x3f)
	to[2] = uint8((g << 4) & 0xf0)
	g = toGCR(from[2])
	to[2] |= uint8((g >> 6) & 0x0f)
	to[3] = uint8((g << 2) & 0xfc)
	g = toGCR(from[3])
	to[3] |= uint8((g >> 8) & 0x03)
	to[4] = uint8(g)
	return to
}

// conv5to4 converts GCR encoded bytes into an array of 4 decoded bytes
func conv5to4(src [5]uint8) [4]uint8 {
	var dst [4]uint8
	t := uint32(src[0])
	t <<= 13
	sourceIdx := 1
	for dstIdx, shift := range []uint8{5, 7, 9, 11} {
		t |= (uint32(src[sourceIdx])) << shift
		dst[dstIdx] = fromGCR(t) << 4 //_gcrFromTable[(t>>16)&0x1f] << 4
		t <<= 5
		dst[dstIdx] |= fromGCR(t) //_gcrFromTable[(t>>16)&0x1f]
		t <<= 5
		sourceIdx++
	}
	return dst
}

func sector2gcr(sector [blockBytesLen]uint8, id1 uint8, id2 uint8, trackIdx uint8, sectorIdx uint8) [gcrSectorLen]uint8 {
	const last = blockBytesLen - 1

	var ret [gcrSectorLen]uint8
	idx := 0
	copy(ret[idx:], _syncBlock[:])
	idx += len(_syncBlock)

	headerData := conv4to5([4]uint8{0x08, uint8(int(sectorIdx) ^ int(trackIdx) ^ int(id2) ^ int(id1)), sectorIdx, trackIdx})
	copy(ret[idx:], headerData[:])
	idx += len(headerData)

	bamData := conv4to5([4]uint8{id2, id1, 0x0f, 0x0f})
	copy(ret[idx:], bamData[:])
	idx += len(bamData)

	copy(ret[idx:], _gapHeaderBlock[:])
	idx += len(_gapHeaderBlock)

	copy(ret[idx:], _syncBlock[:])
	idx += len(_syncBlock)

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

	copy(ret[idx:], _gapInterSectorBlock[:])
	idx += len(_gapInterSectorBlock)

	return ret
}

/*
//Users/tinmr305/Desktop/emu/vice-emu-code-r45201-trunk-vice/src/gcr.c
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
