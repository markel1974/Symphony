package main

import (
	"flag"
	"fmt"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/hardware"
	"github.com/markel1974/c64emu/src/references"
	"github.com/markel1974/c64emu/src/renderers/audio"
	"github.com/markel1974/c64emu/src/renderers/graphics"
	"github.com/markel1974/c64emu/src/shell"
	"github.com/markel1974/c64emu/src/shell/authenticator"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"github.com/markel1974/c64emu/src/version"
	"log"
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

var _c64DefaultHardware = []struct {
	label    string
	id       string
	instance int
}{
	{"c64", "c64", 0},
	{"c64", "iec", 0},
	{"c64", "mos6569", 0},
	{"c64", "mos6526", 0},
	{"c64", "mos6526", 1},
	{"c64", "cartridges_c64", 0},
	{"c64", "mos6510", 0},
	{"c64", "pic_6510", 0},
	{"c64", "dynamic_throttle", 0},
	{"c64", "mos6581", 0},
	{"c64", "roms_c64", 0},
	{"c64", "pla_c64", 0},
	{"c64", "keyboard_c64", 0},
	{"c64", "joystick_c64", 0},
	{"c64", "joystick_c64", 1},
	{"c64", "quartz", 0},

	/*
		{"c1541_8", "roms_c1541", 0},
		{"c1541_8", "pla_c1541", 0},
		{"c1541_8", "pic_6510", 0},
		{"c1541_8", "mos6510", 0},
		{"c1541_8", "mos6522", 0},
		{"c1541_8", "mos6522", 1},

		{"c1541_9", "roms_c1541", 0},
		{"c1541_9", "pla_c1541", 0},
		{"c1541_9", "pic_6510", 0},
		{"c1541_9", "mos6510", 0},
		{"c1541_9", "mos6522", 0},
		{"c1541_9", "mos6522", 1},

		{"c1541_10", "roms_c1541", 0},
		{"c1541_10", "pla_c1541", 0},
		{"c1541_10", "pic_6510", 0},
		{"c1541_10", "mos6510", 0},
		{"c1541_10", "mos6522", 0},
		{"c1541_10", "mos6522", 1},

		{"c1541_11", "roms_c1541", 0},
		{"c1541_11", "pla_c1541", 0},
		{"c1541_11", "pic_6510", 0},
		{"c1541_11", "mos6510", 0},
		{"c1541_11", "mos6522", 0},
		{"c1541_11", "mos6522", 1},

	*/
}

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

	run := func(r interfaces.IContext, cmd *cli.Command, pid int, args []string) error {
		return nil
	}
	t := cli.NewCommand("bin", nil, false, run)
	t.SetHelp("Bin", "Bin")
	_ = t.AddCommand(target)
	auth := authenticator.NewSimpleAuthenticator()
	if err := auth.Setup(user, pass); err != nil {
		return err
	}
	fmt.Println("Starting shell")
	fmt.Println("port", port)
	fmt.Println("secure", secure)
	fmt.Println("user", user)
	k := shell.NewServer(secure, auth, port, false)
	k.SetPrompt(prompt)
	k.SetTemplate(t)
	go func() {
		k.Start()
	}()
	return nil
}

func main() {
	//https://sergetoro.com/posts/golang-round-robin-queue-from-scratch/
	//fifoBlockTest()
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

	graphicsFactory := graphics.NewFactory()
	audioFactory := audio.NewFactory()
	displayRender := graphicsFactory.Create(renderId)
	display, err := displayRender.CreateDisplayBuffer(0x180, 0x110) //mos6569.DisplayX, mos6569.DisplayY
	if err != nil {
		log.Fatal(err)
	}
	audioRender := audioFactory.Create(playerId)
	if err = audioRender.Setup(cfg); err != nil {
		log.Fatal(err)
	}

	hwFactory := hardware.NewFactory(display, audioRender, cfg)

	var boardComponent references.IComponent = nil
	var hw []references.IComponent
	components := make(map[string]references.IComponent)
	for _, h := range _c64DefaultHardware {
		var comp references.IComponent
		if comp, err = hwFactory.Create(boardComponent, h.label, h.id, h.instance); err != nil {
			log.Fatal(err.Error())
		}
		components[comp.HardwareId()] = comp
		if boardComponent == nil {
			boardComponent = comp
		}
		hw = append(hw, comp)
	}

	if boardComponent == nil {
		log.Fatal("board is nil")
	}
	board, ok := boardComponent.(references.IBoard)
	if !ok {
		log.Fatal("invalid board")
	}

	//TODO REMOVE WHEN TREE IS READY
	//BEGIN
	for _, c := range hw {
		if err = c.Setup(); err != nil {
			log.Fatal(err)
		}
	}
	for _, c := range hw {
		if err = c.Connect(); err != nil {
			log.Fatal(err)
		}
	}
	//END

	if err = createShell(boardComponent.GetCommand()); err != nil {
		log.Fatal(err)
	}
	if err = displayRender.Setup(board, cfg); err != nil {
		log.Fatal(err)
	}
	if err = board.Start(); err != nil {
		log.Fatal(err)
	}
	if err = displayRender.Start(); err != nil {
		log.Fatal(err)
	}
}
