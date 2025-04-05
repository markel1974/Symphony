package media_drive

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/iec"
	"github.com/markel1974/c64emu/src/hardware/media_drive/adapters"
	"github.com/markel1974/c64emu/src/references"
	"strings"

	"github.com/markel1974/c64emu/src/config"
)

const errChannel = 15

// MediaDrive represents a media driving component with protocol support and channel management.
// It embeds BaseComponent and IIecDevice, providing IEC device capabilities.
// MediaDrive includes protocol handling, command execution, and configuration settings.
// It manages up to 16 communication channels and an adapter for device interactions.
// Device ID and number are associated with the drive for identification purposes.
type MediaDrive struct {
	*component.BaseComponent
	references.IIecDevice
	protocol       *iec.Protocol
	commands       *Commands
	deviceId       uint8
	deviceNumber   uint8
	cfg            *config.Config
	channels       [16]*Channel
	adapter        adapters.IAdapter
	adapterFactory *adapters.Factory
	matcher        *Matcher
}

// NewBoard creates and initializes a new MediaDrive instance with the specified parent component, component factory, label, and instance number.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *MediaDrive {
	protocol := iec.NewProtocol(factory, parent, label, instance)
	fs := &MediaDrive{
		BaseComponent:  component.NewBaseComponent(),
		IIecDevice:     protocol,
		protocol:       protocol,
		deviceId:       0,
		deviceNumber:   0,
		commands:       NewCommands(),
		cfg:            nil,
		adapterFactory: adapters.NewFactory(),
		adapter:        nil,
		matcher:        NewMatcher(),
	}
	fs.BaseComponent.Register(factory, fs.protocol, Identifier(), fs, references.IdIIecProtocolDevice(fs, label, 0))
	fs.protocol.SetDevice(fs)
	fs.adapter = fs.adapterFactory.Void()
	for idx := range fs.channels {
		fs.channels[idx] = NewChannel()
	}
	return fs
}

// Setup initializes the MediaDrive by configuring its settings and preparing the protocol for use.
func (v *MediaDrive) Setup() error {
	v.cfg = v.GetFactory().GetConfig()
	v.cfg.Bind(v.configChanged)
	if err := v.protocol.Setup(); err != nil {
		return err
	}
	return nil
}

// Bind binds the MediaDrive to the specified device socket, updates its device ID and number, and initializes the adapter and channels.
func (v *MediaDrive) Bind(_ references.IIecDeviceSocket, deviceId uint8, deviceNumber uint8) error {
	path := ""
	if d := v.cfg.Drive(v.deviceId); d != nil {
		path = d.GetId()
	}
	v.deviceId = deviceId
	v.deviceNumber = deviceNumber
	if err := v.protocol.Bind(v, deviceId, deviceNumber); err != nil {
		return err
	}
	adapter, err := v.adapterFactory.Create(path)
	if err != nil {
		return err
	}
	v.adapter = adapter
	for idx := range v.channels {
		v.channels[idx].Reset()
	}
	return nil
}

// Connect establishes a connection using the configured protocol and returns an error if the connection fails.
func (v *MediaDrive) Connect() error {
	if err := v.protocol.Connect(); err != nil {
		return err
	}
	return nil
}

// Internal returns a boolean indicating whether the MediaDrive is an internal device. Always returns false.
func (v *MediaDrive) Internal() bool {
	return false
}

// GetDeviceNumber retrieves the device number associated with the MediaDrive instance.
func (v *MediaDrive) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

// Reset reinitializes the MediaDrive by closing all channels, clearing commands, setting the error index, and resetting the protocol.
func (v *MediaDrive) Reset() {
	v.LedTurnOff()
	for i := uint8(0); i < uint8(len(v.channels)); i++ {
		v.channels[i].Close()
	}
	v.commands.CommandClear()
	v.channels[errChannel].DataSetError(adapters.Error(adapters.ErrStartup))
	v.protocol.Reset()
}

// LedTurnOn activates the LED of the MediaDrive, indicating an active or operational state.
func (v *MediaDrive) LedTurnOn() {
	v.LEDSignal().Emit(uint32(1)<<8 | uint32(v.deviceNumber))
}

// LedTurnOff turns off the LED indicator for the MediaDrive.
func (v *MediaDrive) LedTurnOff() {
	v.LEDSignal().Emit(uint32(0)<<8 | uint32(v.deviceNumber))
}

