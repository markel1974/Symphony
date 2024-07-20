package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/c64emu/src/preferences"
	"github.com/markel1974/c64emu/src/render"
	"github.com/markel1974/c64emu/src/version"
	"strings"
)

func test(a int, b int, c string, d int, e bool) {
	fmt.Println("INSIDE", a, b, c, d, e)
}
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
		for _, c := range strings.Split(cartridge, ";") {
			prefs.AddCartridge(c)
		}
	}
	g := render.New(prefs)
	g.Start()
}
