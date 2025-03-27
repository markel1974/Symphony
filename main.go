package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/renderers/audio"
	"github.com/markel1974/c64emu/src/renderers/graphics"
	"github.com/markel1974/c64emu/src/shell"
	"github.com/markel1974/c64emu/src/shell/authenticator"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/version"
)

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

// -f "fs_drive:/Users/tinmr305/Downloads/c64carts/"

/*
func restoreTest(factory references.IComponentFactory) {
	//state, _ := s.DumpAll()
	//buf, _ := json.MarshalIndent(state, "", " ")
	//fmt.Println(string(buf))
	// err := s.RestoreAll(state); err != nil {
	//	fmt.Println(err)
	//}

	stub := make(map[string]interface{})
	err := json.Unmarshal([]byte(treeStub), &stub)
	if err != nil {
		panic(err)
	}
	out, err := component.Restore(factory, nil, nil, stub)
	if err != nil {
		panic(err)
	}
	for _, v := range out.GetChildren() {
		fmt.Println(v.GetId())
	}
	os.Exit(1)
}
*/

func createShell(target *cli.Command) error {
	const prompt = "symphony" + " " + "1.4.3" + "> "
	const port = 1234
	const user = "u"
	const secure = true
	const pass = "p"

	t := cli.NewCommand()
	_ = t.AddCommand(target)
	auth := authenticator.NewSimpleAuthenticator()
	if err := auth.Setup(user, pass); err != nil {
		return err
	}
	fmt.Println("Starting shell")
	fmt.Println("port", port)
	fmt.Println("secure", secure)
	fmt.Println("user", user)
	k := shell.New(secure, auth, port, false)
	k.SetPrompt(prompt)
	k.SetTemplate(t)
	go func() {
		k.Start()
	}()
	return nil
}

func main() {
	//TestCommand()
	var showHelp bool
	var showVersion bool
	var cartridges string
	var drives string
	var disks string
	var prg string
	var noJiffy bool
	var boardId string
	var playerId string
	var renderId string

	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&cartridges, "c", "", "cartridge path")
	flag.StringVar(&drives, "d", "", "drives path")
	flag.StringVar(&disks, "f", "", "disks")
	flag.StringVar(&boardId, "m", "c64", "hardware: vic20, c64")
	flag.StringVar(&renderId, "r", "gl", "graphics: gl, ascii")
	flag.StringVar(&playerId, "a", "default", "audio: default")
	flag.StringVar(&prg, "p", "", "prg path")
	flag.BoolVar(&noJiffy, "j", false, "disable jiffy")
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
		if err := cfg.BuildCartridges(cartridges); err != nil {
			log.Fatal(err)
		}
	}
	if len(drives) > 0 {
		if err := cfg.BuildDrives(drives); err != nil {
			log.Fatal(err)
		}
	}
	if len(disks) > 0 {
		if err := cfg.BuildSpareDisks(disks); err != nil {
			log.Fatal(err)
		}
	}
	if noJiffy {
		cfg.DisableJiffy()
	}

	hwFactory := hardware.NewFactory(cfg)
	boardComponent, err := hwFactory.Create(nil, boardId, 0)
	if err != nil {
		log.Fatal(err)
	}
	board, ok := boardComponent.(references.IBoard)
	if !ok || board == nil {
		log.Fatal("board is nil")
	}

	graphicsFactory := graphics.NewFactory()
	audioFactory := audio.NewFactory()

	displayRender := graphicsFactory.Create(renderId)
	w, h := board.GetScreenSize()
	display, err := displayRender.CreateDisplayBuffer(w, h)
	if err != nil {
		log.Fatal(err)
	}
	audioRender := audioFactory.Create(playerId)
	if err = audioRender.Setup(cfg); err != nil {
		log.Fatal(err)
	}

	if err = board.Setup(display, audioRender, cfg); err != nil {
		log.Fatal(err)
	}
	if err = boardComponent.Connect(); err != nil {
		log.Fatal(err)
	}
	if err = createShell(boardComponent.GetCommand()); err != nil {
		log.Fatal(err)
	}
	if err = displayRender.Setup(board, cfg); err != nil {
		log.Fatal(err)
	}
	if err = displayRender.Start(); err != nil {
		log.Fatal(err)
	}
}
