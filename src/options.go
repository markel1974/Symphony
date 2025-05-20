package symphony

import (
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/references"
)

type IGraphicFactory interface {
	Create(id string) references.IDisplayRender
}

type IAudioFactory interface {
	Create(id string) references.IAudioRender
}

type Options struct {
	RenderId        string
	PlayerId        string
	Prg             string
	Cartridges      []*config.Cartridge
	Drives          []*config.Drive
	Disks           []*config.Drive
	NoJiffy         bool
	graphicsFactory IGraphicFactory
	audioFactory    IAudioFactory
}

func NewOptions(gFactory IGraphicFactory, aFactory IAudioFactory) *Options {
	return &Options{
		graphicsFactory: gFactory,
		audioFactory:    aFactory,
	}
}

func (o *Options) GetGraphicsFactory() IGraphicFactory {
	return o.graphicsFactory
}

func (o *Options) GetAudioFactory() IAudioFactory {
	return o.audioFactory
}
