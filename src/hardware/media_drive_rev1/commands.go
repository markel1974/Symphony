package media_drive_rev1

import (
	"github.com/markel1974/c64emu/src/hardware/media_drive_rev1/adapters"
	"strconv"
	"strings"
)

// Commands provides a structure for handling command execution and communication with an IAdapter interface.
// It maintains a buffer for command input, tracks the current command length, and interfaces with file operations.
type Commands struct {
	cmdBuf  []uint8 // Buffer for incoming command strings
	cmdLen  int     // Length of received command
	adapter adapters.IAdapter
}

// NewCommands initializes and returns a new Commands instance with the provided IAdapter for file and archival operations.
func NewCommands(adapter adapters.IAdapter) *Commands {
	return &Commands{
		cmdBuf:  make([]uint8, 64),
		cmdLen:  0,
		adapter: adapter,
	}
}

// CommandSet adds a byte of data to the command buffer if the buffer has not reached its maximum capacity.
// Returns false if the buffer is full, otherwise increments the buffer length and returns true.
func (vd *Commands) CommandSet(data uint8) bool {
	if vd.cmdLen >= 58 {
		return false
	}
	vd.cmdBuf[vd.cmdLen] = data
	vd.cmdLen++
	return true
}

// CommandClear resets the command length to zero, effectively clearing the current command buffer.
func (vd *Commands) CommandClear() {
	vd.cmdLen = 0
}

// CommandExecBuf executes the command stored in the internal buffer and resets the buffer length to zero, returning result or error.
func (vd *Commands) CommandExecBuf() (int, error) {
	// TODO IMPLEMENT: This method should be updated to accept channels and adapter
	// or be refactored to use CommandExec directly.
	v, err := vd.CommandExec(vd.cmdBuf[:vd.cmdLen])
	vd.cmdLen = 0
	return v, err
}

