package main

import (
	"flag"
	"fmt"
	symphony "github.com/markel1974/c64emu/src"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/renderers/audio"
	"github.com/markel1974/c64emu/src/renderers/graphics"
	"github.com/markel1974/c64emu/src/shell/cli"
	"github.com/markel1974/c64emu/src/shell/interfaces"
	"github.com/markel1974/c64emu/src/shell/servers"
	"github.com/markel1974/c64emu/src/shell/servers/authenticator"
	"github.com/markel1974/c64emu/src/version"
	"log"
	"os"
	"path"
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

// -f "/Users/tinmr305/Downloads/c64carts/edge_of_disgrace-booze_design_0.d64"

// -f "fs_drive:/Users/tinmr305/Downloads/c64carts/"

// -c "/Users/tinmr305/Downloads/c64carts/1541DiagnosticCart/1541diagcart.crt" -f "/Users/tinmr305/Downloads/c64carts/hessian-xth.d64"

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

	run := func(task interfaces.ITask, args []string) error {
		return nil
	}
	t := cli.NewCommand("bin", interfaces.CommandTypeDirectory, nil, false, run)
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
	k := servers.NewServer(secure, auth, port, false)
	k.SetPrompt(prompt)
	k.SetTemplate(t)
	go func() {
		k.Start()
	}()
	return nil
}

func BuildPrg(prgFile string) ([]byte, error) {
	if len(prgFile) == 0 {
		return []byte{}, nil
	}
	src, err := os.ReadFile(prgFile)
	if err != nil {
		return nil, err
	}
	return src, nil
}

func BuildCartridges(c string) ([]*config.Cartridge, error) {
	if len(c) == 0 {
		return nil, nil
	}
	var cartridges []*config.Cartridge
	for _, v := range config.KeyVal(c) {
		data, _, err := config.ImageFromFile(v.V)
		if err != nil {
			return nil, err
		}
		name := path.Base(v.V)
		cartridge, err := config.NewCartridge(v.K, v.V, name, data)
		if err != nil {
			return nil, err
		}
		cartridges = append(cartridges, cartridge)
	}
	return cartridges, nil
}

// BuildDrives parses a drive configuration string, creates Drive instances, and appends them to the Config's drives list.
// Returns an error if any drive creation fails.
func BuildDrives(d string) ([]*config.Drive, error) {
	if len(d) == 0 {
		return nil, nil
	}
	var drives []*config.Drive
	for _, v := range config.KeyVal(d) {
		data, wp, err := config.ImageFromFile(v.V)
		if err != nil {
			return nil, err
		}
		drive, err := config.NewDrive(v.K, v.V, data, wp)
		if err != nil {
			return nil, err
		}
		drives = append(drives, drive)
	}
	return drives, nil
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
	var playerId string
	var renderId string

	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&cartridges, "c", "", "cartridge path")
	flag.StringVar(&drives, "d", "", "drives path")
	flag.StringVar(&disks, "f", "", "disks")
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

	opt := symphony.NewOptions(graphics.NewFactory(), audio.NewFactory())
	cCartridges, err := BuildCartridges(cartridges)
	if err != nil {
		log.Fatal(err)
	}
	cDrives, err := BuildDrives(drives)
	if err != nil {
		log.Fatal(err)
	}
	cDisks, err := BuildDrives(disks)
	if err != nil {
		log.Fatal(err)
	}
	cPrg, err := BuildPrg(prg)
	if err != nil {
		log.Fatal(err)
	}
	opt.Prg = cPrg
	opt.Cartridges = cCartridges
	opt.Drives = cDrives
	opt.Disks = cDisks
	opt.NoJiffy = noJiffy
	opt.PlayerId = playerId
	opt.RenderId = renderId

	emulator := symphony.New()
	if err = emulator.Setup(opt); err != nil {
		log.Fatal(err)
	}
	board := emulator.GetBoard()
	if err = createShell(board.GetCommand()); err != nil {
		log.Fatal(err)
	}
	if err = emulator.Start(); err != nil {
		log.Fatal(err)
	}
}
