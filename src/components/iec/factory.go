package iec

import (
	c1541board "github.com/markel1974/c64emu/src/c1541/board"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/iec/fsdrive"
	"github.com/markel1974/c64emu/src/components/iec/virtualdrive"
	"github.com/markel1974/c64emu/src/config"
	"strconv"
)

// DriveFactory represents a component responsible for creating and managing drive configurations and operations.
type DriveFactory struct {
	*board.BaseComponent
	cfg *config.Config
}

// NewDriveFactory creates and initializes a new DriveFactory instance, registers it with the provided parent component.
func NewDriveFactory(parent board.IComponent, suffix string) *DriveFactory {
	df := &DriveFactory{
		BaseComponent: board.NewBaseComponent("drive_factory", suffix),
		cfg:           nil,
	}
	board.Register(parent, df)
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

// Create instantiates and returns a new virtual drive instance based on the specified type, options, and device ID.
func (c *DriveFactory) Create(kind string, opts string, deviceId uint8) virtualdrive.IVirtualDrive {
	deviceNumber := deviceId + 8
	suffix := strconv.Itoa(int(deviceNumber))
	var vd virtualdrive.IVirtualDrive
	switch kind {
	case "C1541":
		vd = c1541board.New(c, suffix, deviceId, deviceNumber, opts)
	case "FSDRIVE":
		vd = fsdrive.New(c, suffix, deviceId, deviceNumber, opts)
	default:
		vd = c1541board.New(c, suffix, deviceId, deviceNumber, opts)
	}
	return vd
}
