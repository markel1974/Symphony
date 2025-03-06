package main

import (
	"flag"
	"fmt"
	c64board "github.com/markel1974/c64emu/src/c64/board"
	"github.com/markel1974/c64emu/src/components/board"
	"github.com/markel1974/c64emu/src/components/sid"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/render/asciirender"
	"github.com/markel1974/c64emu/src/render/glrender"
	"github.com/markel1974/c64emu/src/version"
	vic20board "github.com/markel1974/c64emu/src/vic20/board"
	"os"
)

//TODO WASM
// https://garciat.com/posts/go-wasm/
// https://github.com/seqsense/webgl-go/tree/master

// -c "SCPU:;REU16M:/Users/tinmr305/Downloads/c64carts/doom/doom.reu" -p /Users/tinmr305/Downloads/c64carts/doom/loader.prg

// -d "/Users/tinmr305/Downloads/c64carts/C64_disk/blast170.d64" -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Easyflash/d2ef-vol03.crt"

// -d "/Users/tinmr305/Downloads/c64carts/GEOS64/GEOS64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/WRUTIL64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/APPS64.D64;/Users/tinmr305/Downloads/c64carts/GEOS64/SPELL64.D64" -c "REU2M:"

// -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Games/2/Lemmings_ef.crt"

// -c "/Users/tinmr305/Downloads/c64carts/mayhem.crt"

// -c "/Users/tinmr305/Downloads/c64carts/batman.bin"

// -c "/Users/tinmr305/Downloads/c64carts/thepit.bin"

// -c "/Users/tinmr305/Downloads/c64carts/popeye.bin"

// -c "/Users/tinmr305/Downloads/c64carts/frogger.bin"

// -c "/Users/tinmr305/Downloads/c64carts/SuperWonderboyInMonsterland_1989_Activision-EF.crt"

// -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Games/1/defender_of_the_crown.crt"

// -c "/Users/tinmr305/Downloads/c64carts/GhostNGoblinsArcade_N0S_EF.crt"

// -c "/Users/tinmr305/Downloads/c64carts/briley_witch_chronicles_v1.01_[excess].crt"

// -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Games/2/Pinball Spectacular.crt"

// -c "/Users/tinmr305/Downloads/c64carts/c64diag.bin"

// -c "/Users/tinmr305/Downloads/c64carts/Ultimate/Games/3/Tool Collection.crt"

// -f "/Users/tinmr305/Downloads/c64carts/Elvira2_CHR/Elvira2_Chr_A.d64;/Users/tinmr305/Downloads/c64carts/Elvira2_CHR/Elvira2_Chr_B.d64;/Users/tinmr305/Downloads/c64carts/Elvira2_CHR/Elvira2_Chr_C.d64;/Users/tinmr305/Downloads/c64carts/Elvira2_CHR/Elvira2_Chr_D.d64;/Users/tinmr305/Downloads/c64carts/Elvira2_CHR/Elvira2_Chr_E.d64"

// -f "/Users/tinmr305/Downloads/c64carts/mw4_1.d64;/Users/tinmr305/Downloads/c64carts/mw4_2.d64"

// -f "/Users/tinmr305/Downloads/c64carts/0_LOAD-Steel_Ranger_+3CD_LAXITY.d64;/Users/tinmr305/Downloads/c64carts/1_GAME-Steel_Ranger_+3CD_LAXITY.d64"

// -f "/Users/tinmr305/Downloads/c64carts/01SonicTheHedgehog_v1.2+5_-TRIAD+GP.d64;/Users/tinmr305/Downloads/c64carts/02SonicTheHedgehog_v1.2+5_-TRIAD+GP.d64" -c "REU2M:"

// -f "/Users/tinmr305/Downloads/c64carts/hessian-xth.d64"

// -f "/Users/tinmr305/Downloads/c64carts/0_LOAD-Steel_Ranger_+3CD_LAXITY.d64;/Users/tinmr305/Downloads/c64carts/1_GAME-Steel_Ranger_+3CD_LAXITY.d64" -c "/Users/tinmr305/Downloads/c64carts/1541DiagnosticCart/1541diagcart.crt"

// -c "/Users/tinmr305/Downloads/c64carts/wrathdemon.crt"

// -c "/Users/tinmr305/Downloads/c64carts/eye_of_beholder.crt"

// -f "/Users/tinmr305/Downloads/c64carts/SamsJourneySeasonsSpecialV1_1+5D-GP.d64"

func Test() {
	components := board.NewComponents()
	s := mos6581.NewSID("test", "")
	components.Register(s)
	components.Dump(s.GetId())
	os.Exit(1)
}

type IRender interface {
	Start() error
}

func main() {
	var showHelp bool
	var showVersion bool
	var cartridges string
	var drives string
	var disks string
	var prg string
	var noJiffy bool
	var ascii bool
	var mode string
	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&cartridges, "c", "", "cartridge path")
	flag.StringVar(&drives, "d", "", "drives path")
	flag.StringVar(&disks, "f", "", "disks")
	flag.StringVar(&mode, "m", "", "hardware mode: vic20, c64")
	flag.StringVar(&prg, "p", "", "prg path")
	flag.BoolVar(&noJiffy, "j", false, "disable jiffy")
	flag.BoolVar(&ascii, "a", false, "ascii render")
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
		for _, v := range config.KeyVal(drives) {
			cfg.AddDrive(v.K, v.V)
		}
	}
	if len(disks) > 0 {
		if kv := config.KeyVal(disks); len(kv) > 0 {
			if len(cfg.GetDrives()) == 0 {
				cfg.AddDrive(kv[0].K, kv[0].V)
			}
			for _, v := range kv {
				cfg.AddDisk(v.K, v.V)
			}
		}
	}
	if noJiffy {
		cfg.DisableJiffy()
	}

	var g IRender
	var b board.IBoard

	if mode == "vic20" {
		b = vic20board.NewBoard()
	} else {
		b = c64board.NewBoard()
	}

	if ascii {
		g = asciirender.New(b, cfg)
	} else {
		g = glrender.New(b, cfg)
	}

	if err := g.Start(); err != nil {
		fmt.Println(err)
		return
	}
}
