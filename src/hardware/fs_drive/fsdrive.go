package fs_drive

import (
	"fmt"
	"github.com/markel1974/c64emu/src/common/fifo"
	"github.com/markel1974/c64emu/src/component"
	"github.com/markel1974/c64emu/src/hardware/iec"
	"github.com/markel1974/c64emu/src/references"
	"io"
	"os"
	"strings"

	"github.com/markel1974/c64emu/src/config"
)

type FSDrive struct {
	*component.BaseComponent
	references.IIecDevice
	protocol     *iec.Protocol
	commands     *Commands
	deviceId     uint8
	deviceNumber uint8
	path         string
	dirPath      string       // Path to directory
	origDirPath  string       // Original directory path
	dirTitle     string       // Directory title
	file         [16]*os.File // File pointers for each of the 16 channels
	readChar     [16]uint8    // Buffers for one-byte read-ahead
	ready        bool
	cfg          *config.Config
	buffer       [0xf][]byte
	test         *fifo.StaticFifo
}

func NewBoard(parent references.IComponent, factory references.IComponentFactory, label string, instance int) *FSDrive {
	protocol := iec.NewProtocol(factory, parent, label, instance)
	fs := &FSDrive{
		BaseComponent: component.NewBaseComponent(),
		IIecDevice:    protocol,
		protocol:      protocol,
		deviceId:      0,
		deviceNumber:  0,
		path:          "",
		commands:      NewCommands(),
		origDirPath:   "",
		cfg:           nil,
		buffer:        [0xf][]byte{},
		test:          fifo.NewStaticFifo(32),
	}
	fs.protocol.SetDevice(fs)
	fs.BaseComponent.Register(factory, fs.protocol, Identifier(), fs, references.IdIIecProtocolDevice(fs, label, 0))
	return fs
}

func New(parent references.IComponent, factory references.IComponentFactory, label string, instance int) references.IIecDevice {
	return NewBoard(parent, factory, label, instance)
}

func (v *FSDrive) Setup() error {
	v.cfg = v.GetFactory().GetConfig()
	v.cfg.Bind(v.configChanged)
	if err := v.protocol.Setup(); err != nil {
		return err
	}
	return nil
}

func (v *FSDrive) Bind(_ references.IIecDeviceSocket, deviceId uint8, deviceNumber uint8) error {
	v.deviceId = deviceId
	v.deviceNumber = deviceNumber
	v.origDirPath = v.path
	if err := v.protocol.Bind(v, deviceId, deviceNumber); err != nil {
		return err
	}
	if d := v.cfg.Drive(v.deviceId); d != nil {
		v.path = d.GetId()
	}
	if v.changeDirectory(v.origDirPath) {
		for i := 0; i < 16; i++ {
			v.file[i] = nil
		}
		v.ready = true
	}
	return nil
}

func (v *FSDrive) Connect() error {
	if err := v.protocol.Connect(); err != nil {
		return err
	}
	return nil
}

func (v *FSDrive) Internal() bool {
	return false
}

func (v *FSDrive) GetDeviceNumber() uint8 {
	return v.deviceNumber
}

func (v *FSDrive) Reset() {
	v.closeAllChannels()
	v.commands.CommandClear()
	v.commands.SetError(ERR_STARTUP)
	//TODO IN FASE DI RESET CAMBIARE LO STATO DEL BUS
	v.protocol.Reset()
}

func (v *FSDrive) LedTurnOn() {
}

func (v *FSDrive) LedTurnOff() {
}

func (v *FSDrive) GetPath() string {
	return v.dirPath
}

func (v *FSDrive) Listen(sec uint8) {
	fmt.Println("LISTEN", sec)
}

func (v *FSDrive) Unlisten(sec uint8) {
	data := v.openDirectory(sec, "", "PROVA")
	v.test = fifo.NewStaticFifo(uint(len(data)))
	for _, k := range data {
		v.test.Set(int(k))
	}
	fmt.Println("UNLISTEN", sec)
}

func (v *FSDrive) Talk(sec uint8) {
	fmt.Println("TALK", sec)
}

func (v *FSDrive) Untalk(sec uint8) {
	fmt.Println("UNTALK", sec)
}

