package fs_drive

import (
	"errors"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/iec"
	"github.com/markel1974/c64emu/src/references"
	"os"
	"strings"

	"github.com/markel1974/c64emu/src/config"
)

type FSDrive struct {
	*component.BaseComponent
	references.IIecDevice
	protocol     *iec.Protocol
	commands     *Commands
	deviceId     uint8
	deviceNumber uint8
	dirPath      string
	cfg          *config.Config
	channels     [16]*Channel
}

func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *FSDrive {
	protocol := iec.NewProtocol(factory, parent, label, instance)
	fs := &FSDrive{
		BaseComponent: component.NewBaseComponent(),
		IIecDevice:    protocol,
		protocol:      protocol,
		deviceId:      0,
		deviceNumber:  0,
		commands:      NewCommands(),
		cfg:           nil,
	}
	fs.BaseComponent.Register(factory, fs.protocol, Identifier(), fs, references.IdIIecProtocolDevice(fs, label, 0))
	fs.protocol.SetDevice(fs)

	for idx := range fs.channels {
		fs.channels[idx] = NewChannel()
	}

	return fs
}

func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IIecDevice {
	return NewBoard(parent, factory, label, instance)
}

func (v *FSDrive) Setup() error {
	v.cfg = v.GetFactory().GetConfig()
	v.cfg.Bind(v.configChanged)
	if err := v.protocol.Setup(); err != nil {
		return err
	}
	return nil
}

func (v *FSDrive) Bind(_ references.IIecDeviceSocket, deviceId uint8, deviceNumber uint8) error {
	path := ""
	if d := v.cfg.Drive(v.deviceId); d != nil {
		path = d.GetId()
	}
	v.deviceId = deviceId
	v.deviceNumber = deviceNumber
	if err := v.protocol.Bind(v, deviceId, deviceNumber); err != nil {
		return err
	}
	if v.changeDirectory(path) {
		for idx := range v.channels {
			v.channels[idx].Reset()
		}
	}
	return nil
}

func (v *FSDrive) Connect() error {
	if err := v.protocol.Connect(); err != nil {
		return err
	}
	return nil
}

func (v *FSDrive) Internal() bool {
	return false
}

func (v *FSDrive) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

func (v *FSDrive) Reset() {
	v.closeAllChannels()
	v.commands.CommandClear()
	v.commands.SetErrorIdx(ERR_STARTUP)
	v.protocol.Reset()
}

func (v *FSDrive) LedTurnOn() {
}

func (v *FSDrive) LedTurnOff() {
}

func (v *FSDrive) GetPath() string {
	return v.dirPath
}

func (v *FSDrive) Listen(d uint8) uint8 {
	//channel := d & 0xf
	return StOk
}

func (v *FSDrive) Unlisten(d uint8) uint8 {
	//TODO
	//If this is an UNLISTEN that followed an OPEN (0x2_ 0xf_), then
	//device.unlisten will try to open the file with the filename that
	//was received in between the OPEN and now.
	//If the file cannot be opened, it will set st != 0.

	channel := d & 0xf

	mode := v.channels[channel].ModeGet() & 0xf0
	if mode != 0x20 && mode != 0xf0 {
		return StOk
	}

	data := v.channels[channel].BufferGet()

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

	v.channels[channel].Close()

	if len(data) == 0 {
		v.commands.SetErrorIdx(ERR_NOCHANNEL)
		return StOk
	}
	if data[0] == '#' {
		v.commands.SetErrorIdx(ERR_NOCHANNEL)
		return StOk
	}
	if data[0] == '$' {
		dirData, err := v.openDirectory("", v.dirPath)
		if err != nil {
			v.commands.SetError(err)
			return StOk
		}
		v.channels[channel].DataSet(dirData)
		return StOk
	}
	fileData, err := v.openFile(channel, string(data))
	if err != nil {
		v.commands.SetError(err)
		return StOk
	}
	v.channels[channel].DataSet(fileData)
	return StOk
}

func (v *FSDrive) Talk(d uint8) uint8 {
	//channel := d & 0xf
	return StOk
}

func (v *FSDrive) Untalk(d uint8) uint8 {
	//channel := d & 0xf
	return StOk
}

func (v *FSDrive) Open(d uint8) uint8 {
	channel := d & 0xf
	v.channels[channel].Reset()
	v.channels[channel].ModeSet(d)
	return StOk
}

func (v *FSDrive) Close(d uint8) uint8 {
	channel := d & 0xf
	v.LedTurnOff()
	if channel == 15 {
		v.closeAllChannels()
		return StOk
	}
	v.channels[channel].Close()
	return StOk
}

