//go:build js && wasm

package main

import (
	"syscall/js"

	symphony "github.com/markel1974/c64emu/src"
	"github.com/markel1974/c64emu/src/config"
	"log"
)

//cp $GOROOT/misc/wasm/wasm_exec.js in static_content
//cp index.html in static_content
//GOOS=js GOARCH=wasm go build -o  ./src/renderers/graphics/wasm_render/server/static_content/symphony.wasm ./src/renderers/graphics/wasm_render/main

func stubCartridges() ([]*config.Cartridge, error) {
	cart, err := config.NewCartridge("", "mayhem.crt", "mayhem.crt", _stub)
	if err != nil {
		return nil, err
	}
	return []*config.Cartridge{cart}, nil
}

func main() {
	var cCarts []*config.Cartridge = nil
	var err error
	gFactory := NewGraphicsFactory()
	aFactory := NewAudioFactory()

	opt := symphony.NewOptions(gFactory, aFactory)
	symphonyConfig := js.Global().Get("symphonyConfig")
	if symphonyConfig.IsUndefined() {
		log.Println("undefined symphonyConfig")
		return
	}
	if symphonyConfig.Type() != js.TypeObject {
		log.Println("invalid symphonyConfig type")
		return
	}

	var buffer []byte
	var name string
	var cart *config.Cartridge
	jsUseStub := symphonyConfig.Get("useStub")
	if !jsUseStub.IsUndefined() && jsUseStub.Type() == js.TypeBoolean && jsUseStub.Bool() {
		if cCarts, err = stubCartridges(); err != nil {
			log.Println(err)
			return
		}
	} else {
		jsCartName := symphonyConfig.Get("cartridgeName")
		jsCartBuffer := symphonyConfig.Get("cartridgeBuffer")
		if !jsCartName.IsUndefined() && jsCartName.Type() == js.TypeString {
			name = jsCartName.String()
			buffer, err = ConvertJsBufferToGoBytes(jsCartBuffer)
			if err != nil {
				log.Println(err)
				return
			}
			cart, err = config.NewCartridge("", name, name, buffer)
			if err != nil {
				log.Println(err)
				return
			}
			cCarts = []*config.Cartridge{cart}
		}
	}

	opt.Cartridges = cCarts
	emulator := symphony.New()
	if err = emulator.Setup(opt); err != nil {
		log.Println(err)
		return
	}
	if err = emulator.Start(); err != nil {
		log.Println(err)
		return
	}
}
