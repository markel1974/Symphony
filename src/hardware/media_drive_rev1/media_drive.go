package media_drive_rev1

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware/iec_rev1"
	"github.com/markel1974/c64emu/src/hardware/media_drive_rev1/adapters"
	"github.com/markel1974/c64emu/src/references"
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
	protocol       *iec_rev1.Protocol
	commands       *Commands
	deviceId       uint8
	deviceNumber   uint8
	path           string
	cfg            *config.Config
	channels       [16]*Channel
	adapterFactory *adapters.Factory
}

// NewBoard creates and initializes a new MediaDrive instance with the specified parent component, component factory, label, and instance number.
func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *MediaDrive {
	protocol := iec_rev1.NewProtocol(factory, parent, label, instance)
	fs := &MediaDrive{
		BaseComponent:  component.NewBaseComponent(),
		IIecDevice:     protocol,
		protocol:       protocol,
		deviceId:       0,
		deviceNumber:   0,
		commands:       NewCommands(),
		cfg:            nil,
		adapterFactory: adapters.NewFactory(),
		path:           "",
	}
	fs.BaseComponent.Register(factory, fs.protocol, Identifier(), fs, references.IdIIecProtocolDevice(fs, label, 0))
	fs.protocol.SetDevice(fs)
	adapter := fs.adapterFactory.Void()
	for idx := range fs.channels {
		fs.channels[idx] = NewChannel(idx, adapter)
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
	v.path = path
	if err := v.protocol.Bind(v, deviceId, deviceNumber); err != nil {
		return err
	}
	adapter, err := v.adapterFactory.Create(path)
	if err != nil {
		return err
	}
	for idx := range v.channels {
		v.channels[idx].SetAdapter(adapter)
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

// Reset reinitializes the MediaDrive by closing all channels, clearing commands, setting the error index, and resetting the protocol.
func (v *MediaDrive) Reset() {
	v.LedActivity(false)
	for i := uint8(0); i < uint8(len(v.channels)); i++ {
		v.channels[i].Reset()
	}
	v.commands.CommandClear()
	v.channels[errChannel].SetError(adapters.Error(adapters.ErrStartup))
	v.protocol.Reset()
}

// Internal returns a boolean indicating whether the MediaDrive is an internal device. Always returns false.
func (v *MediaDrive) Internal() bool {
	return false
}

// GetPath returns the name of the adapter associated with the MediaDrive.
func (v *MediaDrive) GetPath() string {
	return v.path
}

// GetDeviceNumber retrieves the device number associated with the MediaDrive instance.
func (v *MediaDrive) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

// Talk sends a command to a specific device channel and returns a status code indicating the operation result.
func (v *MediaDrive) Talk(_ uint8) uint8 {
	//channel := d & 0xf
	return adapters.StOk
}

// Untalk disables the "talk" state for a given channel identified by the lower 4 bits of the input parameter d.
// Returns StOk on successful execution.
func (v *MediaDrive) Untalk(_ uint8) uint8 {
	//channel := d & 0xf
	return adapters.StOk
}

// Listen processes the device identifier and returns the standard operational status code.
func (v *MediaDrive) Listen(_ uint8) uint8 {
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

	channelId := d & 0xf
	channel := v.channels[channelId]

	if openMode := channel.OpenModeGet() & 0xf0; openMode != 0x20 && openMode != 0xf0 {
		return adapters.StOk
	}
	v.LedActivity(true)
	data := channel.Buffer()
	if channelId == errChannel {
		action, err := v.commands.CommandExec(data)
		if err != nil {
			v.channels[errChannel].SetError(err)
			return adapters.StOk
		}
		if action == 1 {
			for i := uint8(0); i < uint8(len(v.channels)); i++ {
				v.channels[i].Reset()
			}
			v.commands.CommandClear()
		}
		return adapters.StOk
	}

	channel.Reset()
	if len(data) == 0 {
		v.channels[errChannel].SetError(adapters.Error(adapters.ErrNoChannel))
		return adapters.StOk
	}
	if data[0] == '#' {
		v.channels[errChannel].SetError(adapters.Error(adapters.ErrNoChannel))
		return adapters.StOk
	}
	if data[0] == '$' {
		if err := channel.OpenDirectory(string(data)); err != nil {
			v.channels[errChannel].SetError(err)
			return adapters.StOk
		}
		return adapters.StOk
	}
	if err := channel.OpenFile(string(data)); err != nil {
		v.channels[errChannel].SetError(err)
		return adapters.StOk
	}
	return adapters.StOk
}

// Open initializes the specified channel in the MediaDrive, setting its mode and resetting its state.
func (v *MediaDrive) Open(d uint8) uint8 {
	channelId := d & 0xf
	channel := v.channels[channelId]
	channel.Reset()
	channel.OpenModeSet(d)
	return adapters.StOk
}

// Close shuts down the specified channel `d` and turns off the LED. If `d` equals 15, all channels are closed. Returns status.
func (v *MediaDrive) Close(d uint8) uint8 {
	channelId := d & 0xf
	v.LedActivity(false)
	if channelId == errChannel {
		for i := uint8(0); i < uint8(len(v.channels)); i++ {
			_ = v.channels[i].Close()
		}
		v.commands.CommandClear()
		return adapters.StOk
	}
	if err := v.channels[channelId].Close(); err != nil {
		v.channels[errChannel].SetError(err)
	}
	return adapters.StOk
}

// Read retrieves the next byte of data from the specified channel and returns its status.
// If the channel is 15, an error is processed. When data is available, it is returned along with the status.
// Returns StReadTimeout if no data is available or StEof if the channel is empty.
func (v *MediaDrive) Read(d uint8) (uint8, uint8) {
	channelId := d & 0xf
	channel := v.channels[channelId]
	data, ok := channel.Read()
	if !ok {
		return 0, adapters.StReadTimeout
	}
	if channel.ReadIsEmpty() {
		v.channels[errChannel].SetError(adapters.Error(adapters.ErrOk))
		return data, adapters.StEof
	}
	return data, adapters.StOk
}

// Write writes a byte of data to a specific channel and executes commands for channel 15, returning a status code.
func (v *MediaDrive) Write(d uint8, data uint8) uint8 {
	//TODO EOI
	channelId := d & 0xf
	channel := v.channels[channelId]
	eoi := false
	if channelId == errChannel {
		if !v.commands.CommandSet(data) {
			return adapters.StTimeout
		}
		if eoi {
			if _, err := v.commands.CommandExecBuf(); err != nil {
				v.channels[errChannel].SetError(err)
			} else {
				v.channels[errChannel].SetError(adapters.Error(adapters.ErrOk))
			}
		}
		return adapters.StOk
	}
	channel.BufferAdd(data)
	return adapters.StOk
}

// configChanged is a callback invoked when the component's configuration changes to apply updates dynamically.
func (v *MediaDrive) configChanged() {
}
