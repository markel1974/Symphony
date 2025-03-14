package iec

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	c1541board "github.com/markel1974/c64emu/src/hardware/c1541/board"
	"github.com/markel1974/c64emu/src/hardware/iec/fsdrive"
	"github.com/markel1974/c64emu/src/hardware/iec/iecdevice"
	"github.com/markel1974/c64emu/src/references"
	"strconv"
)

// DriveFactory represents a component responsible for creating and managing drive configurations and operations.
type DriveFactory struct {
	*component.BaseComponent
	cfg    *config.Config
	quartz references.IQuartzSocket
}

// NewDriveFactory creates and initializes a new DriveFactory instance, registers it with the provided parent component.
func NewDriveFactory(parent component.IComponent, suffix string, q references.IQuartzSocket) *DriveFactory {
	df := &DriveFactory{
		BaseComponent: component.NewBaseComponent("iec_drive_factory", suffix),
		quartz:        q,
		cfg:           nil,
	}
	component.Register(parent, df)
	return df
}

// Reset clears the internal state of the DriveFactory and prepares it for reinitialization or reuse.
func (c *DriveFactory) Reset() {
}

// Setup initializes the DriveFactory with the provided configuration and binds the configChanged function to the config signal.
func (c *DriveFactory) Setup(cfg *config.Config) {
	c.cfg = cfg
	c.cfg.Bind(c.configChanged)
}

// configChanged handles updates to the configuration when bound to configuration changes in the Setup method.
func (c *DriveFactory) configChanged() {
	//TODO IMPLEMENT
}

// Create instantiating and returns a new virtual drive instance based on the specified type, options, and device ID.
func (c *DriveFactory) Create(kind string, opts string, deviceId uint8) iecdevice.IIecDevice {
	deviceNumber := deviceId + 8
	suffix := strconv.Itoa(int(deviceNumber))
	var vd iecdevice.IIecDevice
	switch kind {
	case "C1541":
		vd = c1541board.New(c, suffix, c.quartz, deviceId, deviceNumber, opts)
	case "FSDRIVE":
		vd = fsdrive.New(c, suffix, c.quartz, deviceId, deviceNumber, opts)
	default:
		vd = c1541board.New(c, suffix, c.quartz, deviceId, deviceNumber, opts)
	}
	return vd
}
