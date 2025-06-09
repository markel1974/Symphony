package gcr

import (
	"fmt"
	"log"
)

const (
	//zone3Sectors = 21
	//zone2Sectors = 19
	//zone1Sectors = 18
	//zone0Sectors = 17
	//tailGap = 19
	//sectorGap = 15
	//z3 = ((gcrSectorLen + sectorGap) * zone3Sectors) + tailGap
	//z2 = ((gcrSectorLen + sectorGap) * zone2Sectors) + tailGap
	//z1 = ((gcrSectorLen + sectorGap) * zone1Sectors) + tailGap
	//z0 = ((gcrSectorLen + sectorGap) * zone0Sectors) + tailGap
	realTrackSizeBytesZone3 = 7936.0 // Per tracce 1-17 (21 settori)
	realTrackSizeBytesZone2 = 7680.0 // Per tracce 18-24 (19 settori)
	realTrackSizeBytesZone1 = 7168.0 // Per tracce 25-30 (18 settori)
	realTrackSizeBytesZone0 = 6656.0 // Per tracce 31-35 (17 settori)
	rpm                     = 300.0
	systemClock             = 1000000.0 //985248.0.
	rotationTimeCycles      = (60.0 / rpm) * systemClock
	cyclesPerByteZone3      = rotationTimeCycles / realTrackSizeBytesZone3 // Risultato: ~25.2 cicli
	cyclesPerByteZone2      = rotationTimeCycles / realTrackSizeBytesZone2 // Risultato: ~26.0 cicli
	cyclesPerByteZone1      = rotationTimeCycles / realTrackSizeBytesZone1 // Risultato: ~27.9 cicli
	cyclesPerByteZone0      = rotationTimeCycles / realTrackSizeBytesZone0 // Risultato: ~30.0 cicli
)

// _tracks stores a list of TrackData pointers representing track metadata, including sectors, speed zones, and offsets.
var _tracks []*TrackData

// _totalSectors holds the total count of sectors across all tracks, computed during the initialization process.
var _totalSectors uint

// init initializes the track data and updates track properties based on predefined ranges and configurations.
func init() {
	r := systemClock / 985248.0
	_totalSectors = 0
	currentOffset := uint16(0)
	for trackIdx := uint8(0); trackIdx <= 35; trackIdx++ {
		sectors, cyclesPerByte := getTrackInfo(trackIdx)
		cyclesPerByte *= r
		//if cyclesPerByte > 0 {
		//	cyclesPerByte = 20
		//}
		track := NewTrackData(trackIdx, currentOffset, sectors, cyclesPerByte)
		currentOffset += uint16(track.sectors)
		_tracks = append(_tracks, track)
		_totalSectors += uint(track.sectors)
	}
}

func getTrackInfo(trackIdx uint8) (sectors uint8, cyclesPerByte float64) {
	if trackIdx >= 1 && trackIdx <= 17 {
		return 21, cyclesPerByteZone3
	} else if trackIdx >= 18 && trackIdx <= 24 {
		return 19, cyclesPerByteZone2
	} else if trackIdx >= 25 && trackIdx <= 30 {
		return 18, cyclesPerByteZone1
	} else if trackIdx >= 31 && trackIdx <= 35 {
		return 17, cyclesPerByteZone0
	}
	// Aggiungi qui la logica per le tracce 36-40 se vuoi supportarle
	return 0, 0
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
func getMicroSecPerByte(idx uint8) float64 {
	if idx >= uint8(len(_tracks)) {
		return 0
	}
	return _tracks[idx].microSecPerByte
}

// Aggiungi questa funzione helper da qualche parte nel tuo package
func getInterleaveTable(trackIdx uint8, numSectors uint8) (uint8, []uint8) {
	if numSectors <= 0 {
		return 0, nil
	}
	skew := 10 // Interleave di 10 per le tracce dati
	if trackIdx == 18 {
		skew = 3 // Interleave di 3 per la traccia directory
	}
	table := make([]uint8, numSectors)
	for i := range table {
		table[i] = 0xff // Inizializza con un valore non valido
	}
	currentSector := 0

	for idx := 0; idx < int(numSectors); idx++ {
		for table[currentSector] != 0xff {
			currentSector = (currentSector + 1) % int(numSectors)
		}
		table[currentSector] = uint8(idx)
		currentSector = (currentSector + skew) % int(numSectors)
	}

	last := table[len(table)-1]
	//fmt.Println(table)
	return last, table
}

// createGap creates and returns a byte slice of specified size filled with a constant gap byte value.
func createGap(s int) []byte {
	gap := make([]byte, s)
	for i := range gap {
		gap[i] = gapByte
	}
	return gap
}

// rawSector extracts a specific sector from a disk image based on track offset, header length, and sector index.
// Returns a sector-sized buffer and an error if the extraction fails due to bounds or other inconsistencies.
func rawSector(disk []uint8, headerLen uint8, trackOffset uint16, sectorIdx uint8) ([blockBytesLen]uint8, error) {
	var buffer [blockBytesLen]uint8
	rOffset := (int(trackOffset) + int(sectorIdx)) << 8
	begin := rOffset + int(headerLen)
	end := begin + blockBytesLen
	if begin > len(disk) || end > len(disk) {
		log.Printf("invalid start/end: %d - %d", begin, end)
		return buffer, fmt.Errorf("sector index out of range")
	}
	copy(buffer[:], disk[begin:end])
	return buffer, nil
}
