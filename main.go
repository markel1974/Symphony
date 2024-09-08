package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/render"
	"github.com/markel1974/c64emu/src/version"
)

// -c "SCPU:;REU16M:/Users/tinmr305/Downloads/c64carts/doom/doom.reu" -p /Users/tinmr305/Downloads/c64carts/doom/loader.prg

// -d "/Users/tinmr305/Downloads/c64carts/C64_disk/blast170.d64" -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Easyflash/d2ef-vol03.crt"

// -d "/Users/tinmr305/Downloads/c64carts/GEOS64/GEOS64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/WRUTIL64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/APPS64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/SPELL64.D64" -c "REU2M:"

// -d "/Users/tinmr305/Downloads/c64carts/GEOS64/GEOS64.D64" -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Games/2/Lemmings_ef.crt"

// -c "/Users/tinmr305/Downloads/c64carts/mayhem.crt"

// -c "/Users/tinmr305/Downloads/c64carts/batman.bin"

// -c "/Users/tinmr305/Downloads/c64carts/thepit.bin"

// -c "/Users/tinmr305/Downloads/c64carts/popeye.bin"

// -c "/Users/tinmr305/Downloads/c64carts/SuperWonderboyInMonsterland_1989_Activision-EF.crt"

// -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Games/1/defender_of_the_crown.crt"

// -c "/Users/tinmr305/Downloads/c64carts/GhostNGoblinsArcade_N0S_EF.crt"

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
	cfg := config.New()
	if len(prg) > 0 {
		cfg.SetPrg(prg)
	}
	if len(cartridges) > 0 {
		kv := config.KeyVal(cartridges)
		for _, v := range kv {
			cfg.AddCartridge(v.K, v.V)
		}
	}
	if len(drives) > 0 {
		kv := config.KeyVal(drives)
		for _, v := range kv {
			cfg.AddDrive(v.K, v.V)
		}
	}
	g := render.New(cfg)
	g.Start()
}
