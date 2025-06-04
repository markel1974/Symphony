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

	var cart *config.Cartridge

	useStub, _ := JsBooleanToGoBool(symphonyConfig.Get("useStub"))
	if useStub {
		if cCarts, err = stubCartridges(); err != nil {
			log.Println(err)
			return
		}
	} else {
		cartName, _ := JsStringToGoString(symphonyConfig.Get("cartridgeName"))
		if len(cartName) > 0 {
			cartBuffer, cErr := JsBufferToGoBytes(symphonyConfig.Get("cartridgeBuffer"))
			if cErr != nil {
				log.Println(cErr)
				return
			}
			cart, err = config.NewCartridge("", cartName, cartName, cartBuffer)
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
