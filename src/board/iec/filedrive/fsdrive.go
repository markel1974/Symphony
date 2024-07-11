package filedrive

import "os"

type DriveHelper struct {
	deviceNumber   uint8
	_dir_path      string       // Path to directory
	_orig_dir_path string       // Original directory path
	_dir_title     string       // Directory title
	_file          [16]*os.File // File pointers for each of the 16 channels
	_read_char     [16]uint8    // Buffers for one-byte read-ahead
}

func NewDriveHelper(deviceNumber uint8, path string) *DriveHelper {
	return &DriveHelper{deviceNumber: deviceNumber, _dir_path: path}
}

func (v *DriveHelper) Reset() {
	//c.closeAllChannels()
	//c.cmd_len = 0
	//setError(ERR_STARTUP)
}

func (v *DriveHelper) Ready() bool {
	return false
}

func (v *DriveHelper) GetPath() string {
	return v._dir_path
}

func (v *DriveHelper) Open(uint8, []byte) uint8 {
	return 0
}

func (v *DriveHelper) Close(uint8) uint8 {
	return 0
}

func (v *DriveHelper) Read(channel uint8, data *uint8) uint8 {
	return 0

	/*
		var c int
		// Channel 15: Error channel
		if channel == 15 {
			*data = *c.error_ptr
			c.error_ptr++
			if *data != '\r' {
				return ST_OK
			} else { // End of message
				c.setError(ERR_OK)
				return ST_EOF
			}
		}

		if !c._file[channel] {
			return ST_READ_TIMEOUT
		}

		// Read one byte
		*data = c._read_char[channel]
		c = fgetc(_file[channel])
		if c == EOF {
			return ST_EOF
		} else {
			c._read_char[channel] = (uint8)(c)
			return ST_OK
		}
	*/
}

func (v *DriveHelper) Write(channel uint8, data uint8, eoi bool) uint8 {
	return 0
	// Channel 15: Collect chars and execute command on EOI
	/*
		if channel == 15 {
			if c.cmd_len >= 58 {
				return ST_TIMEOUT
			}
			c.cmd_buf[c.cmd_len] = data
			c.cmd_len++
			if eoi {
				c.executeCmd(c.cmd_buf, c.cmd_len)
				c.cmd_len = 0
			}
			return ST_OK
		}
		if !v_file[channel] {
			c.setError(ERR_FILENOTOPEN)
			return ST_TIMEOUT
		}
		if v.putc(data, v._file[channel]) == EOF {
			setError(ERR_WRITE25)
			return ST_TIMEOUT
		}
		return ST_OK
	*/
}

func (v *DriveHelper) initializeCmd() {
	//v.closeAllChannels()
}

func (v *DriveHelper) validateCmd() {
}

func CreateDrive(deviceNumber uint8, path string) *DriveHelper {
	v := NewDriveHelper(deviceNumber, path)
	return v
}
