//go:build js && wasm

package main

import (
	"syscall/js"

	symphony "github.com/markel1974/c64emu/src"
	"github.com/markel1974/c64emu/src/config"
	"log"
)

//GO <= 1.20 cp $GOROOT/misc/wasm/wasm_exec.js src/renderers/graphics/wasm_render/server/static_content/
//GO  > 1.20 cp $GOROOT/lib/wasm/wasm_exec.js src/renderers/graphics/wasm_render/server/static_content/
//cp index.html in static_content
//GOOS=js GOARCH=wasm go build -o  ./src/renderers/graphics/wasm_render/server/static_content/symphony.wasm ./src/renderers/graphics/wasm_render/main

func main() {
	var cCarts []*config.Cartridge = nil

	symphonyConfig := js.Global().Get("symphonyConfig")
	if symphonyConfig.IsUndefined() {
		log.Println("undefined symphonyConfig")
		return
	}
	if symphonyConfig.Type() != js.TypeObject {
		log.Println("invalid symphonyConfig type")
		return
	}

	useStub, _ := JsBooleanToGoBool(symphonyConfig.Get("useStub"))
	if useStub {
		var err error
		cart, err := config.NewCartridge("", "mayhem.crt", "mayhem.crt", _stub)
		if err != nil {
			log.Println(err)
			return
		}
		cCarts = []*config.Cartridge{cart}
	} else {
		cartName, _ := JsStringToGoString(symphonyConfig.Get("cartridgeName"))
		if len(cartName) > 0 {
			cartBuffer, err := JsBufferToGoBytes(symphonyConfig.Get("cartridgeBuffer"))
			if err != nil {
				log.Println(err)
				return
			}
			var cart *config.Cartridge
			cart, err = config.NewCartridge("", cartName, cartName, cartBuffer)
			if err != nil {
				log.Println(err)
				return
			}
			cCarts = []*config.Cartridge{cart}
		}
	}

	gFactory := NewGraphicsFactory()
	aFactory := NewAudioFactory()

	opt := symphony.NewOptions(gFactory, aFactory)

	opt.Cartridges = cCarts
	emulator := symphony.New()
	if err := emulator.Setup(opt); err != nil {
		log.Println(err)
		return
	}
	if err := emulator.Start(); err != nil {
		log.Println(err)
		return
	}
}
