package drives

import (
	"github.com/markel1974/c64emu/src/board/iec/virtualdrive"
	"os"
	"strings"
)

type DriveHelper struct {
	*virtualdrive.VirtualDrive
	_dir_path      string       // Path to directory
	_orig_dir_path string       // Original directory path
	_dir_title     string       // Directory title
	_file          [16]*os.File // File pointers for each of the 16 channels
	_read_char     [16]uint8    // Buffers for one-byte read-ahead
	_ready         bool

	cmd_len int
}

func NewDriveHelper(deviceNumber uint8, path string) *DriveHelper {
	d := &DriveHelper{
		VirtualDrive:   virtualdrive.NewVirtualDrive(deviceNumber),
		_orig_dir_path: path,
	}
	if d.changeDirectory(d._orig_dir_path) {
		for i := 0; i < 16; i++ {
			d._file[i] = nil
		}
		d.Reset()
		d._ready = true
	}
	return d
}

func (v *DriveHelper) Reset() {
	v.closeAllChannels()
	v.cmd_len = 0
	v.SetError(virtualdrive.ERR_STARTUP)
}

func (v *DriveHelper) Ready() bool {
	return false
}

func (v *DriveHelper) GetPath() string {
	return v._dir_path
}

func (v *DriveHelper) Open(channel uint8, name []byte) uint8 {
	//TODO IMPLEMENT
	return 0
}

/*
func (v *DriveHelper) Open(channel uint8, name string) uint8 {
	v.LedTurnOn()
	// Channel 15: Execute file name as command
	if channel == 15 {
		v.executeCmd(name)
		return virtualdrive.StOk
	}
	// Close previous file if still open
	if v._file[channel] != nil {
		v._file[channel].Close()
		v._file[channel] = nil
	}
	if name[0] == '#' {
		v.setError(virtualdrive.ERR_NOCHANNEL)
		return virtualdrive.StOk
	}
	if name[0] == '$' {
		return v.openDirectory(channel, name)
	}
	return v.openFile(channel, name)
}
*/

func (v *DriveHelper) openFile(channel uint8, name string) uint8 {
	var plainName string
	//var plainNameLen int
	mode := virtualdrive.FMODE_READ
	kind := virtualdrive.FTYPE_PRG
	//recLen := 0
	//v.parseFileName(name, plainName, plainNameLen, mode, kind, recLen, true)
	// Channel 0 is READ, channel 1 is WRITE
	if channel == 0 || channel == 1 {
		mode = virtualdrive.FMODE_READ
		if channel != 0 {
			mode = virtualdrive.FMODE_WRITE
		}
		if kind == virtualdrive.FTYPE_DEL {
			kind = virtualdrive.FTYPE_PRG
		}
	}
	writing := mode == virtualdrive.FMODE_WRITE || mode == virtualdrive.FMODE_APPEND
	if strings.Contains(plainName, "*") || strings.Contains(plainName, "?") {
		if writing {
			v.SetError(virtualdrive.ERR_SYNTAX33)
			return virtualdrive.StOk
		} else {
			v.findFirstFile(plainName)
		}
	}
	if kind == virtualdrive.FTYPE_REL {
		v.SetError(virtualdrive.ERR_UNIMPLEMENTED)
		return virtualdrive.StOk
	}
	flags := os.O_RDONLY
	perm := os.FileMode(0)
	switch mode {
	case virtualdrive.FMODE_WRITE:
		perm = os.FileMode(0666)
		flags = os.O_RDWR
	case virtualdrive.FMODE_APPEND:
		perm = os.FileMode(0666)
		flags = os.O_RDWR | os.O_APPEND
	}
	completeFileName := v._dir_path + string(os.PathSeparator) + plainName
	f, err := os.OpenFile(completeFileName, flags, perm)
	if err != nil {
		v.SetError(virtualdrive.ERR_FILENOTFOUND)
	} else {
		v._file[channel] = f
		data := make([]byte, 1)
		_, _ = f.Read(data)
		v._read_char[channel] = data[0]
	}
	return virtualdrive.StOk
}

func (v *DriveHelper) Close(channel uint8) uint8 {
	v.LedTurnOff()
	if channel == 15 {
		v.closeAllChannels()
		return virtualdrive.StOk
	}
	if v._file[channel] != nil {
		v._file[channel].Close()
		v._file[channel] = nil
	}
	return virtualdrive.StOk
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

func (v *DriveHelper) changeDirectory(dirPath string) bool {
	d, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	if !d.IsDir() {
		return false
	}
	v._dir_path = dirPath
	v._dir_title = v._dir_path
	if len(v._dir_title) > 16 {
		v._dir_title = v._dir_title[:16]
	}
	return true
}

func (v *DriveHelper) findFirstFile(pattern string) (string, bool) {
	items, _ := os.ReadDir(v._dir_path)
	matcher := virtualdrive.NewMatcher()
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		if b, err := matcher.Match(pattern, item.Name()); err == nil && b {
			return item.Name(), false
		}
	}
	return "", false
}

func (v *DriveHelper) closeAllChannels() {
	for i := uint8(0); i < 15; i++ {
		v.Close(i)
	}
	v.cmd_len = 0
}

func (v *DriveHelper) openDirectory(channel uint8, name string) uint8 {
	return virtualdrive.StOk
}
