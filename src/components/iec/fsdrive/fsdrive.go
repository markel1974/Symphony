package fsdrive

import (
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/iec/iecdevice"
	"github.com/markel1974/c64emu/src/components/iec/iecprotocol"
	"github.com/markel1974/c64emu/src/config"
	"io"
	"os"
	"strings"
)

const (
	ATN_IN  = 0x80
	CLK_IN  = 0x04
	DATA_IN = 0x01

	DATA_OUT = 0x02
	CLK_OUT  = 0x08
	ATN_A    = 0x10
)

type FSDrive struct {
	*board.BaseComponent
	*iecprotocol.Protocol
	commands     *Commands
	deviceId     uint8
	deviceNumber uint8
	path         string
	dirPath      string       // Path to directory
	origDirPath  string       // Original directory path
	dirTitle     string       // Directory title
	file         [16]*os.File // File pointers for each of the 16 channels
	readChar     [16]uint8    // Buffers for one-byte read-ahead
	ready        bool
	cfg          *config.Config
}

func New(parent board.IComponent, suffix string, deviceId uint8, deviceNumber uint8, path string) *FSDrive {
	fs := &FSDrive{
		BaseComponent: board.NewBaseComponent("fs_drive", suffix),
		deviceId:      deviceId,
		deviceNumber:  deviceNumber,
		path:          path,
		commands:      NewCommands(),
		origDirPath:   "",
		cfg:           nil,
	}
	fs.Protocol = iecprotocol.NewProtocol(parent, suffix, deviceNumber, fs)
	board.Register(fs.Protocol, fs)
	return fs
}

func (v *FSDrive) Setup(iec iecdevice.IIec, cfg *config.Config) {
	v.Protocol.Setup(iec, cfg)
	v.cfg = cfg
	v.cfg.Bind(v.configChanged)
	v.origDirPath = v.path
	if v.changeDirectory(v.origDirPath) {
		for i := 0; i < 16; i++ {
			v.file[i] = nil
		}
		v.ready = true
	}
	v.Reset()
}

func (v *FSDrive) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

func (v *FSDrive) Reset() {
	v.closeAllChannels()
	v.commands.CommandClear()
	v.commands.SetError(ERR_STARTUP)

	v.Protocol.Reset()
	//TODO IN FASE DI RESET CAMBIARE LO STATO DEL BUS
}

func (v *FSDrive) LedTurnOn() {
}

func (v *FSDrive) LedTurnOff() {
}

func (v *FSDrive) GetPath() string {
	return v.dirPath
}

func (v *FSDrive) Listen(sec uint8) {
}

func (v *FSDrive) Unlisten(sec uint8) {
}

func (v *FSDrive) Talk(sec uint8) {
}

func (v *FSDrive) Untalk(sec uint8) {
}

func (v *FSDrive) Open(channel uint8) uint8 {
	//TODO DATA
	var data []uint8
	v.LedTurnOn()
	// Channel 15: Execute file name as command
	if channel == 15 {
		action, ok := v.commands.CommandExec(data)
		if ok {
			if action == 1 {
				v.Reset()
			}
		}
		return StOk
	}
	// Close previous file if still open
	if v.file[channel] != nil {
		v.file[channel].Close()
		v.file[channel] = nil
	}
	if data[0] == '#' {
		v.commands.SetError(ERR_NOCHANNEL)
		return StOk
	}
	if data[0] == '$' {
		return v.openDirectory(channel, string(data))
	}
	return v.openFile(channel, string(data))
}

func (v *FSDrive) Close(channel uint8) uint8 {
	v.LedTurnOff()
	if channel == 15 {
		v.closeAllChannels()
		return StOk
	}
	if v.file[channel] != nil {
		v.file[channel].Close()
		v.file[channel] = nil
	}
	return StOk
}

func (v *FSDrive) Read(channel uint8) (uint8, uint8) {
	// Channel 15: Error channel
	if channel == 15 {
		data := v.commands.RetrieveError()
		if data != '\r' {
			return data, StOk
		}
		// End of message
		v.commands.SetError(ERR_OK)
		return data, StEof
	}

	if v.file[channel] == nil {
		return 0, StReadTimeout
	}

	// Read one byte
	data := v.readChar[channel]
	buffer := make([]uint8, 1)
	c, err := v.file[channel].Read(buffer)
	if err == io.EOF {
		return data, StEof
	}
	v.readChar[channel] = (uint8)(c)
	return data, StOk
}

