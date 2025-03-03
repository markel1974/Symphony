package fsdrive

import (
	"bytes"
	"fmt"
	"github.com/markel1974/c64emu/src/components/iec/virtualdrive"
)

type Commands struct {
	cmdBuf    []uint8 // Buffer for incoming command strings
	cmdLen    int     // Length of received command
	errorData []byte
	errorIdx  int
}

func NewCommands() *Commands {
	return &Commands{
		cmdBuf: make([]uint8, 64),
		cmdLen: 0,
	}
}

func (vd *Commands) SetError(e int) {
	vd.errorData = virtualdrive.Errors[e]
	vd.errorIdx = 0
}

func (vd *Commands) RetrieveError() uint8 {
	if vd.errorIdx < len(vd.errorData) {
		v := vd.errorData[vd.errorIdx]
		vd.errorIdx++
		return v
	}
	return '\r'
}

func (vd *Commands) CommandSet(data uint8) bool {
	if vd.cmdLen >= 58 {
		return false
	}
	vd.cmdBuf[vd.cmdLen] = data
	vd.cmdLen++
	return true
}

func (vd *Commands) CommandClear() {
	vd.cmdLen = 0
}

func (vd *Commands) CommandExecBuf() {
	//TODO IMPLEMENT
	vd.CommandExec(vd.cmdBuf)
	vd.cmdLen = 0
}

