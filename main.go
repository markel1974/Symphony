package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/c64emu/src/preferences"
	"github.com/markel1974/c64emu/src/render"
	"github.com/markel1974/c64emu/src/version"
)

func main() {
	var showHelp bool
	var showVersion bool
	var cartridge string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&cartridge, "c", "", "cartridge path")
	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println(version.AppName, version.AppVersion)
		return
	}

	prefs := preferences.NewPrefs()
	if len(cartridge) > 0 {
		prefs.SetCartridge(cartridge)
	}
	g := render.New(prefs)
	g.Start()
}
