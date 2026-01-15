package config

import (
	"fmt"
	"github.com/markel1974/symphony/src/common/signals"
	"io"
	"io/fs"
	"os"
)

type IAssetInfo interface {
	fs.FileInfo
}

type IAssetReader interface {
	io.ReadCloser
	io.ReaderAt
}

type IAssetWriter interface {
	io.WriteCloser
}

// Config represents a configuration structure for managing cartridges, drives, disks, and various related options.
type Config struct {
	cartridges []*Cartridge
	drives     []*Drive
	spareDisks []*Drive
	changed    *signals.Signal
	prg        []byte
	diskIndex  int
	jiffy      bool
}

// New initializes and returns a Config instance with default values.
func New() *Config {
	c := &Config{
		cartridges: []*Cartridge{},
		drives:     []*Drive{},
		spareDisks: []*Drive{},
		diskIndex:  0,
		changed:    signals.NewSignal(),
		jiffy:      true,
	}
	return c
}

// Bind associates the provided function with the Config's signal, triggering it whenever the signal is emitted.
func (p *Config) Bind(changed func()) {
	p.changed.Bind(changed)
}

// AddDrive adds a new Drive instance to the Config's drives list and returns an error if any issues occur during the process.
func (p *Config) AddDrive(drive *Drive) error {
	p.drives = append(p.drives, drive)
	return nil
}

// AddSpareDisk adds the specified Drive to the spareDisks list in the Config structure and updates drives if empty.
func (p *Config) AddSpareDisk(drive *Drive) error {
	if len(p.drives) == 0 {
		p.drives = append(p.spareDisks, drive)
	}
	p.spareDisks = append(p.spareDisks, drive)
	return nil
}

// Drives returns the list of drives configured in the Config structure.
func (p *Config) Drives() []*Drive {
	return p.drives
}

// Drive retrieves the Drive instance corresponding to the given id from the Config's drive list. Returns nil if not found.
func (p *Config) Drive(id uint8) *Drive {
	if int(id) < len(p.drives) {
		return p.drives[id]
	}
	return nil
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
func (p *Config) SetPrg(prg []byte) {
	p.prg = prg
}

// Prg returns the value of the `prg` field from the Config struct.
func (p *Config) Prg() []byte {
	return p.prg
}

func (p *Config) AddCartridge(crt *Cartridge) error {
	p.cartridges = append(p.cartridges, crt)
	return nil
}

func (p *Config) AssetInfo(asset string) (IAssetInfo, error) {
	//default implementation
	s, err := os.Stat(asset)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (p *Config) AssetWrite(asset string, data []byte) error {
	//default implementation
	if err := os.WriteFile(asset, data, 0644); err != nil {
		return err
	}
	return nil
}

func (p *Config) AssetRead(asset string) ([]byte, error) {
	//default implementation
	dat, err := os.ReadFile(asset)
	return dat, err
}

func (p *Config) AssetReader(asset string) (IAssetReader, int64, error) {
	f, err := os.Open(asset)
	if err != nil {
		return nil, 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, err
	}
	return f, fi.Size(), nil
}

func (p *Config) AssetDir(asset string) ([]IAssetInfo, error) {
	items, e := os.ReadDir(asset)
	if e != nil {
		return nil, e
	}
	var out []IAssetInfo
	for _, item := range items {
		info, err := item.Info()
		if err != nil {
			return nil, err
		}
		out = append(out, info)
	}
	return out, nil
}

func (p *Config) AssetPathSeparator() string {
	return string(os.PathSeparator)
}

func (p *Config) AssetWriter(asset string) (IAssetWriter, error) {
	f, err := os.OpenFile(asset, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Cartridges returns the list of configured cartridges in the Config structure.
func (p *Config) Cartridges() []*Cartridge {
	return p.cartridges
}

// C1541RomAsset returns the file path of the 1541 ROM as a string.
func (p *Config) C1541RomAsset() string {
	return ""
}

// C64RomKernalAsset returns the file path of the Kernal ROM as a string.
func (p *Config) C64RomKernalAsset() string {
	return ""
}

// C64RomBasicAsset returns the file path of the Kernal ROM as a string.
func (p *Config) C64RomBasicAsset() string {
	return ""
}

// C64RomCharAsset returns the file path of the Kernal ROM as a string.
func (p *Config) C64RomCharAsset() string {
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