func (vd *Commands) CommandExec(cmd []uint8) (int, bool) {
	//TODO IMPLEMENT
	action := 0
	//vd.executeCmd(vd.cmdBuf, vd.cmdLen)
	if len(cmd) == 0 {
		vd.SetError(virtualdrive.ERR_SYNTAX31)
		return 0, false
	}
	// Strip trailing CRs
	cmdLen := len(cmd)
	for cmdLen > 0 && cmd[cmdLen-1] == 0x0d {
		cmdLen--
	}
	var colon []byte
	var equal []byte
	var comma []byte
	var minus []byte
	if pos := bytes.Index(cmd, []byte(":")); pos >= 0 {
		p0 := cmd[:pos]
		colon = cmd[pos:]
		if pos1 := bytes.Index(p0, []byte("=")); pos1 >= 0 {
			equal = p0[pos1:]
			//equal := colon ? (const uint8*)memchr(colon, '=', cmdLen - (colon - cmd)) : NULL
		}
	}
	if pos := bytes.Index(cmd, []byte(",")); pos >= 0 {
		comma = cmd[pos:]
	}
	if pos := bytes.Index(cmd, []byte("-")); pos >= 0 {
		minus = cmd[pos:]
	}

	// Find token delimiters
	//colon := (const uint8*)memchr(cmd, ':', cmdLen)
	//equal := colon ? (const uint8*)memchr(colon, '=', cmdLen - (colon - cmd)) : NULL
	//comma := (const uint8*)memchr(cmd, ',', cmdLen)
	//minus := (const uint8*)memchr(cmd, '-', cmdLen)

	// Parse command name
	//vd.SetError(ERR_OK)
	switch cmd[0] {
	case 'B': // Block/buffer
		if len(minus) == 0 {
			vd.SetError(virtualdrive.ERR_SYNTAX31)
			return 0, false
		}
		// Parse arguments (up to 4 decimal numbers separated by
		// space, cursor right or comma)
		//TODO IMPLEMENT
		fmt.Println("CommandExec: TODO IMPLEMENT")
		l := 0
		//l := cmd + 3
		//if len(colon) > 0 {
		//	l = colon + 1
		//}
		p := colon[l:]
		//const uint8* p = colon ? colon + 1 : cmd + 3
		arg1, arg2, arg3, arg4 := vd.parseBlockCmdArgs(p)
		// Switch on command
		switch minus[1] {
		case 'R':
			vd.blockReadCmd(arg1, arg3, arg4, false)
		case 'W':
			vd.blockWriteCmd(arg1, arg3, arg4, false)
		case 'E':
			vd.blockExecuteCmd(arg1, arg3, arg4)
		case 'A':
			vd.blockAllocateCmd(arg2, arg3)
		case 'F':
			vd.blockFreeCmd(arg2, arg3)
		case 'P':
			vd.bufferPointerCmd(arg1, arg2)
		default:
			vd.SetError(virtualdrive.ERR_SYNTAX31)
			return 0, false
		}

	case 'M':
		// Memory
		if cmd[1] != '-' {
			vd.SetError(virtualdrive.ERR_SYNTAX31)
			return 0, false
		}
		// Read parameters
		cLen := cmd[5]
		addr := uint16(cmd[3]) | (uint16(cmd[4]) << 8)
		// Switch on command
		switch cmd[2] {
		case 'R':
			l := cLen
			if cmdLen < 6 {
				l = 1
			}
			vd.memReadCmd(addr, l)
		case 'W':
			vd.memWriteCmd(addr, cLen, cmd[6:])
		case 'E':
			vd.memExecuteCmd(addr)
		default:
			vd.SetError(virtualdrive.ERR_SYNTAX31)
			return 0, false
		}

	case 'C':
		// Copy
		//TODO VERIFY
		fmt.Println("CommandExec: TODO VERIFY")
		test := len(comma) > 0 && len(comma) < len(equal)
		if len(colon) == 0 {
			vd.SetError(virtualdrive.ERR_SYNTAX31)
			return 0, false
		}
		if len(equal) == 0 || bytes.Index(cmd, []byte("*")) > 0 || bytes.Index(cmd, []byte("?")) > 0 || test {
			vd.SetError(virtualdrive.ERR_SYNTAX30)
			return 0, false
		}
		//TODO IMPLEMENT
		fmt.Println("CommandExec: TODO IMPLEMENT")
		var newFile []byte
		var oldFile []byte
		//tn := equal - colon - 1
		//newFile := colon[1:tn]
		//to := cmdLen - (equal + 1 - cmd)
		//oldFile := equal[1:to]
		vd.copyCmd(newFile, oldFile)

	case 'R': // Rename
		if len(colon) == 0 {
			vd.SetError(virtualdrive.ERR_SYNTAX34)
			return 0, false
		}
		if len(equal) == 0 || len(comma) > 0 || bytes.Index(cmd, []byte("*")) > 0 || bytes.Index(cmd, []byte("?")) > 0 {
			vd.SetError(virtualdrive.ERR_SYNTAX30)
			return 0, false
		}
		//TODO IMPLEMENT
		fmt.Println("CommandExec: TODO IMPLEMENT")
		var newFile []byte
		var oldFile []byte
		//tn := equal - colon - 1
		//newFile := colon[1:tn]
		//to := cmdLen - (equal + 1 - cmd)
		//oldFile := equal[1:to]
		vd.renameCmd(newFile, oldFile)

	case 'S':
		// Scratch
		if len(colon) == 0 {
			vd.SetError(virtualdrive.ERR_SYNTAX34)
			return 0, false
		}
		//TODO IMPLEMENT
		fmt.Println("CommandExec: TODO IMPLEMENT")
		var t []byte
		//l := cmdLen - (colon + 1 - cmd)
		//t := colon[1:l]
		vd.scratchCmd(t)

	case 'P':
		// Position
		vd.positionCmd(cmd[1:])

	case 'I': // Initialize
		vd.initializeCmd()

	case 'N': // New (format)
		if len(colon) == 0 {
			vd.SetError(virtualdrive.ERR_SYNTAX34)
			return 0, false
		}
		//TODO IMPLEMENT
		fmt.Println("CommandExec: TODO IMPLEMENT")
		var t []byte
		//l := comma - colon - 1
		//if len(comma) > 0 {
		//	l = cmdLen - (colon + 1 - cmd)
		//}
		//t := colon[1:l]
		vd.newCmd(t, comma)

	case 'V': // Validate
		vd.validateCmd()

	case 'U': // User
		if cmd[1] == '0' {
			//Nothing to do...
		} else {
			switch cmd[1] & 0x0f {
			case 1:
				// U1/UA: Read block
				p := cmd[2:]
				if len(colon) > 0 {
					p = colon[1:]
				}
				arg1, _, arg3, arg4 := vd.parseBlockCmdArgs(p)
				vd.blockReadCmd(arg1, arg3, arg4, true)
			case 2:
				// U2/UB: Write block
				p := cmd[2:]
				if len(colon) > 0 {
					p = colon[1:]
				}
				arg1, _, arg3, arg4 := vd.parseBlockCmdArgs(p)
				vd.blockWriteCmd(arg1, arg3, arg4, true)
			case 9:
				// U9/UI: C64/VC20 mode switch
				if cmd[2] != '+' && cmd[2] != '-' {
					action = 1 //RESET
				}
			case 10: // U:/UJ: Reset
				action = 1 //RESET
			default:
				vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
				return 0, false
			}
		}
	default:
		vd.SetError(virtualdrive.ERR_SYNTAX31)
		return 0, false
	}
	vd.SetError(virtualdrive.ERR_OK)
	return action, true
}

