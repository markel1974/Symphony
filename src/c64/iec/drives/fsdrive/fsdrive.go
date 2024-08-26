package fsdrive

import (
	"fmt"
	"github.com/markel1974/c64emu/src/c64/iec/virtualdrive"
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

const (
	NOOP = -1
)

type FSDrive struct {
	iec            virtualdrive.IIec
	commands       *virtualdrive.Commands
	deviceId       uint8
	deviceNumber   uint8
	path           string
	respond        *Queue
	_dir_path      string       // Path to directory
	_orig_dir_path string       // Original directory path
	_dir_title     string       // Directory title
	_file          [16]*os.File // File pointers for each of the 16 channels
	_read_char     [16]uint8    // Buffers for one-byte read-ahead
	_ready         bool
	atn            bool
	state          int
	cfg            *config.Config
	//test           int64
}

func New(iec virtualdrive.IIec, deviceId uint8, deviceNumber uint8, path string) *FSDrive {
	v := &FSDrive{
		iec:            iec,
		deviceId:       deviceId,
		deviceNumber:   deviceNumber,
		path:           path,
		respond:        NewQueue(512),
		commands:       virtualdrive.NewCommands(),
		_orig_dir_path: "",
		atn:            false,
		state:          0,
		cfg:            nil,
	}
	return v
}

func (v *FSDrive) Setup(cfg *config.Config) {
	v.cfg = cfg
	v.cfg.Bind(v.configChanged)
	v._orig_dir_path = v.path
	if v.changeDirectory(v._orig_dir_path) {
		for i := 0; i < 16; i++ {
			v._file[i] = nil
		}
		v.Reset()
		v._ready = true
	}
}

func (v *FSDrive) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

func (v *FSDrive) Reset() {
	v.closeAllChannels()
	v.commands.CommandClear()
	v.commands.SetError(virtualdrive.ERR_STARTUP)

	//TODO IN FASE DI RESET CAMBIARE LO STATO DEL BUS
}

func (v *FSDrive) AtnStateChanged(atnPrev bool, atn bool) {
	//https://www.pagetable.com/?p=1135
	// All devices on the bus have to respond to ATN by pulling DATA within 1000 µs (“ATN Response Timing”),
	// and also eventually release CLK, because they are now receivers.
	// Devices usually implement this in hardware by automatically answering ATN=1 with DATA=1,
	// so that they can participate in receiving the command even when the CPU is busy and cannot be interrupted. /atn = !atn
	//atn = !atn
	fmt.Printf("ATN received %v", atn)
	v.atn = atn
	if atn {
		v.state = 1
		value := uint8(0)
		value |= DATA_OUT
		value |= CLK_OUT //ERROR
		//value |= ATN_A

		//value := uint8(DATA_OUT)
		//value = ^value

		v.respond.AddMulti(NOOP, 8)
		v.respond.Add(int(value))
		fmt.Printf("...atn is on, responding %03d [%08b]\n", value, value)
	} else {
		v.state = 0
		//v.respond.AddMulti(NOOP, 8)
		//v.respond.Add(int(0))
		fmt.Printf("...atn is now off\n")
	}
}

func (v *FSDrive) BusStateChanged(data uint8) {
	fmt.Printf("DTA received %03d [%08b] ATN: [%v], %s", data, data, v.atn, data2string(data))
	clkIn := data&CLK_IN != 0
	dataIn := data&DATA_IN != 0
	switch v.state {
	case 1:
		if clkIn && !dataIn {
			//value := 0xfe
			//value := uint8(DATA_OUT)

			//value := data | DATA_OUT
			//value = ^value

			//value := 0xfe
			value := ATN_A | DATA_OUT
			v.respond.AddMulti(NOOP, 16)
			v.respond.Add(value)
			v.state = 2
			fmt.Printf("...responding I'm Here %03d [%08b]\n", value, value)
			return
		}
	case 2:
		if !clkIn && dataIn {
			value := ATN_A
			v.respond.AddMulti(NOOP, 16)
			v.respond.Add(value)
			v.state = 3
			fmt.Printf("...responding Ready %03d [%08b]\n", value, value)
			return
		}
	}
	fmt.Printf("...NOT RESPONDING!\n")

}

func (v *FSDrive) Emulate() {
	if v.respond.Len() == 0 {
		return
	}
	data := v.respond.Next()
	if data == NOOP {
		return
	}
	fmt.Printf("SENDING %03d\n", uint8(data))
	//if time.Now().UnixMilli() > v.test { //qCycle&0x800 != 0 {
	v.iec.PeripheralWrite(v.deviceNumber, uint8(data))
	//d := v.iec.PeripheralRead(v.deviceNumber)
	//fmt.Printf("STATE: %s\n", strconv.FormatInt(int64(d), 2))
	//	v.test = 0
	//}
}

func (v *FSDrive) Ready() bool {
	return true
}

func (v *FSDrive) LedTurnOn() {
}

func (v *FSDrive) LedTurnOff() {
}

func (v *FSDrive) GetPath() string {
	return v._dir_path
}

func (v *FSDrive) Open(channel uint8, data []uint8) uint8 {
	v.LedTurnOn()
	// Channel 15: Execute file name as command
	if channel == 15 {
		action, ok := v.commands.CommandExec(data)
		if ok {
			if action == 1 {
				v.Reset()
			}
		}
		return virtualdrive.StOk
	}
	// Close previous file if still open
	if v._file[channel] != nil {
		v._file[channel].Close()
		v._file[channel] = nil
	}
	if data[0] == '#' {
		v.commands.SetError(virtualdrive.ERR_NOCHANNEL)
		return virtualdrive.StOk
	}
	if data[0] == '$' {
		return v.openDirectory(channel, string(data))
	}
	return v.openFile(channel, string(data))
}

func (v *FSDrive) openFile(channel uint8, name string) uint8 {
	plainName, mode, kind, _ := virtualdrive.ParseFileName(name, true)
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
			v.commands.SetError(virtualdrive.ERR_SYNTAX33)
			return virtualdrive.StOk
		} else {
			v.findFirstFile(plainName)
		}
	}
	if kind == virtualdrive.FTYPE_REL {
		v.commands.SetError(virtualdrive.ERR_UNIMPLEMENTED)
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
		v.commands.SetError(virtualdrive.ERR_FILENOTFOUND)
	} else {
		v._file[channel] = f
		data := make([]byte, 1)
		_, _ = f.Read(data)
		v._read_char[channel] = data[0]
	}
	return virtualdrive.StOk
}

