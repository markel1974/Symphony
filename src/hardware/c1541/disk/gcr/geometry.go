package gcr

import (
	"fmt"
	"log"
)

//Zone	Track	Sectors	Data(362*N) Tail Gap 	Capacity	Value Timer VIA
//3   	1-17	21  	7602 byte	90  byte	7692 byte	32 / $20
//2   	18-24	19   	6878 byte	262 byte	7140 byte	30 / $1E
//1     25-30	18   	6516 byte	150 byte	6666 byte	28 / $1C
//0   	31-35	17  	6154 byte	93  byte	6247 byte	26 / $1A

const (
	zone3Sectors = 21
	zone2Sectors = 19
	zone1Sectors = 18
	zone0Sectors = 17

	zone3Gap = 90
	zone2Gap = 262
	zone1Gap = 150
	zone0Gap = 93

	rpm                = 300.0
	systemClockReal    = 1000000.0 //  985248.0 1_022_727.0
	systemClockFactor  = 1.0       //systemClockReal / 985248.0
	systemClock        = systemClockReal * systemClockFactor
	rotationTimeCycles = (60.0 / rpm) * systemClock

	realTrackSizeBytesZone3 = (gcrSectorLen * zone3Sectors) + zone3Gap // 7692 // Per tracce 1-17 (21 settori * 362)
	realTrackSizeBytesZone2 = (gcrSectorLen * zone2Sectors) + zone2Gap // 7140 // Per tracce 18-24 (19 settori)
	realTrackSizeBytesZone1 = (gcrSectorLen * zone1Sectors) + zone1Gap // 6666 // Per tracce 25-30 (18 settori)
	realTrackSizeBytesZone0 = (gcrSectorLen * zone0Sectors) + zone0Gap // 6247 // Per tracce 31-35 (17 settori)

	cyclesPerByteZone3 = rotationTimeCycles / realTrackSizeBytesZone3
	cyclesPerByteZone2 = rotationTimeCycles / realTrackSizeBytesZone2
	cyclesPerByteZone1 = rotationTimeCycles / realTrackSizeBytesZone1
	cyclesPerByteZone0 = rotationTimeCycles / realTrackSizeBytesZone0
)

// _tracks stores a list of TrackData pointers representing track metadata, including sectors, speed zones, and offsets.
var _tracks []*TrackData

// _totalSectors holds the total count of sectors across all tracks, computed during the initialization process.
var _totalSectors uint

// init initializes the track data and updates track properties based on predefined ranges and configurations.
func init() {
	_totalSectors = 0
	currentOffset := uint16(0)
	for trackIdx := uint8(0); trackIdx <= 35; trackIdx++ {
		sectors, cyclesPerByte, _ := getTrackInfo(trackIdx)
		track := NewTrackData(trackIdx, currentOffset, sectors, cyclesPerByte)
		currentOffset += uint16(track.sectors)
		_tracks = append(_tracks, track)
		_totalSectors += uint(track.sectors)
	}
}

func getTrackInfo(trackIdx uint8) (sectors uint8, cyclesPerByte float64, size int) {
	if trackIdx >= 1 && trackIdx <= 17 {
		return zone3Sectors, cyclesPerByteZone3, realTrackSizeBytesZone3
	} else if trackIdx >= 18 && trackIdx <= 24 {
		return zone2Sectors, cyclesPerByteZone2, realTrackSizeBytesZone2
	} else if trackIdx >= 25 && trackIdx <= 30 {
		return zone1Sectors, cyclesPerByteZone1, realTrackSizeBytesZone1
	} else if trackIdx >= 31 && trackIdx <= 35 {
		return zone0Sectors, cyclesPerByteZone0, realTrackSizeBytesZone0
	}
	return 0, 0, 0
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

// getInterleaveTable calculates and returns the interleave table for a track and the last sector index in the table.
// trackIdx specifies the track index for which the interleave table is generated.
// numSectors defines the total number of sectors in the track.
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
