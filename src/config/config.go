package config

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/signals"
)

// Config represents a configuration structure for managing cartridges, drives, disks, and various related options.
type Config struct {
	cartridges []*Cartridge
	drives     []*Drive
	spareDisks []*Drive
	changed    *signals.Signal
	prg        string
	diskIndex  int
	jiffy      bool
}

// New initializes and returns a Config instance with default values.
func New() *Config {
	return &Config{
		cartridges: nil,
		diskIndex:  0,
		changed:    signals.NewSignal(),
		jiffy:      true,
	}
}

// Bind associates the provided function with the Config's signal, triggering it whenever the signal is emitted.
func (p *Config) Bind(changed func()) {
	p.changed.Bind(changed)
}

// AddDrive adds a new drive to the Config instance with the specified kind and options.
func (p *Config) AddDrive(kind string, filePath string) error {
	drive, err := NewDrive(kind, filePath)
	if err != nil {
		return err
	}
	p.drives = append(p.drives, drive)
	return nil
}

// AddDisk appends a new disk of the specified kind and options to the Config's disks list.
func (p *Config) AddDisk(kind string, filePath string) error {
	drive, err := NewDrive(kind, filePath)
	if err != nil {
		return err
	}
	p.spareDisks = append(p.spareDisks, drive)
	return nil
}

// GetDrives returns the list of drives configured in the Config structure.
func (p *Config) GetDrives() []*Drive {
	return p.drives
}

// GetDrive retrieves the options of a drive by its ID if it exists. Returns the options and true if found, else an empty string and false.
func (p *Config) GetDrive(id uint8) (*Drive, bool) {
	if int(id) < len(p.drives) {
		return p.drives[id], true
	}
	return nil, false
}

// SwitchDisk cycles through the available disks, updates the active drive's options, and emits a configuration change signal.
func (p *Config) SwitchDisk() (string, error) {
	if len(p.drives) == 0 || len(p.spareDisks) == 0 {
		return "", fmt.Errorf("nil disk")
	}
	p.diskIndex++
	driveIndex := p.diskIndex % len(p.spareDisks)
	p.drives[0] = p.spareDisks[driveIndex]
	p.changed.Emit()
	return p.spareDisks[driveIndex].GetId(), nil
}

// SetPrg sets the program path in the Config instance.
func (p *Config) SetPrg(prg string) {
	p.prg = prg
}

// GetPrg returns the value of the `prg` field from the Config struct.
func (p *Config) GetPrg() string {
	return p.prg
}

// AddCartridge appends a new Cartridge with the specified kind and path to the cartridges slice of the Config instance.
func (p *Config) AddCartridge(kind string, filePath string) error {
	crt, err := NewCartridge(kind, filePath)
	if err != nil {
		return err
	}
	p.cartridges = append(p.cartridges, crt)
	return nil
}

// GetCartridges returns the list of configured cartridges in the Config structure.
func (p *Config) GetCartridges() []*Cartridge {
	return p.cartridges
}

// Get1541RomPath returns the file path of the 1541 ROM as a string.
func (p *Config) Get1541RomPath() string {
	return ""
}

// GetKernalRomPath returns the file path of the Kernal ROM as a string.
func (p *Config) GetKernalRomPath() string {
	return ""
}

// DisableJiffy disables the Jiffy mode by setting the `jiffy` field to `false`.
func (p *Config) DisableJiffy() {
	p.jiffy = false
}

// UseJiffy returns the current state of the jiffy mode in the configuration.
func (p *Config) UseJiffy() bool {
	return p.jiffy
}
