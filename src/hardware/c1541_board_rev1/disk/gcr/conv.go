package gcr

//http://www.unusedino.de/ec64/technical/formats/g64.html
//http://www.baltissen.org/newhtm/1541c.htm

// blockBytesLen defines the length of a block in bytes.
// dataBlockLen represents the total length of a data block including id, bytes, checksum, and padding.
// syncLen specifies the length of the synchronization sequence.
// headerLen specifies the length of the header for a single sector, including ID, checksum, and other metadata.
// gapHeaderLen defines the length of the gap following the header section.
// gapInterSectorLen defines the length of the gap between two sectors.
// gcrDataLen specifies the encoded data length for GCR based on the data block length.
// gcrSectorLen represents the total length of a GCR-encoded sector, including sync, header, gaps, and data.
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

// syncMarker is a constant used as a synchronization marker with a value of 0xff.
// gapByte is a uint8 constant with a value of 0x55, potentially used as a padding or gap indicator.
const (
	syncMarker       = 0xff
	gapByte    uint8 = 0x55
)

// _gcrTable is a lookup table that maps 4-bit input values to their corresponding 5-bit GCR-encoded values.
var _gcrTable = []uint16{
	0x0a, 0x0b, 0x12, 0x13, 0x0e, 0x0f, 0x16, 0x17,
	0x09, 0x19, 0x1a, 0x1b, 0x0d, 0x1d, 0x1e, 0x15,
}

// _gcrFromTable is a lookup table mapping 5-bit GCR values to their corresponding 4-bit decoded values.
var _gcrFromTable = []uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 1, 0, 12, 4, 5, 0,
	0, 2, 3, 0, 15, 6, 7, 0, 9, 10, 11, 0, 13, 14, 0,
}

// _syncBlock is an array representing a synchronization block, initialized with `syncMarker`, used in data encoding processes.
var _syncBlock [syncLen]uint8

// _gapHeaderBlock is a predefined array filled with `gapByte` values, used to create spacing within data blocks.
var _gapHeaderBlock [gapHeaderLen]uint8

// _gapInterSectorBlock is a predefined array used to define a gap between sectors in a GCR-encoded data block.
var _gapInterSectorBlock [gapInterSectorLen]uint8

// init initializes the synchronization blocks and gap blocks with predefined marker and gap byte values respectively.
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

// toGCR converts a 4-bit data value into its corresponding 5-bit GCR encoding using a predefined lookup table.
func toGCR(from uint8) uint16 {
	to := (_gcrTable[(from>>4)&0xf] << 5) | _gcrTable[from&0xf]
	return to
}

// fromGCR maps a 5-bit encoded value from `from` to an 8-bit value using the `_gcrFromTable`.
func fromGCR(from uint32) uint8 {
	g := _gcrFromTable[(from>>16)&0x1f]
	return g
}

// conv4to5 converts a 4-byte input array into a 5-byte output array using GCR transformation.
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

// conv5to4 converts a 5-byte input array into a 4-byte output array using GCR decoding and bit-shifting operations.
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

// sector2gcr converts a disk sector into its GCR-encoded representation including header, data, checksum, and padding.
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
// BuildTrackImage costruisce l'immagine binaria completa di una singola traccia.
// Prende come input l'indice della traccia, una mappa dei dati dei settori e l'ID del disco.
// Restituisce una slice di byte che rappresenta la traccia GCR completa, pronta per essere letta dal mechanic.
func BuildTrackImage(trackIdx uint8, sectorsData map[uint8][blockBytesLen]uint8, id1, id2 uint8) []byte {
	numSectors, usPerByte := getTrackInfo(trackIdx)
	if numSectors == 0 {
		return nil // Traccia non valida o vuota
	}

	// --- PASSO 1: Calcola la dimensione totale corretta della traccia in byte ---
	// Questa è la correzione fondamentale al tuo approccio statico.
	totalTrackBytes := int(rotationTimeUs / usPerByte)

	// Usiamo un buffer per costruire dinamicamente la nostra traccia.
	var trackBuffer bytes.Buffer
	trackBuffer.Grow(totalTrackBytes)

	// --- PASSO 2: Scrivi tutti i settori formattati nel buffer ---
	for sectorIdx := uint8(0); sectorIdx < numSectors; sectorIdx++ {
		sectorData, ok := sectorsData[sectorIdx]
		if !ok {
			// Se mancano i dati per un settore, lo creiamo vuoto.
			sectorData = [blockBytesLen]uint8{}
		}

		// Usiamo la tua funzione esistente per creare il blocco GCR del settore.
		gcrSectorBlock := sector2gcr(sectorData, id1, id2, trackIdx, sectorIdx)
		trackBuffer.Write(gcrSectorBlock[:])
	}

	// --- PASSO 3: Calcola e aggiungi il "Tail Gap" finale ---
	// Questo è lo spazio vuoto alla fine della traccia per completare i 200,000µs.
	currentLength := trackBuffer.Len()
	tailGapSize := totalTrackBytes - currentLength

	if tailGapSize > 0 {
		tailGap := make([]byte, tailGapSize)
		for i := range tailGap {
			tailGap[i] = gapByte // Riempiamo il gap con 0x55
		}
		trackBuffer.Write(tailGap)
	} else if tailGapSize < 0 {
		// Questo è un segnale di errore: i settori occupano più spazio di quello disponibile sulla traccia.
		fmt.Printf("ATTENZIONE: La Traccia %d è troppo lunga di %d byte!\n", trackIdx, -tailGapSize)
	}

	return trackBuffer.Bytes()
}
*/

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
