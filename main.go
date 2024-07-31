package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/c64emu/src/config"
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
	prefs := config.New()
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