func (v *FSDrive) Close(channel uint8) uint8 {
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

func (v *FSDrive) Read(channel uint8) (uint8, uint8) {
	// Channel 15: Error channel
	if channel == 15 {
		data := v.commands.RetrieveError()
		if data != '\r' {
			return virtualdrive.StOk, data
		}
		// End of message
		v.commands.SetError(virtualdrive.ERR_OK)
		return virtualdrive.StEof, data
	}

	if v._file[channel] == nil {
		return virtualdrive.StReadTimeout, 0
	}

	// Read one byte
	data := v._read_char[channel]
	buffer := make([]uint8, 1)
	c, err := v._file[channel].Read(buffer)
	if err == io.EOF {
		return virtualdrive.StEof, data
	}
	v._read_char[channel] = (uint8)(c)
	return virtualdrive.StOk, data
}

func (v *FSDrive) Write(channel uint8, data uint8, eoi bool) uint8 {
	// Channel 15: Collect chars and execute command on EOI
	if channel == 15 {
		if !v.commands.CommandSet(data) {
			return virtualdrive.StTimeout
		}
		if eoi {
			v.commands.CommandExecBuf()
		}
		return virtualdrive.StOk
	}
	if v._file[channel] == nil {
		v.commands.SetError(virtualdrive.ERR_FILENOTOPEN)
		return virtualdrive.StTimeout
	}
	if _, err := v._file[channel].Write([]byte{data}); err == io.EOF {
		v.commands.SetError(virtualdrive.ERR_WRITE25)
		return virtualdrive.StTimeout
	}
	return virtualdrive.StOk
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
	v._dir_path = dirPath
	v._dir_title = v._dir_path
	if len(v._dir_title) > 16 {
		v._dir_title = v._dir_title[:16]
	}
	return true
}

func (v *FSDrive) findFirstFile(pattern string) (string, bool) {
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

func (v *FSDrive) closeAllChannels() {
	for i := uint8(0); i < 15; i++ {
		v.Close(i)
	}
	v.commands.CommandClear()
}

func (v *FSDrive) openDirectory(channel uint8, name string) uint8 {
	return virtualdrive.StOk
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