// CommandExec parses and executes a given command string, returning an action code or an error on failure.
func (vd *Commands) CommandExec(cmd []uint8) (int, error) {
	action := 0
	if len(cmd) == 0 {
		return 0, adapters.Error(adapters.ErrSyntax31)
	}
	// Strip trailing Carriage Returns (CR)
	cmdLen := len(cmd)
	for cmdLen > 0 && cmd[cmdLen-1] == 0x0d {
		cmdLen--
	}
	cmd = cmd[:cmdLen]
	cmdStr := strings.ToUpper(string(cmd))
	parts := strings.SplitN(cmdStr, ":", 2)
	commandName := parts[0]
	var commandData string
	if len(parts) > 1 {
		commandData = parts[1]
	}
	switch commandName[0] {
	case 'S': // SCRATCH
		// 1. Validate syntax: must have a filename after the colon.
		if commandData == "" {
			return 0, adapters.Error(adapters.ErrSyntax34)
		}
		if err := vd.adapter.ScratchFile(commandData); err != nil {
			return 0, err
		}
		return 0, nil

	case 'R': // RENAME
		// 1. Validate syntax: must contain a colon and an equals sign.
		renameParts := strings.SplitN(commandData, "=", 2)
		if len(renameParts) != 2 {
			return 0, adapters.Error(adapters.ErrSyntax30)
		}
		newName := renameParts[0]
		oldName := renameParts[1]
		if err := vd.adapter.RenameFile(oldName, newName); err != nil {
			return 0, err
		}
		return 0, nil
	case 'C': // COPY
		// 1. Validate syntax: must contain a colon and an equals sign.
		copyParts := strings.SplitN(commandData, "=", 2)
		if len(copyParts) != 2 {
			return 0, adapters.Error(adapters.ErrSyntax30)
		}
		newName := copyParts[0]
		oldFilesStr := copyParts[1]
		oldFiles := strings.Split(oldFilesStr, ",")
		var combinedData []byte
		for _, oldFile := range oldFiles {
			data, err := vd.adapter.ReadFile(oldFile)
			if err != nil {
				return 0, err
				//channels[errChannel].SetError(adapters.Error(adapters.ErrFileNotFound))
			}
			combinedData = append(combinedData, data...)
		}
		if err := vd.adapter.WriteFile(newName, combinedData); err != nil {
			return 0, err
		}
		return 0, nil

	case 'N': // NEW (Format)
		formatParts := strings.SplitN(commandData, ",", 2)
		diskName := formatParts[0]
		var diskId string
		if len(formatParts) > 1 {
			diskId = formatParts[1]
		}
		if err := vd.adapter.Format(diskName, diskId); err != nil {
			return 0, err
		}
		return 0, nil

	case 'I': // INITIALIZE
		// This command resets the drive state, rereads the BAM (Block Availability Map).
		// For the MediaDrive, this could mean clearing caches or re-checking the backend.
		// For now, we can consider it a soft reset.
		action = 1 // Signal a RESET action
		if err := vd.adapter.Reset(); err != nil {
			return 0, err
		}
		return 0, nil

	case 'V': // VALIDATE
		// Similar to INITIALIZE, but more thorough. In a real 1541, it would rebuild the BAM.
		// For a virtual drive, this might not have a direct equivalent.
		// We can treat it as an OK command.
		if err := vd.adapter.Validate(); err != nil {
			return 0, err
		}
		return 0, nil
	// Note: Block and Memory commands (B-* and M-*) are for low-level access
	// and are harder to map to a high-level adapter. They would typically
	// work with a block-based adapter (like a D64 image).
	// The current adapter is file-based.
	case 'B': // Block/buffer
		// 1. Esegui il parsing degli argomenti (canale, traccia, settore)
		parts1 := strings.Split(commandData, ",")
		if len(parts1) < 3 {
			return 0, adapters.Error(adapters.ErrSyntax31)
		}
		channel, _ := strconv.Atoi(parts1[0])
		track, _ := strconv.Atoi(parts1[1])
		sector, _ := strconv.Atoi(parts1[2])
		_ = channel
		_ = track
		_ = sector
		subCmd := commandName[2:] // Assume formato B-R, B-W...
		switch subCmd {
		case "R":
			//err := vd.adapter.BlockRead(channels[channel], track, sector)
			//return 0, err
			return 0, adapters.Error(adapters.ErrUnimplemented)
		case "W":
			//err := vd.adapter.BlockWrite(channels[channel], track, sector)
			//return 0, err
			return 0, adapters.Error(adapters.ErrUnimplemented)
		default:
			return 0, adapters.Error(adapters.ErrSyntax31)
		}

	case 'M': // Memory
		// 1. Esegui il parsing (indirizzo, lunghezza)
		subCmd := commandName[2:3]
		addr, _ := strconv.ParseInt(commandData[0:4], 16, 16)
		length, _ := strconv.Atoi(commandData[4:])
		_ = addr
		_ = length
		switch subCmd {
		case "R":
			//data, err := vd.adapter.MemoryRead(uint16(addr), length)
			//if err != nil {
			//	return 0, err
			//}
			// M-R mette i dati letti nel canale di errore per essere letti dal C64
			//channels[errChannel].dataSet(data)
			//return 0, nil
			return 0, adapters.Error(adapters.ErrUnimplemented)
		case "W":
			// Per M-W, i dati dovrebbero essere già nel buffer di un canale.
			// La logica qui sarebbe più complessa, ma la delega è la stessa.
			// return 0, vd.adapter.MemoryWrite(uint16(addr), dataFromChannel)
			return 0, adapters.Error(adapters.ErrUnimplemented) // Lasciamo non implementato per ora
		default:
			return 0, adapters.Error(adapters.ErrSyntax31)
		}

	case 'P': // Position
		// 1. Esegui il parsing (canale, posizione)
		parts1 := strings.Split(commandData, ",")
		if len(parts1) < 2 {
			return 0, adapters.Error(adapters.ErrSyntax31)
		}
		channel, _ := strconv.Atoi(parts1[0])
		position, _ := strconv.Atoi(parts1[1])
		_ = channel
		_ = position
		//err := vd.adapter.Position(channels[channel], position)
		//return 0, err
		return 0, adapters.Error(adapters.ErrUnimplemented)

	case 'U': // User Command
		if len(commandName) > 1 {
			switch commandName[1] {
			case 'J', ':': // UJ or U: is Reset
				action = 1 //RESET
				return action, nil
			}
		}
		return 0, adapters.Error(adapters.ErrUnimplemented)

	default:
		return 0, adapters.Error(adapters.ErrSyntax31)
	}
}
