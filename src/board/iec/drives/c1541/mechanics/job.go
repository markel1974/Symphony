package mechanics

import (
	"github.com/markel1974/c64emu/src/board/iec/drives/c1541/ram"
	"github.com/markel1974/c64emu/src/preferences"
	"io"
	"os"
)

type Job struct {
	*Core

	_image_header int // Length of .d64/.x64 file header

	_id1 uint8 // ID of disk

	_id2 uint8 // ID of disk

	_error_info []uint8 // Sector error information (1 byte/sector)

	_gcr_data []uint8 // Pointer to GCR encoded disk data

	_gcr_ptr int // Pointer to GCR data under R/W head

	_gcr_track_start int // Pointer to start of GCR data of current track

	_gcr_track_end int // Pointer to end of GCR data of current track

	_data []uint8

	_filePath string

	ram *ram.Ram

	deviceNumber uint8
}

func NewJob(ram *ram.Ram, deviceNumber uint8) *Job {
	j := &Job{
		ram:          ram,
		Core:         NewCore(),
		deviceNumber: deviceNumber,
	}
	j._id1 = 0
	j._id2 = 0
	j._image_header = 0
	j._gcr_data = make([]uint8, GCR_DISK_SIZE)
	j._gcr_ptr = 0
	j._gcr_track_start = 0
	j._gcr_track_end = j._gcr_track_start + GCR_TRACK_SIZE
	j._current_halftrack = 2
	j._error_info = make([]uint8, NUM_SECTORS)
	return j
}

func (j *Job) Setup(prefs *preferences.Prefs) {
	//if (prefs.Emul1541Proc()) {
	//	_filePath = prefs.GetDrivePath(board->GetDeviceNumber() - 8)
	//	openFile(_filePath)
	//}
	//if !prefs.Emul1541Proc() {
	//	j.closeFile()
	//	return
	//}
	filePath := prefs.GetDrivePath(int(j.deviceNumber - 8))
	if !j.HasDisk() {
		j._filePath = filePath
		j.openFile()
	} else if j._filePath != filePath {
		j._filePath = filePath
		j.closeFile()
		j.openFile()
		j._disk_changed = true
	}
}

func (j *Job) HasDisk() bool {
	return j._data != nil
}

func (j *Job) WriteProtectionState() uint8 {
	r := uint8(0)
	if j._disk_changed {
		// Disk change -> WP sensor strobe
		j._disk_changed = false
		if j._write_protected {
			r = 0x10
		}
	} else {
		if !j._write_protected {
			r = 0x10
		}
	}
	return r
}

func (j *Job) SyncFound() bool {
	if j._gcr_data[j._gcr_ptr] == 0xff {
		return true
	}
	// Rotate disk
	j._gcr_ptr++
	if j._gcr_ptr == j._gcr_track_end {
		j._gcr_ptr = j._gcr_track_start
	}
	return false

}

func (j *Job) WriteSector() {
	track := j.ram.Read(0x18)
	sector := j.ram.Read(0x19)
	start := uint16(j.ram.Read(0x30)) | (uint16(j.ram.Read(0x31)) << 8)
	if start <= 0x0700 {
		block := j.ram.Interval(start, BLOCK_SIZE)
		if j.writeTrackSector(int(track), int(sector), block) {
			j.sector2gcr(int(track), int(sector))
		}
	}
}

func (j *Job) FormatTrack() {
	track := j.ram.Read(0x51)
	// Get new ID
	bufNum := j.ram.Read(0x3d)
	j._id1 = j.ram.Read(0x12 + uint16(bufNum))
	j._id2 = j.ram.Read(0x13 + uint16(bufNum))

	// Create empty block
	buf := make([]uint8, BLOCK_SIZE)
	buf[0] = 0x4b

	// Write block to all sectors on track
	for sector := 0; sector < int(__num_sectors[track]); sector++ {
		j.writeTrackSector(int(track), sector, buf)
		j.sector2gcr(int(track), sector)
	}

	// Clear error info (all sectors no error)
	if track == 35 {
		for x := range j._error_info {
			j._error_info[x] = 1
		}
		// Write error_info to disk?
	}
}

