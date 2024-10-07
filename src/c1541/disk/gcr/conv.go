package gcr

import (
	"fmt"
	"log"
)

//http://www.unusedino.de/ec64/technical/formats/g64.html

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

func rawSector(disk []uint8, headerLen uint8, trackOffset uint16, sectorIdx uint8) ([blockSize]uint8, error) {
	var buffer [blockSize]uint8
	rOffset := (int(trackOffset) + int(sectorIdx)) << 8
	begin := rOffset + int(headerLen)
	end := begin + blockSize
	if begin > len(disk) || end > len(disk) {
		log.Printf("invalid start/end: %d - %d", begin, end)
		return buffer, fmt.Errorf("sector index out of range")
	}
	copy(buffer[:], disk[begin:end])
	return buffer, nil
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

func sector2gcr(sector [blockSize]uint8, bam1 uint8, bam2 uint8, trackIdx uint8, sectorIdx uint8) [sectorSize]uint8 {
	const last = blockSize - 1
	const fillByte uint8 = 0x55
	const syncMarker = 0xff

	var ret [sectorSize]uint8
	idx := 0
	ret[idx] = syncMarker
	idx++

	headerData := conv4to5([4]uint8{0x08, uint8(int(sectorIdx) ^ int(trackIdx) ^ int(bam2) ^ int(bam1)), sectorIdx, trackIdx})
	copy(ret[idx:], headerData[:])
	idx += len(headerData)

	bamData := conv4to5([4]uint8{bam2, bam1, 0x0f, 0x0f})
	copy(ret[idx:], bamData[:])
	idx += len(bamData)

	fillData := [9]uint8{fillByte, fillByte, fillByte, fillByte, fillByte, fillByte, fillByte, fillByte, fillByte}
	copy(ret[idx:], fillData[:])
	idx += len(fillData)

	ret[idx] = syncMarker
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

	endData := [8]uint8{fillByte, fillByte, fillByte, fillByte, fillByte, fillByte, fillByte, fillByte}
	copy(ret[idx:], endData[:])
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