func (v *FSDrive) Read(d uint8) (uint8, uint8) {
	channel := d & 0xf
	if channel == 15 {
		//TODO ERROR channel
		//data := v.commands.RetrieveError()
		//if data != '\r' {
		//	return data, StOk
		//}
		// End of message
		//v.commands.SetError(ERR_OK)
		//return data, StEof
	}
	b, ok := v.channels[channel].DataNext()
	if !ok {
		return 0, StReadTimeout
	}
	if v.channels[channel].DataIsEmpty() {
		v.commands.SetErrorIdx(ERR_OK)
		return b, StEof
	}
	return b, StOk
}

func (v *FSDrive) Write(d uint8, data uint8) uint8 {
	//TODO EOI, eoi bool
	channel := d & 0xf
	eoi := false

	if channel == 15 {
		if !v.commands.CommandSet(data) {
			return StTimeout
		}
		if eoi {
			v.commands.CommandExecBuf()
		}
		return StOk
	}
	//if v.buffer[channel] == nil {
	//	v.commands.SetError(ERR_FILENOTOPEN)
	//	return StTimeout
	//}

	v.channels[channel].BufferAdd(data)
	return StOk

	/*
		//TODO EOI, eoi bool
		eoi := false

		log.Printf("FSDrive received: %s\n", string(data))

		return StOk

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

	*/
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

func (v *FSDrive) openFile(channel uint8, name string) ([]uint8, error) {
	plainName, mode, kind, _ := ParseFileName(name)
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
	if strings.Contains(plainName, "*") || strings.Contains(plainName, "?") {
		if mode == FMODE_WRITE || mode == FMODE_APPEND {
			return nil, errors.New(string(Errors[ERR_SYNTAX33]))
		}
		n, ok := v.findFirstFile(plainName)
		if !ok {
			return nil, errors.New(string(Errors[ERR_FILENOTFOUND]))
		}
		plainName = n
	}
	if kind == FTYPE_REL {
		return nil, errors.New(string(Errors[ERR_UNIMPLEMENTED]))
	}
	completeFileName := v.dirPath + string(os.PathSeparator) + plainName
	data, err := os.ReadFile(completeFileName)
	if err != nil {
		return nil, errors.New(string(Errors[ERR_FILENOTFOUND]))
	}
	return data, nil
}

func (v *FSDrive) openDirectory(pattern string, dirName string) ([]byte, error) {
	// Special treatment for "$0"
	if len(pattern) > 0 {
		if pattern[0] == '0' && len(pattern) == 1 {
			pattern = ""
		}
	}

	if p := strings.Index(pattern, ":"); p >= 0 {
		p++
		if len(pattern) < p {
			pattern = pattern[p:]
		}
	}

	const titleStart = "\001\004\001\001\000\000\022\""
	const titleEnd = "\" 00 2A"
	const blocksFreeStart = "\001\001\000\000"
	const blockFreeEnd = "\000\000"

	title := CreateFileNameFilled(dirName, ' ')

	var buf []byte

	fullTile := titleStart + string(title) + titleEnd

	buf = append(buf, fullTile...)
	buf = append(buf, 0)

	entries, err := os.ReadDir(v.dirPath)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		fInfo, err := e.Info()
		if err != nil {
			continue
		}
		if fInfo.IsDir() {
			continue
		}
		z := v.createFileEntry(e.Name(), int(fInfo.Size()), 0)
		buf = append(buf, z...)
	}
	buf = append(buf, blocksFreeStart...)
	buf = append(buf, "BLOCKS FREE.             "...)
	buf = append(buf, blockFreeEnd...)
	buf = append(buf, 0)

	return buf, nil
}

func (v *FSDrive) createFileEntry(name string, size int, kind int) []byte {
	const dirEntryMax = 32
	vName := CreateFileName(name)
	n := (size + 254) / 254
	ret := make([]byte, dirEntryMax)
	for x := range ret {
		ret[x] = ' '
	}
	ret[0] = 0x1
	ret[1] = 0x1
	ret[2] = uint8(n & 0xff)
	ret[3] = uint8((n >> 8) & 0xff)
	nameIdx := 4
	if n < 10 {
		nameIdx++
	}
	if n < 100 {
		nameIdx++
	}
	ret[nameIdx] = '"'
	nameIdx++
	ret[nameIdx+len(vName)] = '"'
	for x, i := range vName {
		ret[nameIdx+x] = i
	}
	ret[28] = 'P'
	ret[29] = 'R'
	ret[30] = 'G'
	ret[31] = 0
	return ret
}

func (v *FSDrive) configChanged() {
}