func (v *FSDrive) Open(channel uint8) uint8 {
	//TODO initialize channel

	//for _, c := range "PROVA" {
	//	v.test.Set(int(c))
	//}

	v.buffer[channel] = []byte{}
	return StOk
	/*
		//TODO DATA
		var data []uint8
		v.LedTurnOn()
		// Channel 15: Execute file name as command
		if channel == 15 {
			action, ok := v.commands.CommandExec(data)
			if ok {
				if action == 1 {
					v.Reset()
				}
			}
			return StOk
		}
		// Close previous file if still open
		if v.file[channel] != nil {
			v.file[channel].Close()
			v.file[channel] = nil
		}
		if len(data) == 0 {
			v.commands.SetError(ERR_NOCHANNEL)
			return StOk
		}
		if data[0] == '#' {
			v.commands.SetError(ERR_NOCHANNEL)
			return StOk
		}
		if data[0] == '$' {
			return v.openDirectory(channel, string(data))
		}
		return v.openFile(channel, string(data))

	*/
}

func (v *FSDrive) Close(channel uint8) uint8 {
	v.LedTurnOff()
	if channel == 15 {
		v.closeAllChannels()
		return StOk
	}
	if v.file[channel] != nil {
		v.file[channel].Close()
		v.file[channel] = nil
	}
	return StOk
}

func (v *FSDrive) Read(channel uint8) (uint8, uint8) {
	d, ok := v.test.Next()
	if !ok {
		return 0, StReadTimeout
	}
	fmt.Println("Read:", string(byte(d)))
	if v.test.Len() == 0 {
		v.commands.SetError(ERR_OK)
		return uint8(d), StEof
	}
	return uint8(d), StOk

	return 0, StReadTimeout
	// Channel 15: Error channel
	if channel == 15 {
		data := v.commands.RetrieveError()
		if data != '\r' {
			return data, StOk
		}
		// End of message
		v.commands.SetError(ERR_OK)
		return data, StEof
	}

	if v.file[channel] == nil {
		return 0, StReadTimeout
	}

	// Read one byte
	data := v.readChar[channel]
	buffer := make([]uint8, 1)
	c, err := v.file[channel].Read(buffer)
	if err == io.EOF {
		return data, StEof
	}
	v.readChar[channel] = (uint8)(c)
	return data, StOk
}

func (v *FSDrive) Write(channel uint8, data uint8) uint8 {
	//TODO EOI, eoi bool
	eoi := false

	if channel == 15 {
		if !v.commands.CommandSet(data) {
			return StTimeout
		}
		if eoi {
			v.commands.CommandExecBuf()
		}
		return StOk
	}
	//if v.buffer[channel] == nil {
	//	v.commands.SetError(ERR_FILENOTOPEN)
	//	return StTimeout
	//}

	v.buffer[channel] = append(v.buffer[channel], data)
	return StOk

	/*
		//TODO EOI, eoi bool
		eoi := false

		fmt.Printf("fsdrive_received: %s\n", string(data))

		return StOk

		// Channel 15: Collect chars and execute command on EOI
		if channel == 15 {
			if !v.commands.CommandSet(data) {

				return StTimeout
			}
			if eoi {
				v.commands.CommandExecBuf()
			}
			return StOk
		}
		if v.file[channel] == nil {
			v.commands.SetError(ERR_FILENOTOPEN)
			return StTimeout
		}
		if _, err := v.file[channel].Write([]byte{data}); err == io.EOF {
			v.commands.SetError(ERR_WRITE25)
			return StTimeout
		}
		return StOk

	*/
}

func (v *FSDrive) initializeCmd() {
	//v.closeAllChannels()
}

func (v *FSDrive) validateCmd() {
}

func (v *FSDrive) changeDirectory(dirPath string) bool {
	d, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	if !d.IsDir() {
		return false
	}
	v.dirPath = dirPath
	v.dirTitle = v.dirPath
	if len(v.dirTitle) > 16 {
		v.dirTitle = v.dirTitle[:16]
	}
	return true
}

func (v *FSDrive) findFirstFile(pattern string) (string, bool) {
	items, _ := os.ReadDir(v.dirPath)
	matcher := NewMatcher()
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		if b, err := matcher.Match(pattern, item.Name()); err == nil && b {
			return item.Name(), false
		}
	}
	return "", false
}

func (v *FSDrive) closeAllChannels() {
	for i := uint8(0); i < 15; i++ {
		v.Close(i)
	}
	v.commands.CommandClear()
}

