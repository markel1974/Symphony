package session

import (
	"strconv"
)

// IOCode represents an integer code used for defining input/output operations or control commands.
type IOCode int

// byteToCode maps a byte value to its corresponding IOCode representation.
var byteToCode map[byte]IOCode

// codeToByte maps IOCode values to their corresponding byte representations for protocol command handling.
var codeToByte map[IOCode]byte

// NUL represents NULL or no operation.
// ECHO represents the Echo operation.
// SGA represents Suppress-go-ahead.
// ST represents Status.
// TM represents Timing mark.
// BEL represents Bell.
// BS represents Backspace.
// HT represents Horizontal tab.
// LF represents Line feed.
// FF represents Form feed.
// CR represents Carriage return.
// TT represents Terminal type.
// WS represents Window size.
// TS represents Terminal speed.
// RFC represents Remote flow control.
// LM represents Line mode.
// EV represents Environment variables.
// SE represents End of sub negotiation parameters.
// NOP represents No operation.
// DM represents Data Mark, accompanied by a TCP Urgent notification.
// BRK represents Break, or NVT character BRK.
// IP represents Interrupt Process.
// AO represents Abort Output.
// AYT represents Are You There.
// EC represents Erase Character.
// EL represents Erase Line.
// GA represents Go Ahead signal.
// SB indicates sub negotiation of the indicated option.
// WILL indicates a desire to begin or confirmation of performing the indicated option.
// WONT indicates refusal or cessation of performing the indicated option.
// DO indicates a request or confirmation expecting the other party to perform the indicated option.
// DONT indicates a demand or confirmation the other party stops performing the indicated option.
// IAC represents Interpret as Command.
// CMP1 represents MCCP Compress.
// CMP2 represents MCCP Compress2.
// AARD represents Aardwolf MUD out of band communication.
// ATCP represents Achaea Telnet Client Protocol.
// GMCP represents Generic Mud Communication Protocol.
const (
	NUL  IOCode = iota // NULL, no operation
	ECHO IOCode = iota // Echo
	SGA  IOCode = iota // Suppress go ahead
	ST   IOCode = iota // Status
	TM   IOCode = iota // Timing mark
	BEL  IOCode = iota // Bell
	BS   IOCode = iota // Backspace
	HT   IOCode = iota // Horizontal tab
	LF   IOCode = iota // Line feed
	FF   IOCode = iota // Form feed
	CR   IOCode = iota // Carriage return
	TT   IOCode = iota // Terminal type
	WS   IOCode = iota // Window size
	TS   IOCode = iota // Terminal speed
	RFC  IOCode = iota // Remote flow control
	LM   IOCode = iota // Line mode
	EV   IOCode = iota // Environment variables
	SE   IOCode = iota // End of sub negotiation parameters.
	NOP  IOCode = iota // No operation.
	DM   IOCode = iota // Data Mark. The data stream portion of a Sync. This should always be accompanied by a TCP Urgent notification.
	BRK  IOCode = iota // Break. NVT character BRK.
	IP   IOCode = iota // Interrupt Process
	AO   IOCode = iota // Abort output
	AYT  IOCode = iota // Are you there
	EC   IOCode = iota // Erase character
	EL   IOCode = iota // Erase line
	GA   IOCode = iota // Go ahead signal
	SB   IOCode = iota // Indicates that what follows is sub negotiation of the indicated option.
	WILL IOCode = iota // Indicates the desire to begin performing, or confirmation that you are now performing, the indicated option.
	WONT IOCode = iota // Indicates the refusal to perform, or continue performing, the indicated option.
	DO   IOCode = iota // Indicates the request that the other party perform, or confirmation that you are expecting the other party to perform, the indicated option.
	DONT IOCode = iota // Indicates the demand that the other party stop performing, or confirmation that you are no longer expecting the other party to perform, the indicated option.
	IAC  IOCode = iota // Interpret as command

	// Non-standard codes:
	CMP1 IOCode = iota // MCCP Compress
	CMP2 IOCode = iota // MCCP Compress2
	AARD IOCode = iota // Aardwolf MUD out of band communication, http://www.aardwolf.com/blog/2008/07/10/telnet-negotiation-control-mud-client-interaction/
	ATCP IOCode = iota // Achaea Telnet Client Protocol, http://www.ironrealms.com/rapture/manual/files/FeatATCP-txt.html
	GMCP IOCode = iota // Generic Mud Communication Protocol
)