// GetPath returns the name of the adapter associated with the MediaDrive.
func (v *MediaDrive) GetPath() string {
	return v.adapter.Name()
}

// Listen processes the device identifier and returns the standard operational status code.
func (v *MediaDrive) Listen(d uint8) uint8 {
	//channel := d & 0xf
	return adapters.StOk
}

// Unlisten handles the UNLISTEN command for a specified channel in the MediaDrive.
// Returns a status code indicating the operation's success or failure.
func (v *MediaDrive) Unlisten(d uint8) uint8 {
	//TODO
	//If this is an UNLISTEN that followed an OPEN (0x2_ 0xf_), then
	//device.unlisten will try to open the file with the filename that
	//was received in between the OPEN and now.
	//If the file cannot be opened, it will set st != 0.

	channel := d & 0xf

	mode := v.channels[channel].ModeGet() & 0xf0
	if mode != 0x20 && mode != 0xf0 {
		return adapters.StOk
	}

	data := v.channels[channel].BufferGet()

	v.LedTurnOn()
	// Channel 15: Execute file name as command
	if channel == errChannel {
		if action, err := v.commands.CommandExec(data); err != nil {
			v.channels[errChannel].DataSetError(err)
		} else {
			if action == 1 {
				v.Reset()
			}
		}
		return adapters.StOk
	}
	v.channels[channel].Close()

	if len(data) == 0 {
		v.channels[errChannel].DataSetError(adapters.Error(adapters.ErrNoChannel))
		return adapters.StOk
	}
	if data[0] == '#' {
		v.channels[errChannel].DataSetError(adapters.Error(adapters.ErrNoChannel))
		return adapters.StOk
	}
	if data[0] == '$' {
		dirData, err := v.openDirectory("")
		if err != nil {
			v.channels[errChannel].DataSetError(err)
			return adapters.StOk
		}
		v.channels[channel].DataSet(dirData)
		return adapters.StOk
	}
	fileData, err := v.openFile(channel, string(data))
	if err != nil {
		v.channels[errChannel].DataSetError(err)
		return adapters.StOk
	}
	v.channels[channel].DataSet(fileData)
	return adapters.StOk
}

// Talk sends a command to a specific device channel and returns a status code indicating the operation result.
func (v *MediaDrive) Talk(d uint8) uint8 {
	//channel := d & 0xf
	return adapters.StOk
}

// Untalk disables the "talk" state for a given channel identified by the lower 4 bits of the input parameter d.
// Returns StOk on successful execution.
func (v *MediaDrive) Untalk(d uint8) uint8 {
	//channel := d & 0xf
	return adapters.StOk
}

// Open initializes the specified channel in the MediaDrive, setting its mode and resetting its state.
func (v *MediaDrive) Open(d uint8) uint8 {
	channel := d & 0xf
	v.channels[channel].Reset()
	v.channels[channel].ModeSet(d)
	return adapters.StOk
}

// Close shuts down the specified channel `d` and turns off the LED. If `d` equals 15, all channels are closed. Returns status.
func (v *MediaDrive) Close(d uint8) uint8 {
	channel := d & 0xf
	v.LedTurnOff()
	if channel == errChannel {
		for i := uint8(0); i < uint8(len(v.channels)); i++ {
			v.channels[i].Close()
		}
		v.commands.CommandClear()
		return adapters.StOk
	}
	v.channels[channel].Close()
	return adapters.StOk
}

// Read retrieves the next byte of data from the specified channel and returns its status.
// If the channel is 15, an error is processed. When data is available, it is returned along with the status.
// Returns StReadTimeout if no data is available or StEof if the channel is empty.
func (v *MediaDrive) Read(d uint8) (uint8, uint8) {
	channel := d & 0xf
	b, ok := v.channels[channel].DataNext()
	if !ok {
		return 0, adapters.StReadTimeout
	}
	if v.channels[channel].DataIsEmpty() {
		v.channels[errChannel].DataSetError(adapters.Error(adapters.ErrOk))
		return b, adapters.StEof
	}
	return b, adapters.StOk
}