func (vd *Commands) parseBlockCmdArgs(p []uint8) (int, int, int, int) {
	arg1 := 0
	arg2 := 0
	arg3 := 0
	arg4 := 0
	pIdx := 0
	if len(p) == 0 {
		return 0, 0, 0, 0
	}
	for p[pIdx] == ' ' || p[pIdx] == 0x1d || p[pIdx] == ',' {
		pIdx++
	}
	for p[pIdx] >= '0' && p[pIdx] < '@' {
		arg1 = arg1*10 + (int(p[pIdx]) & 0x0f)
		pIdx++
	}
	for p[pIdx] == ' ' || p[pIdx] == 0x1d || p[pIdx] == ',' {
		pIdx++
	}
	for p[pIdx] >= '0' && p[pIdx] < '@' {
		arg2 = arg2*10 + (int(p[pIdx]) & 0x0f)
		pIdx++
	}
	for p[pIdx] == ' ' || p[pIdx] == 0x1d || p[pIdx] == ',' {
		pIdx++
	}
	for p[pIdx] >= '0' && p[pIdx] < '@' {
		arg3 = arg3*10 + (int(p[pIdx]) & 0x0f)
		pIdx++
	}
	for p[pIdx] == ' ' || p[pIdx] == 0x1d || p[pIdx] == ',' {
		pIdx++
	}
	for p[pIdx] >= '0' && p[pIdx] < '@' {
		arg4 = arg4*10 + (int(p[pIdx]) & 0x0f)
		pIdx++
	}
	return arg1, arg2, arg3, arg4
}

func (vd *Commands) blockReadCmd(channel int, track int, sector int, userCmd bool) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) blockWriteCmd(channel int, track int, sector int, userCmd bool) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) blockExecuteCmd(channel int, track int, sector int) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

// BLOCK-ALLOCATE:0,track,sector
func (vd *Commands) blockAllocateCmd(track int, sector int) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

// BLOCK-FREE:0,track,sector
func (vd *Commands) blockFreeCmd(track int, sector int) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) bufferPointerCmd(channel int, pos int) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

// M-R<adr low><adr high>[<number>]
func (vd *Commands) memReadCmd(addr uint16, length uint8) {
	vd.unsupportedCmd()
	//TODO
	//error_ptr = error_buf
	//error_buf[0] = 0
	//error_len = 0
	vd.SetError(virtualdrive.ERR_OK)
}

// M-W<adr low><adr high><number><data...>
func (vd *Commands) memWriteCmd(addr uint16, length uint8, p []uint8) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

// M-E<adr low><adr high>
func (vd *Commands) memExecuteCmd(addr uint16) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) copyCmd(newFile []uint8, oldFiles []uint8) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) renameCmd(newFile []uint8, oldFile []uint8) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) scratchCmd(file []uint8) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

// P<channel><record low><record high><byte>
func (vd *Commands) positionCmd(cmd []uint8) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

// INITIALIZE
func (vd *Commands) initializeCmd() {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) newCmd(name []uint8, comma []uint8) {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) validateCmd() {
	vd.SetError(virtualdrive.ERR_UNIMPLEMENTED)
}

func (vd *Commands) unsupportedCmd() {
}
