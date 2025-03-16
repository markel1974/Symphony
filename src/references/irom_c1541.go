package references

import "github.com/markel1974/c64emu/src/config"

type IROMLoaderC1541 interface {
	Setup(cfg *config.Config) error

	Load() []byte
}