// init initializes the mappings between byte values and IOCode enums for bidirectional lookup.
func init() {
	byteToCode = map[byte]IOCode{}
	codeToByte = map[IOCode]byte{}

	codeToByte[NUL] = '\x00'
	codeToByte[ECHO] = '\x01'
	codeToByte[SGA] = '\x03'
	codeToByte[ST] = '\x05'
	codeToByte[TM] = '\x06'
	codeToByte[BEL] = '\x07'
	codeToByte[BS] = '\x08'
	codeToByte[HT] = '\x09'
	codeToByte[LF] = '\x0a'
	codeToByte[FF] = '\x0c'
	codeToByte[CR] = '\x0d'
	codeToByte[TT] = '\x18'
	codeToByte[WS] = '\x1F'
	codeToByte[TS] = '\x20'
	codeToByte[RFC] = '\x21'
	codeToByte[LM] = '\x22'
	codeToByte[EV] = '\x24'
	codeToByte[SE] = '\xf0'
	codeToByte[NOP] = '\xf1'
	codeToByte[DM] = '\xf2'
	codeToByte[BRK] = '\xf3'
	codeToByte[IP] = '\xf4'
	codeToByte[AO] = '\xf5'
	codeToByte[AYT] = '\xf6'
	codeToByte[EC] = '\xf7'
	codeToByte[EL] = '\xf8'
	codeToByte[GA] = '\xf9'
	codeToByte[SB] = '\xfa'
	codeToByte[WILL] = '\xfb'
	codeToByte[WONT] = '\xfc'
	codeToByte[DO] = '\xfd'
	codeToByte[DONT] = '\xfe'
	codeToByte[IAC] = '\xff'

	codeToByte[CMP1] = '\x55'
	codeToByte[CMP2] = '\x56'
	codeToByte[AARD] = '\x66'
	codeToByte[ATCP] = '\xc8'
	codeToByte[GMCP] = '\xc9'

	for enum, code := range codeToByte {
		byteToCode[code] = enum
	}
}

// _ converts a slice of bytes into a string representation, converting each byte using ByteToCodeString and separating with spaces.
func _(bytes []byte) string {
	str := ""
	for _, b := range bytes {
		if str != "" {
			str = str + " "
		}
		str = str + ByteToCodeString(b)
	}
	return str
}

// ByteToCodeString converts a given byte to its corresponding code string representation or a fallback formatted string if not found.
func ByteToCodeString(b byte) string {
	code, found := byteToCode[b]
	if !found {
		return "??(" + strconv.Itoa(int(b)) + ")"
	}
	return CodeToString(code)
}

// CodeToString converts an IOCode value to its corresponding string representation. If not matched, returns an empty string.
func CodeToString(code IOCode) string {
	switch code {
	case NUL:
		return "NUL"
	case ECHO:
		return "ECHO"
	case SGA:
		return "SGA"
	case ST:
		return "ST"
	case TM:
		return "TM"
	case BEL:
		return "BEL"
	case BS:
		return "BS"
	case HT:
		return "HT"
	case LF:
		return "LF"
	case FF:
		return "FF"
	case CR:
		return "CR"
	case TT:
		return "TT"
	case WS:
		return "WS"
	case TS:
		return "TS"
	case RFC:
		return "RFC"
	case LM:
		return "LM"
	case EV:
		return "EV"
	case SE:
		return "SE"
	case NOP:
		return "NOP"
	case DM:
		return "DM"
	case BRK:
		return "BRK"
	case IP:
		return "IP"
	case AO:
		return "AO"
	case AYT:
		return "AYT"
	case EC:
		return "EC"
	case EL:
		return "EL"
	case GA:
		return "GA"
	case SB:
		return "SB"
	case WILL:
		return "WILL"
	case WONT:
		return "WONT"
	case DO:
		return "DO"
	case DONT:
		return "DONT"
	case IAC:
		return "IAC"
	case CMP1:
		return "CMP1"
	case CMP2:
		return "CMP2"
	case AARD:
		return "AARD"
	case ATCP:
		return "ATCP"
	case GMCP:
		return "GMCP"
	}

	return ""
}
