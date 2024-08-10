package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/render"
	"github.com/markel1974/c64emu/src/version"
	"strings"
)

//-c "SCPU:;REU16M:/Users/tinmr305/Downloads/c64carts/doom/doom.reu" -p /Users/tinmr305/Downloads/c64carts/doom/loader.prg

func main() {
	var showHelp bool
	var showVersion bool
	var cartridge string
	var prg string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&cartridge, "c", "", "cartridge path")
	flag.StringVar(&prg, "p", "", "prg path")
	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println(version.AppName, version.AppVersion)
		return
	}
	prefs := config.New()
	if len(prg) > 0 {
		prefs.SetPrg(prg)
	}
	if len(cartridge) > 0 {
		for _, c := range strings.Split(cartridge, ";") {
			kind := ""
			path := c
			if opts := strings.Split(c, ":"); len(opts) > 1 {
				kind = opts[0]
				path = opts[1]
			}
			prefs.AddCartridge(kind, path)
		}
	}
	g := render.New(prefs)
	g.Start()
}