// Write writes a byte of data to a specific channel and executes commands for channel 15, returning a status code.
func (v *MediaDrive) Write(d uint8, data uint8) uint8 {
	//TODO EOI, eoi bool
	channel := d & 0xf
	eoi := false

	if channel == errChannel {
		if !v.commands.CommandSet(data) {
			return adapters.StTimeout
		}
		if eoi {
			if _, err := v.commands.CommandExecBuf(); err != nil {
				v.channels[errChannel].DataSetError(err)
			} else {
				v.channels[errChannel].DataSetError(adapters.Error(adapters.ErrOk))
			}
		}
		return adapters.StOk
	}
	//if v.buffer[channel] == nil {
	//	v.commands.SetError(ErrFileNotOpen)
	//	return StTimeout
	//}

	v.channels[channel].BufferAdd(data)
	return adapters.StOk

	/*
		//TODO EOI, eoi bool
		eoi := false

		log.Printf("MediaDrive received: %s\n", string(data))

		return StOk

		// Channel 15: Collect chars and execute command on EOI
		if channel == errChannel {
			if !v.commands.CommandSet(data) {

				return StTimeout
			}
			if eoi {
				v.commands.CommandExecBuf()
			}
			return StOk
		}
		if v.file[channel] == nil {
			v.commands.SetError(ErrFileNotOpen)
			return StTimeout
		}
		if _, err := v.file[channel].Write([]byte{data}); err == io.EOF {
			v.commands.SetError(ErrWrite25)
			return StTimeout
		}
		return StOk

	*/
}

// initializeCmd sets up and initializes necessary commands or operations for the MediaDrive instance.
func (v *MediaDrive) initializeCmd() {
	//v.closeAllChannels()
}

// validateCmd ensures that the current state or command of the MediaDrive is valid according to its protocol or configuration.
func (v *MediaDrive) validateCmd() {
}

// openFile attempts to open a file based on the specified channel and name, returning its content or an error if unsuccessful.
// It handles file name parsing, mode checks (read/write), and wildcards, returning appropriate errors for invalid cases.
// Errors include syntax issues, file not found, or unimplemented file types like relative files.
func (v *MediaDrive) openFile(channel uint8, name string) ([]uint8, error) {
	plainName, mode, kind, _ := adapters.ParseFileName(name)
	// Channel 0 is READ, channel 1 is WRITE
	if channel == 0 || channel == 1 {
		mode = adapters.FModeRead
		if channel != 0 {
			mode = adapters.FModeWrite
		}
		if kind == adapters.FTypeDel {
			kind = adapters.FTypePrg
		}
	}
	if v.matcher.Contains(plainName) {
		if mode == adapters.FModeWrite || mode == adapters.FModeAppend {
			return nil, adapters.Error(adapters.ErrSyntax33)
		}
		items, err := v.adapter.ReadDir()
		if err != nil {
			return nil, adapters.Error(adapters.ErrFileNotFound)
		}
		found := false
		for _, item := range items {
			if !item.IsDir() {
				if found = v.matcher.Match(plainName, item.Name()); found {
					plainName = item.Name()
					break
				}
			}
		}
		if !found {
			return nil, adapters.Error(adapters.ErrFileNotFound)
		}
	}
	if kind == adapters.FTypeRel {
		return nil, adapters.Error(adapters.ErrUnimplemented)
	}
	data, err := v.adapter.ReadFile(plainName)
	if err != nil {
		return nil, adapters.Error(adapters.ErrFileNotFound)
	}
	return data, nil
}

// openDirectory generates a directory listing based on a pattern and returns it as a byte slice, or returns an error if failed.
func (v *MediaDrive) openDirectory(pattern string) ([]byte, error) {
	const titleStart = "\001\004\001\001\000\000\022\""
	const titleEnd = "\" 00 2A"
	const blocksFreeStart = "\001\001\000\000"
	const blockFreeEnd = "\000\000"
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
	title := adapters.CreateFileNameFilled(v.adapter.Name(), ' ')
	fullTile := titleStart + string(title) + titleEnd
	var buf []byte
	buf = append(buf, fullTile...)
	buf = append(buf, 0)

	entries, err := v.adapter.ReadDir()
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		z := v.createFileEntry(e.Name(), int(e.Size()), 0)
		buf = append(buf, z...)
	}
	buf = append(buf, blocksFreeStart...)
	buf = append(buf, "BLOCKS FREE.             "...)
	buf = append(buf, blockFreeEnd...)
	buf = append(buf, 0)

	return buf, nil
}

// createFileEntry generates a directory entry for a file with the specified name, size, and type, returning a byte slice.
func (v *MediaDrive) createFileEntry(name string, size int, kind int) []byte {
	const dirEntryMax = 32
	vName := adapters.CreateFileName(name)
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

// configChanged is a callback invoked when the component's configuration changes to apply updates dynamically.
func (v *MediaDrive) configChanged() {
}
