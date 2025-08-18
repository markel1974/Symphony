package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/markel1974/c64emu/src"
	"github.com/markel1974/c64emu/src/config"
	"github.com/markel1974/c64emu/src/kernel/compiler"
	"github.com/markel1974/c64emu/src/kernel/compiler/sdk"
	"github.com/markel1974/c64emu/src/kernel/compiler/stub"
	"github.com/markel1974/c64emu/src/kernel/component"
	"github.com/markel1974/c64emu/src/kernel/frontend"
	"github.com/markel1974/c64emu/src/kernel/frontend/authenticator"
	"github.com/markel1974/c64emu/src/kernel/interfaces"
	"github.com/markel1974/c64emu/src/kernel/process"
	"github.com/markel1974/c64emu/src/kernel/vm"
	"github.com/markel1974/c64emu/src/kernel/vm/bytecode"
	"github.com/markel1974/c64emu/src/renderers/audio"
	"github.com/markel1974/c64emu/src/renderers/graphics"
	"github.com/markel1974/c64emu/src/version"
)

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

func createShell(target interfaces.ICommand) error {
	const prompt = "symphony" + " " + "1.4.3" + "> "
	const port = 1234
	const user = "u"
	const secure = true
	const pass = "p"

	run := func(proc interfaces.IUserProcess, args []string) error {
		return nil
	}
	t := process.NewCommand("bin", interfaces.CommandTypeDirectory, nil, false, run)
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
	k := frontend.NewFrontend(secure, auth, port, false)
	k.SetPrompt(" % ")
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

func vmTest() {
	comp := compiler.New()
	bc, err := comp.Compile("example.go", stub.Source5)
	if err != nil {
		log.Fatalf("compiler error: %s", err)
	}

	d := bytecode.NewDisassembler(bc)
	d.Disassemble(log.Writer())

	loader := sdk.NewLoader()
	machine, err := vm.NewVM(loader, nil, bc, 100000)
	if err != nil {
		log.Fatalf("VM creation error: %s", err)
	}
	if err = machine.Run("main", 1, 2); err != nil {
		machine.Print(log.Writer())
		log.Fatalf("VM runtime error: %s", err)
	}
	machine.Print(log.Writer())
	os.Exit(0)
}

type NilWriter struct {
}

func (nw *NilWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func main() {
	//benchmark.VIC(1000000, 20, 10, 1)

	vmTest()

	var showHelp bool
	var showVersion bool
	var cartridges string
	var drives string
	var disks string
	var prg string
	var noJiffy bool
	var playerId string
	var renderId string
	var noShell bool
	var reflector string

	flag.BoolVar(&showHelp, "h", false, "show this help")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.StringVar(&cartridges, "c", "", "cartridge path")
	flag.StringVar(&drives, "d", "", "drives path")
	flag.StringVar(&disks, "f", "", "disks")
	flag.StringVar(&renderId, "r", "gl", "graphics: gl, ascii")
	flag.StringVar(&playerId, "a", "default", "audio: default")
	flag.StringVar(&prg, "p", "", "prg path")
	flag.BoolVar(&noJiffy, "j", false, "disable jiffy")
	flag.BoolVar(&noShell, "k", false, "disable shell")
	flag.StringVar(&reflector, "g", "", "generate reflect file")
	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println(version.AppName, version.AppVersion)
		return
	}

	if len(reflector) > 0 {
		for _, v := range strings.Split(reflector, ",") {
			gen := component.NewGenerator(v, true)
			//gen.SetOutput(os.Stdout)
			if err := gen.ParseAndGenerate(); err != nil {
				log.Fatal(err)
			}
		}
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
	if !noShell {
		if err = createShell(board.GetCommand()); err != nil {
			log.Fatal(err)
		}
	}
	if err = emulator.Start(); err != nil {
		log.Fatal(err)
	}
}