func (v *FSDrive) Write(channel uint8, data uint8) uint8 {
	//TODO EOI, eoi bool
	eoi := false
	// Channel 15: Collect chars and execute command on EOI
	if channel == 15 {
		if !v.commands.CommandSet(data) {

			return StTimeout
		}
		if eoi {
			v.commands.CommandExecBuf()
		}
		return StOk
	}
	if v.file[channel] == nil {
		v.commands.SetError(ERR_FILENOTOPEN)
		return StTimeout
	}
	if _, err := v.file[channel].Write([]byte{data}); err == io.EOF {
		v.commands.SetError(ERR_WRITE25)
		return StTimeout
	}
	return StOk
}

func (v *FSDrive) initializeCmd() {
	//v.closeAllChannels()
}

func (v *FSDrive) validateCmd() {
}

func (v *FSDrive) changeDirectory(dirPath string) bool {
	d, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	if !d.IsDir() {
		return false
	}
	v.dirPath = dirPath
	v.dirTitle = v.dirPath
	if len(v.dirTitle) > 16 {
		v.dirTitle = v.dirTitle[:16]
	}
	return true
}

func (v *FSDrive) findFirstFile(pattern string) (string, bool) {
	items, _ := os.ReadDir(v.dirPath)
	matcher := NewMatcher()
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

func (v *FSDrive) closeAllChannels() {
	for i := uint8(0); i < 15; i++ {
		v.Close(i)
	}
	v.commands.CommandClear()
}

func (v *FSDrive) openFile(channel uint8, name string) uint8 {
	plainName, mode, kind, _ := ParseFileName(name, true)
	// Channel 0 is READ, channel 1 is WRITE
	if channel == 0 || channel == 1 {
		mode = FMODE_READ
		if channel != 0 {
			mode = FMODE_WRITE
		}
		if kind == FTYPE_DEL {
			kind = FTYPE_PRG
		}
	}
	writing := mode == FMODE_WRITE || mode == FMODE_APPEND
	if strings.Contains(plainName, "*") || strings.Contains(plainName, "?") {
		if writing {
			v.commands.SetError(ERR_SYNTAX33)
			return StOk
		} else {
			v.findFirstFile(plainName)
		}
	}
	if kind == FTYPE_REL {
		v.commands.SetError(ERR_UNIMPLEMENTED)
		return StOk
	}
	flags := os.O_RDONLY
	perm := os.FileMode(0)
	switch mode {
	case FMODE_WRITE:
		perm = os.FileMode(0666)
		flags = os.O_RDWR
	case FMODE_APPEND:
		perm = os.FileMode(0666)
		flags = os.O_RDWR | os.O_APPEND
	default:
		panic("unhandled default case")
	}
	completeFileName := v.dirPath + string(os.PathSeparator) + plainName
	f, err := os.OpenFile(completeFileName, flags, perm)
	if err != nil {
		v.commands.SetError(ERR_FILENOTFOUND)
	} else {
		v.file[channel] = f
		data := make([]byte, 1)
		_, _ = f.Read(data)
		v.readChar[channel] = data[0]
	}
	return StOk
}

func (v *FSDrive) openDirectory(channel uint8, name string) uint8 {
	return StOk
}

func (v *FSDrive) configChanged() {
}

func data2string(data uint8) string {
	var message []string
	if data&ATN_IN != 0 {
		message = append(message, "[ATN_IN]")
	}
	if data&0x40 != 0 {
		message = append(message, "[UNKNOWN BIT 7]")
	}
	if data&0x20 != 0 {
		message = append(message, "[UNKNOWN BIT 6]")
	}
	if data&ATN_A != 0 {
		message = append(message, "[ATN_A]")
	}
	if data&CLK_OUT != 0 {
		message = append(message, "[CLK_OUT]")
	}
	if data&CLK_IN != 0 {
		message = append(message, "[CLK_IN]")
	}
	if data&DATA_OUT != 0 {
		message = append(message, "[DATA_OUT]")
	}
	if data&DATA_IN != 0 {
		message = append(message, "[DATA_IN]")
	}
	return strings.Join(message, " ")
}
