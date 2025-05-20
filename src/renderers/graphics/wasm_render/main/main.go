package main

import (
	symphony "github.com/markel1974/c64emu/src"
	"log"
)

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
