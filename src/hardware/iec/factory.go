package iec

import (
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/config"
	c1541board "github.com/markel1974/c64emu/src/hardware/c1541/board"
	"github.com/markel1974/c64emu/src/hardware/iec/fsdrive"
	"github.com/markel1974/c64emu/src/references"
)

// DriveFactory represents a component responsible for creating and managing drive configurations and operations.
type DriveFactory struct {
	*component.BaseComponent
	factory references.IComponentFactory
	cfg     *config.Config
}

// NewDriveFactory creates and initializes a new DriveFactory instance, registers it with the provided parent component.
func NewDriveFactory(parent component.IComponent, factory references.IComponentFactory, suffix string) *DriveFactory {
	df := &DriveFactory{
		BaseComponent: component.NewBaseComponent("iec_drive_factory", suffix),
		factory:       factory,
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
func (c *DriveFactory) Create(kind string, suffix string) references.IIecDevice {
	var vd references.IIecDevice
	switch kind {
	case "C1541":
		vd = c1541board.New(c, c.factory, suffix)
	case "FSDRIVE":
		vd = fsdrive.New(c, c.factory, suffix)
	default:
		vd = c1541board.New(c, c.factory, suffix)
	}
	return vd
}
