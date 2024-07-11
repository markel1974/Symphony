package virtualdrive

const (
	StOk          = 0    // No error
	StReadTimeout = 0x02 // Timeout on reading
	StTimeout     = 0x03 // Timeout
	StEof         = 0x40 // End of file
	StNotPresent  = 0x80 // Device not present
)

const (
	DRVLED_OFF   = iota // Inactive LED off
	DRVLED_ON           // Active LED on
	DRVLED_ERROR        // Error blink LED
)

const (
	FILE_IMAGE_TYPE = iota
	FILE_ARCH_TYPE
	FILE_UNKNOWN_TYPE
)

const (
	ERR_OK            = iota // 00 OK
	ERR_SCRATCHED            // 01 FILES SCRATCHED
	ERR_UNIMPLEMENTED        // 03 UNIMPLEMENTED
	ERR_READ20               // 20 READ ERROR (block header not found)
	ERR_READ21               // 21 READ ERROR (no sync character)
	ERR_READ22               // 22 READ ERROR (data block not present)
	ERR_READ23               // 23 READ ERROR (checksum error in data block)
	ERR_READ24               // 24 READ ERROR (byte decoding error)
	ERR_WRITE25              // 25 WRITE ERROR (write-verify error)
	ERR_WRITEPROTECT         // 26 WRITE PROTECT ON
	ERR_READ27               // 27 READ ERROR (checksum error in header)
	ERR_WRITE28              // 28 WRITE ERROR (long data block)
	ERR_DISKID               // 29 DISK ID MISMATCH
	ERR_SYNTAX30             // 30 SYNTAX ERROR (general syntax)
	ERR_SYNTAX31             // 31 SYNTAX ERROR (invalid command)
	ERR_SYNTAX32             // 32 SYNTAX ERROR (command too long)
	ERR_SYNTAX33             // 33 SYNTAX ERROR (wildcards on writing)
	ERR_SYNTAX34             // 34 SYNTAX ERROR (missing file name)
	ERR_WRITEFILEOPEN        // 60 WRITE FILE OPEN
	ERR_FILENOTOPEN          // 61 FILE NOT OPEN
	ERR_FILENOTFOUND         // 62 FILE NOT FOUND
	ERR_FILEEXISTS           // 63 FILE EXISTS
	ERR_FILETYPE             // 64 FILE TYPE MISMATCH
	ERR_NOBLOCK              // 65 NO BLOCK
	ERR_ILLEGALTS            // 66 ILLEGAL TRACK OR SECTOR
	ERR_NOCHANNEL            // 70 NO CHANNEL
	ERR_DIRERROR             // 71 DIR ERROR
	ERR_DISKFULL             // 72 DISK FULL
	ERR_STARTUP              // 73 Power-up message
	ERR_NOTREADY             // 74 DRIVE NOT READY
)

// 1541 file types
const (
	FTYPE_DEL = iota // Deleted
	FTYPE_SEQ        // Sequential
	FTYPE_PRG        // Program
	FTYPE_USR        // User
	FTYPE_REL        // Relative
	FTYPE_UNKNOWN
)

const (
	FMODE_READ   = iota // Read
	FMODE_WRITE         // Write
	FMODE_APPEND        // Append
	FMODE_M             // Read open file
)

var Errors = []string{
	"00, OK,%02d,%02d\x0d",
	"01, FILES SCRATCHED,%02d,%02d\x0d",
	"03, UNIMPLEMENTED,%02d,%02d\x0d",
	"20, READ ERROR,%02d,%02d\x0d",
	"21, READ ERROR,%02d,%02d\x0d",
	"22, READ ERROR,%02d,%02d\x0d",
	"23, READ ERROR,%02d,%02d\x0d",
	"24, READ ERROR,%02d,%02d\x0d",
	"25, WRITE ERROR,%02d,%02d\x0d",
	"26, WRITE PROTECT ON,%02d,%02d\x0d",
	"27, READ ERROR,%02d,%02d\x0d",
	"28, WRITE ERROR,%02d,%02d\x0d",
	"29, DISK ID MISMATCH,%02d,%02d\x0d",
	"30, SYNTAX ERROR,%02d,%02d\x0d",
	"31, SYNTAX ERROR,%02d,%02d\x0d",
	"32, SYNTAX ERROR,%02d,%02d\x0d",
	"33, SYNTAX ERROR,%02d,%02d\x0d",
	"34, SYNTAX ERROR,%02d,%02d\x0d",
	"60, WRITE FILE OPEN,%02d,%02d\x0d",
	"61, FILE NOT OPEN,%02d,%02d\x0d",
	"62, FILE NOT FOUND,%02d,%02d\x0d",
	"63, FILE EXISTS,%02d,%02d\x0d",
	"64, FILE TYPE MISMATCH,%02d,%02d\x0d",
	"65, NO BLOCK,%02d,%02d\x0d",
	"66, ILLEGAL TRACK OR SECTOR,%02d,%02d\x0d",
	"70, NO CHANNEL,%02d,%02d\x0d",
	"71, DIR ERROR,%02d,%02d\x0d",
	"72, DISK FULL,%02d,%02d\x0d",
	"73, CBM DOS V2.6 1541,%02d,%02d\x0d",
	"74, DRIVE NOT READY,%02d,%02d\x0d",
}
