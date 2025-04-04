package media_drive

import (
	"errors"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/iec"
	"github.com/markel1974/c64emu/src/hardware/media_drive/adapters"
	"github.com/markel1974/c64emu/src/references"
	"strings"

	"github.com/markel1974/c64emu/src/config"
)

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
	v.closeAllChannels()
	v.commands.CommandClear()
	v.commands.SetErrorIdx(ERR_STARTUP)
	v.protocol.Reset()
}

// LedTurnOn activates the LED of the MediaDrive, indicating an active or operational state.
func (v *MediaDrive) LedTurnOn() {
}

// LedTurnOff turns off the LED indicator for the MediaDrive.
func (v *MediaDrive) LedTurnOff() {
}

// GetPath returns the name of the adapter associated with the MediaDrive.
func (v *MediaDrive) GetPath() string {
	return v.adapter.Name()
}

// Listen processes the device identifier and returns the standard operational status code.
func (v *MediaDrive) Listen(d uint8) uint8 {
	//channel := d & 0xf
	return StOk
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
		dirData, err := v.openDirectory("")
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

// Talk sends a command to a specific device channel and returns a status code indicating the operation result.
func (v *MediaDrive) Talk(d uint8) uint8 {
	//channel := d & 0xf
	return StOk
}

// Untalk disables the "talk" state for a given channel identified by the lower 4 bits of the input parameter d.
// Returns StOk on successful execution.
func (v *MediaDrive) Untalk(d uint8) uint8 {
	//channel := d & 0xf
	return StOk
}

// Open initializes the specified channel in the MediaDrive, setting its mode and resetting its state.
func (v *MediaDrive) Open(d uint8) uint8 {
	channel := d & 0xf
	v.channels[channel].Reset()
	v.channels[channel].ModeSet(d)
	return StOk
}

// Close shuts down the specified channel `d` and turns off the LED. If `d` equals 15, all channels are closed. Returns status.
func (v *MediaDrive) Close(d uint8) uint8 {
	channel := d & 0xf
	v.LedTurnOff()
	if channel == 15 {
		v.closeAllChannels()
		return StOk
	}
	v.channels[channel].Close()
	return StOk
}

// Read retrieves the next byte of data from the specified channel and returns its status.
// If the channel is 15, an error is processed. When data is available, it is returned along with the status.
// Returns StReadTimeout if no data is available or StEof if the channel is empty.
func (v *MediaDrive) Read(d uint8) (uint8, uint8) {
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

// Write writes a byte of data to a specific channel and executes commands for channel 15, returning a status code.
func (v *MediaDrive) Write(d uint8, data uint8) uint8 {
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

		log.Printf("MediaDrive received: %s\n", string(data))

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

// initializeCmd sets up and initializes necessary commands or operations for the MediaDrive instance.
func (v *MediaDrive) initializeCmd() {
	//v.closeAllChannels()
}

// validateCmd ensures that the current state or command of the MediaDrive is valid according to its protocol or configuration.
func (v *MediaDrive) validateCmd() {
}

// findFirstFile searches for the first file matching the given pattern in the directory and returns its name and success status.
func (v *MediaDrive) findFirstFile(pattern string) (string, bool) {
	items, err := v.adapter.ReadDir()
	if err != nil {
		return "", false
	}
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

// closeAllChannels closes all active channels and clears the command queue in the MediaDrive instance.
func (v *MediaDrive) closeAllChannels() {
	for i := uint8(0); i < 15; i++ {
		v.Close(i)
	}
	v.commands.CommandClear()
}

// openFile attempts to open a file based on the specified channel and name, returning its content or an error if unsuccessful.
// It handles file name parsing, mode checks (read/write), and wildcards, returning appropriate errors for invalid cases.
// Errors include syntax issues, file not found, or unimplemented file types like relative files.
func (v *MediaDrive) openFile(channel uint8, name string) ([]uint8, error) {
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

	data, err := v.adapter.ReadFile(plainName)
	if err != nil {
		return nil, errors.New(string(Errors[ERR_FILENOTFOUND]))
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
	title := CreateFileNameFilled(v.adapter.Name(), ' ')
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

// configChanged is a callback invoked when the component's configuration changes to apply updates dynamically.
func (v *MediaDrive) configChanged() {
}