func (j *Job) ReadGCRByte() uint8 {
	data := j._gcr_data[j._gcr_ptr]
	j._gcr_ptr++ // Rotate disk
	if j._gcr_ptr == j._gcr_track_end {
		j._gcr_ptr = j._gcr_track_start
	}
	return data
}

func (j *Job) UpdateLEDs(l int) {
	//TODO
	// LedStateChangedEvent.Emit(_board->GetDeviceNumber(), state);
}

func (j *Job) MoveHeadOut() {
	if j._current_halftrack == 2 {
		//NOTHING TO DO
	} else {
		j._current_halftrack--
		idx := ((j._current_halftrack >> 1) - 1) * GCR_TRACK_SIZE
		//data := j._gcr_data[idx]
		j._gcr_track_start = idx //data
		j._gcr_ptr = idx         //data
		trackLength := int(__num_sectors[j._current_halftrack>>1]) * GCR_SECTOR_SIZE
		j._gcr_track_end = j._gcr_track_start + trackLength
	}

	//TODO
	// HeadPosChangedEvent.Emit(_board->GetDeviceNumber(), _current_halftrack);
}

func (j *Job) MoveHeadIn() {
	if j._current_halftrack == NUM_TRACKS*2 {
		//NOTHING TO DO
	} else {
		j._current_halftrack++
		idx := ((j._current_halftrack >> 1) - 1) * GCR_TRACK_SIZE
		j._gcr_track_start = idx
		j._gcr_ptr = idx
		trackLength := int(__num_sectors[j._current_halftrack>>1]) * GCR_SECTOR_SIZE
		j._gcr_track_end = j._gcr_track_start + trackLength
	}
	//TODO
	// HeadPosChangedEvent.Emit(_board->GetDeviceNumber(), _current_halftrack);
}

func (j *Job) openFile() bool {
	for x := range j._gcr_data {
		j._gcr_data[x] = 0x55
	}
	j._write_protected = false
	fd, err := os.OpenFile(j._filePath, os.O_RDWR, 0)
	if err != nil {
		j._write_protected = true
		fd, err = os.OpenFile(j._filePath, os.O_RDONLY, 0)
	}
	if err != nil {
		return false
	}
	defer fd.Close()
	j._data, err = io.ReadAll(fd)
	if err != nil {
		return false
	}
	size := len(j._data)
	if size < NUM_SECTORS*BLOCK_SIZE {
		return false
	}
	if j._data[0] == 0x43 && j._data[1] == 0x15 && j._data[2] == 0x41 && j._data[3] == 0x64 {
		j._image_header = 64
	} else {
		j._image_header = 0
	}
	for x := range j._error_info {
		j._error_info[x] = 1
	}
	// Load sector error info from .d64 file, if present
	if j._image_header == 0 && size == NUM_SECTORS*257 {
		copy(j._error_info, j._data[NUM_SECTORS*BLOCK_SIZE:])
	}
	// Read BAM and get ID
	bam := j.readSector(18, 0)
	if bam == nil {
		return false
	}

	j._id1 = bam[162]
	j._id2 = bam[163]

	// Create GCR encoded disk data from image
	j.disk2gcr()
	return true
}

func (j *Job) closeFile() {
	j._data = nil
	//TODO
}

func (j *Job) readSector(track int, sector int) []uint8 {
	// Convert track/sector to byte offset in file
	offset := j.offsetFromTrackSector(track, sector)
	if offset < 0 {
		return nil
	}
	start := offset + j._image_header
	if end := start + BLOCK_SIZE; end >= len(j._data) {
		return nil
	}
	buffer := make([]uint8, BLOCK_SIZE)
	copy(buffer, j._data[start:])
	return buffer
}

