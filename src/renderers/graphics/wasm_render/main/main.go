package main

import (
	symphony "github.com/markel1974/c64emu/src"
	"log"
)

//cp $GOROOT/misc/wasm/wasm_exec.js static_content
//GOOS=js GOARCH=wasm go build -o  ./src/renderers/graphics/wasm_render/server/static_content/symphony.wasm ./src/renderers/graphics/wasm_render/main

func main() {
	gFactory := NewGraphicsFactory()
	aFactory := NewAudioFactory()

	opt := symphony.NewOptions(gFactory, aFactory)
	emulator := symphony.New()
	if err := emulator.Setup(opt); err != nil {
		log.Fatal(err)
	}
	if err := emulator.Start(); err != nil {
		log.Fatal(err)
	}
}
