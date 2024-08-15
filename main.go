package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/render"
	"github.com/markel1974/c64emu/src/version"
)

//-c "SCPU:;REU16M:/Users/tinmr305/Downloads/c64carts/doom/doom.reu" -p /Users/tinmr305/Downloads/c64carts/doom/loader.prg

//-d "/Users/tinmr305/Downloads/c64carts/C64_disk/blast170.d64" -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Easyflash/d2ef-vol50.crt"

//-d "/Users/tinmr305/Downloads/c64carts/GEOS64/GEOS64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/WRUTIL64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/APPS64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/SPELL64.D64"

func main() {
	var showHelp bool
	var showVersion bool
	var cartridges string
	var drives string
	var prg string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&cartridges, "c", "", "cartridge path")
	flag.StringVar(&drives, "d", "", "drives path")
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
	if len(cartridges) > 0 {
		kv := config.KeyVal(cartridges)
		for _, v := range kv {
			prefs.AddCartridge(v.K, v.V)
		}
	}
	if len(drives) > 0 {
		kv := config.KeyVal(drives)
		for _, v := range kv {
			prefs.AddDrive(v.K, v.V)
		}
	}

	g := render.New(prefs)
	g.Start()
}