func (j *Job) writeTrackSector(track int, sector int, buffer []uint8) bool {
	offset := j.offsetFromTrackSector(track, sector)
	// Convert track/sector to byte offset in file
	if offset < 0 {
		return false
	}
	copy(j._data[offset+j._image_header:], buffer)
	_ = os.WriteFile("a", j._data, 0644)
	return true
}

/*
 *  Convert track/sector to offset
 */

func (j *Job) secNumFromTs(track int, sector int) int {
	return __sector_offset[track] + sector
}

func (j *Job) offsetFromTrackSector(track int, sector int) int {
	if (track < 1) || (track > NUM_TRACKS) || (sector < 0) || (sector >= int(__num_sectors[track])) {
		return -1
	}
	return (__sector_offset[track] + sector) << 8
}

/*
 *  Convert 4 bytes to 5 GCR encoded bytes
 */

func (j *Job) gcrConv4(from []uint8, to []uint8) {
	g := (__gcr_table[from[0]>>4] << 5) | __gcr_table[from[0]&15]
	to[0] = uint8(g >> 2)
	to[1] = uint8((g << 6) & 0xc0)
	g = (__gcr_table[from[1]>>4] << 5) | __gcr_table[from[1]&15]
	to[1] |= uint8((g >> 4) & 0x3f)
	to[2] = uint8((g << 4) & 0xf0)
	g = (__gcr_table[from[2]>>4] << 5) | __gcr_table[from[2]&15]
	to[2] |= uint8((g >> 6) & 0x0f)
	to[3] = uint8((g << 2) & 0xfc)
	g = (__gcr_table[from[3]>>4] << 5) | __gcr_table[from[3]&15]
	to[3] |= uint8((g >> 8) & 0x03)
	to[4] = uint8(g)
}

func (j *Job) sector2gcr(track int, sector int) {
	buf := make([]uint8, 4)
	p := (track-1)*GCR_TRACK_SIZE + sector*GCR_SECTOR_SIZE
	block := j.readSector(track, sector)
	if block == nil {
		return
	}
	// Create GCR header
	// SYNC
	j._gcr_data[p] = 0xff
	p++
	// Header mark
	buf[0] = 0x08
	// Checksum
	buf[1] = uint8(sector ^ track ^ int(j._id2) ^ int(j._id1))
	buf[2] = uint8(sector)
	buf[3] = uint8(track)
	j.gcrConv4(buf, j._gcr_data[p:])
	buf[0] = j._id2
	buf[1] = j._id1
	buf[2] = 0x0f
	buf[3] = 0x0f
	j.gcrConv4(buf, j._gcr_data[p+5:])
	p += 10
	for x := 0; x < 9; x++ {
		j._gcr_data[p+x] = 0x55
	}
	p += 9
	// Create GCR data
	// SYNC
	j._gcr_data[p] = 0xff
	p++
	// Data mark
	buf[0] = 0x07
	buf[1] = block[0]
	sum := buf[1]
	buf[2] = block[1]
	sum ^= buf[2]
	buf[3] = block[2]
	sum ^= buf[3]
	j.gcrConv4(buf, j._gcr_data[p:])
	p += 5
	for i := 3; i < 255; i += 4 {
		buf[0] = block[i]
		sum ^= buf[0]
		buf[1] = block[i+1]
		sum ^= buf[1]
		buf[2] = block[i+2]
		sum ^= buf[2]
		buf[3] = block[i+3]
		sum ^= buf[3]
		j.gcrConv4(buf, j._gcr_data[p:])
		p += 5
	}
	buf[0] = block[255]
	sum ^= buf[0]
	// Checksum
	buf[1] = sum
	buf[2] = 0
	buf[3] = 0
	j.gcrConv4(buf, j._gcr_data[p:])
	p += 5
	for x := 0; x < 8; x++ {
		j._gcr_data[p+x] = 0x55
	}
}

func (j *Job) disk2gcr() {
	// Convert all tracks and sectors
	for track := 1; track <= NUM_TRACKS; track++ {
		for sector := 0; sector < int(__num_sectors[track]); sector++ {
			j.sector2gcr(track, sector)
		}
	}
}