func (v *FSDrive) openFile(channel uint8, name string) uint8 {
	plainName, mode, kind, _ := ParseFileName(name, true)
	// Channel 0 is READ, channel 1 is WRITE
	if channel == 0 || channel == 1 {
		mode = FMODE_READ
		if channel != 0 {
			mode = FMODE_WRITE
		}
		if kind == FTYPE_DEL {
			kind = FTYPE_PRG
		}
	}
	writing := mode == FMODE_WRITE || mode == FMODE_APPEND
	if strings.Contains(plainName, "*") || strings.Contains(plainName, "?") {
		if writing {
			v.commands.SetError(ERR_SYNTAX33)
			return StOk
		} else {
			v.findFirstFile(plainName)
		}
	}
	if kind == FTYPE_REL {
		v.commands.SetError(ERR_UNIMPLEMENTED)
		return StOk
	}
	flags := os.O_RDONLY
	perm := os.FileMode(0)
	switch mode {
	case FMODE_WRITE:
		perm = os.FileMode(0666)
		flags = os.O_RDWR
	case FMODE_APPEND:
		perm = os.FileMode(0666)
		flags = os.O_RDWR | os.O_APPEND
	default:
		panic("unhandled default case")
	}
	completeFileName := v.dirPath + string(os.PathSeparator) + plainName
	f, err := os.OpenFile(completeFileName, flags, perm)
	if err != nil {
		v.commands.SetError(ERR_FILENOTFOUND)
	} else {
		v.file[channel] = f
		data := make([]byte, 1)
		_, _ = f.Read(data)
		v.readChar[channel] = data[0]
	}
	return StOk
}

func (v *FSDrive) openDirectory(channel uint8, pattern string, dirName string) []byte {
	// Special treatment for "$0"
	if len(pattern) > 0 {
		if pattern[0] == '0' && len(pattern) == 1 {
			pattern = ""
		}
	}

	if p := strings.Index(pattern, ":"); p >= 0 {
		p++
		if len(pattern) < p {
			pattern = pattern[p:]
		}
	}

	const titleStart = "\001\004\001\001\000\000\022\042"
	const titleEnd = "\042 00 2A"
	const blocksFreeStart = "\001\001\000\000"
	const blockFreeEnd = "\000\000"

	title := make([]byte, 16)
	for idx := range title {
		title[idx] = ' '
		if idx < len(dirName) {
			title[idx] = dirName[idx]
		}
	}

	var buf []byte

	buf = append(buf, titleStart...)
	buf = append(buf, title...)
	buf = append(buf, titleEnd...)

	/*

	   // Create and write one line for every directory entry
	   std::vector<C64DirEntry>::const_iterator i, end = _file_info.end();
	   for (i = _file_info.begin(); i != end; i++) {
	           // Include only files matching the pattern
	           if (pattern_len == 0 || match(pattern, pattern_len, (uint8*)i->name)) {
	                   // Clear line with spaces and terminate with null byte
	                   memset(buf, ' ', 31);
	                   buf[31] = 0;

	                   uint8* p = (uint8*)buf;
	                   *p++ = 0x01;    // Dummy line link
	                   *p++ = 0x01;

	                   // Calculate size in blocks (254 bytes each)
	                   int n = (i->size + 254) / 254;
	                   *p++ = n & 0xff;
	                   *p++ = (n >> 8) & 0xff;

	                   p++;
	                   if (n < 10) p++;        // Less than 10: add one space
	                   if (n < 100) p++;       // Less than 100: add another space

	                   // Convert and insert file name
	                   *p++ = '\"';
	                   uint8* q = p;
	                   for (int j = 0; j < 16 && i->name[j]; j++) {
	                           *q++ = i->name[j];
	                   }
	                   *q++ = '\"';
	                   p += 18;

	                   // File type
	                   switch (i->type) {
	                           case FTYPE_DEL:
	                           *p++ = 'D';
	                           *p++ = 'E';
	                           *p++ = 'L';
	                           break;
	                           case FTYPE_SEQ:
	                           *p++ = 'S';
	                           *p++ = 'E';
	                           *p++ = 'Q';
	                           break;
	                           case FTYPE_PRG:
	                           *p++ = 'P';
	                           *p++ = 'R';
	                           *p++ = 'G';
	                           break;
	                           case FTYPE_USR:
	                           *p++ = 'U';
	                           *p++ = 'S';
	                           *p++ = 'R';
	                           break;
	                           case FTYPE_REL:
	                           *p++ = 'R';
	                           *p++ = 'E';
	                           *p++ = 'L';
	                           break;
	                           default:
	                           *p++ = '?';
	                           *p++ = '?';
	                           *p++ = '?';
	                           break;
	                   }

	                   // Write line
	                   fwrite(buf, 1, 32, _fileChannel[channel]);
	           }
	   }

	*/

	buf = append(buf, blocksFreeStart...)
	buf = append(buf, "BLOCKS FREE.             "...)
	buf = append(buf, blockFreeEnd...)

	// Rewind file for reading and read first byte
	//rewind(_fileChannel[channel]);
	//_read_char[channel] = (uint8)getc(_fileChannel[channel]);

	return buf
}

func (v *FSDrive) configChanged() {
}
