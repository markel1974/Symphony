package config

import "github.com/markel1974/c64emu/src/signals"

// Cartridge represents a cartridge with a specified kind and file path.
type Cartridge struct {
	Kind string
	Path string
}

// Drive represents a storage device with a specified type and configuration options.
// Kind specifies the type of the drive.
// Opts defines additional options or settings for the drive.
type Drive struct {
	Kind string
	Opts string
}

// Config represents a configuration structure for managing cartridges, drives, disks, and various related options.
type Config struct {
	cartridges []Cartridge
	drives     []Drive
	disks      []Drive
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
func (p *Config) AddDrive(kind string, opts string) {
	p.drives = append(p.drives, Drive{Kind: kind, Opts: opts})
}

// AddDisk appends a new disk of the specified kind and options to the Config's disks list.
func (p *Config) AddDisk(kind string, opts string) {
	p.disks = append(p.disks, Drive{Kind: kind, Opts: opts})
}

// GetDrives returns the list of drives configured in the Config structure.
func (p *Config) GetDrives() []Drive {
	return p.drives
}

// GetDrivesOpt retrieves the options of a drive by its ID if it exists. Returns the options and true if found, else an empty string and false.
func (p *Config) GetDrivesOpt(id uint8) (string, bool) {
	if int(id) < len(p.drives) {
		return p.drives[id].Opts, true
	}
	return "", false
}

// SetDriveOpt updates the options for a specific drive identified by id and emits a signal if the update is successful.
// Returns true if the update is performed; otherwise, it returns false. Only applicable if id is within drive array bounds.
func (p *Config) SetDriveOpt(opt string, id uint8) bool {
	if int(id) < len(p.drives) {
		p.drives[id].Opts = opt
		p.changed.Emit()
		return true
	}
	return false
}

// SwitchDisk cycles through the available disks, updates the active drive's options, and emits a configuration change signal.
func (p *Config) SwitchDisk() {
	if len(p.disks) == 0 {
		return
	}
	p.diskIndex++
	driveIndex := p.diskIndex % len(p.disks)
	p.SetDriveOpt(p.disks[driveIndex].Opts, 0)
	p.changed.Emit()
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
func (p *Config) AddCartridge(kind string, path string) {
	p.cartridges = append(p.cartridges, Cartridge{Kind: kind, Path: path})
}

// GetCartridges returns the list of configured cartridges in the Config structure.
func (p *Config) GetCartridges() []Cartridge {
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
