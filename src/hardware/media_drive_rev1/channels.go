package media_drive_rev1

import (
	"github.com/markel1974/c64emu/src/hardware/media_drive_rev1/adapters"
	"strconv"
	"strings"
)

type Channels struct {
	channels [16]*Channel
	adapter  adapters.IAdapter
	cmdLen   int
	cmdBuf   []uint8
}

func NewChannels(adapter adapters.IAdapter) *Channels {
	c := &Channels{
		adapter: adapter,
		cmdBuf:  make([]uint8, 64),
		cmdLen:  0,
	}
	for idx := range c.channels {
		c.channels[idx] = NewChannel(idx, adapter)
		c.channels[idx].Reset()
	}
	return c
}

func (c *Channels) Reset() {
	for _, z := range c.channels {
		z.Reset()
	}
	c.cmdLen = 0
}

func (c *Channels) Close() {
	for _, z := range c.channels {
		_ = z.Close()
	}
	c.cmdLen = 0
}

func (c *Channels) SetError(idx uint8, err error) {
	idx &= 0xf
	c.channels[idx].SetError(err)
}

func (c *Channels) Get(idx uint8) *Channel {
	idx &= 0xf
	channel := c.channels[idx]
	return channel
}

// CommandSet adds a byte of data to the command buffer if the buffer has not reached its maximum capacity.
// Returns false if the buffer is full, otherwise increments the buffer length and returns true.
func (c *Channels) CommandSet(data uint8) bool {
	if c.cmdLen >= 58 {
		return false
	}
	c.cmdBuf[c.cmdLen] = data
	c.cmdLen++
	return true
}

// CommandExecBuf executes the command stored in the internal buffer and resets the buffer length to zero, returning result or error.
func (c *Channels) CommandExecBuf() (int, error) {
	// TODO IMPLEMENT: This method should be updated to accept channels and adapter
	// or be refactored to use CommandExec directly.
	v, err := c.CommandExec(c.cmdBuf[:c.cmdLen])
	c.cmdLen = 0
	return v, err
}

// CommandExec parses and executes a given command string, returning an action code or an error on failure.
func (c *Channels) CommandExec(cmd []uint8) (int, error) {
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
		if commandData == "" {
			return 0, adapters.Error(adapters.ErrSyntax34)
		}
		if err := c.adapter.ScratchFile(commandData); err != nil {
			return 0, err
		}
		return 0, nil

	case 'R': // RENAME
		renameParts := strings.SplitN(commandData, "=", 2)
		if len(renameParts) != 2 {
			return 0, adapters.Error(adapters.ErrSyntax30)
		}
		newName := renameParts[0]
		oldName := renameParts[1]
		if err := c.adapter.RenameFile(oldName, newName); err != nil {
			return 0, err
		}
		return 0, nil
	case 'C': // COPY
		copyParts := strings.SplitN(commandData, "=", 2)
		if len(copyParts) != 2 {
			return 0, adapters.Error(adapters.ErrSyntax30)
		}
		newName := copyParts[0]
		oldFilesStr := copyParts[1]
		oldFiles := strings.Split(oldFilesStr, ",")
		var combinedData []byte
		for _, oldFile := range oldFiles {
			data, err := c.adapter.ReadFile(oldFile)
			if err != nil {
				return 0, err
				//channels[errChannel].Shutdown(adapters.Error(adapters.ErrFileNotFound))
			}
			combinedData = append(combinedData, data...)
		}
		if err := c.adapter.WriteFile(newName, combinedData); err != nil {
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
		if err := c.adapter.Format(diskName, diskId); err != nil {
			return 0, err
		}
		return 0, nil

	case 'I': // INITIALIZE
		// This command resets the drive state, rereads the BAM (Block Availability Map).
		// For the MediaDrive, this could mean clearing caches or re-checking the backend.
		// For now, we can consider it a soft reset.
		action = 1 // Signal a RESET action
		if err := c.adapter.Reset(); err != nil {
			return 0, err
		}
		return 0, nil

	case 'V': // VALIDATE
		// Similar to INITIALIZE, but more thorough. In a real 1541, it would rebuild the BAM.
		// For a virtual drive, this might not have a direct equivalent.
		// We can treat it as an OK command.
		if err := c.adapter.Validate(); err != nil {
			return 0, err
		}
		return 0, nil
	// Note: Block and Memory commands (B-* and M-*) are for low-level access
	// and are harder to map to a high-level adapter. They would typically
	// work with a block-based adapter (like a D64 image).
	// The current adapter is file-based.
	case 'B': // Block/buffer
		subParts := strings.Split(commandData, ",")
		if len(subParts) < 3 {
			return 0, adapters.Error(adapters.ErrSyntax31)
		}
		channelIdx, _ := strconv.Atoi(subParts[0])
		track, _ := strconv.Atoi(subParts[1])
		sector, _ := strconv.Atoi(subParts[2])
		channel := c.Get(uint8(channelIdx))
		subCmd := commandName[2:] // Assume formato B-R, B-W...
		switch subCmd {
		case "R":
			err := c.adapter.BlockRead(channel, track, sector)
			return 0, err
		case "W":
			err := c.adapter.BlockWrite(channel, track, sector)
			return 0, err
		default:
			return 0, adapters.Error(adapters.ErrSyntax31)
		}

	case 'M': // Memory
		// 1. Esegui il parsing (indirizzo, lunghezza)
		subCmd := commandName[2:3]
		addr, _ := strconv.ParseInt(commandData[0:4], 16, 16)
		length, _ := strconv.Atoi(commandData[4:])
		switch subCmd {
		case "R":
			// M-R mette i dati letti nel canale di errore per essere letti
			data, err := c.adapter.MemoryRead(uint16(addr), length)
			if err != nil {
				return 0, err
			}
			channel := c.Get(errChannel)
			channel.dataSet(data)
			return 0, nil
		case "W":
			if len(cmd) < 6+length {
				return 0, adapters.Error(adapters.ErrSyntax31) // Not enough data sent
			}
			dataFromChannel := cmd[6 : 6+length]
			err := c.adapter.MemoryWrite(uint16(addr), dataFromChannel)
			if err != nil {
				return 0, err
			}
			return 0, nil
		case "E": // Memory-Execute
			return 0, c.adapter.MemoryExec(uint16(addr))
		default:
			return 0, adapters.Error(adapters.ErrSyntax31)
		}

	case 'P': // Position
		subParts := strings.Split(commandData, ",")
		if len(subParts) < 2 {
			return 0, adapters.Error(adapters.ErrSyntax31)
		}
		channelIdx, _ := strconv.Atoi(subParts[0])
		position, _ := strconv.Atoi(subParts[1])
		channel := c.Get(uint8(channelIdx))
		err := c.adapter.Position(channel, position)
		return 0, err

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
