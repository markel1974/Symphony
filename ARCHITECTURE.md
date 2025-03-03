# g64 - Commodore 64 Emulator in Go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# g64 - Commodore 64 Emulator - Architectural Overview

## 1. Introduction

g64 is a Commodore 64 emulator written in Go. It aims for high accuracy, including cycle-accurate emulation of the VIC-II video chip and the 6510 CPU. The emulator is designed to be modular, with distinct components responsible for different aspects of the C64 hardware. It does not rely on external emulation libraries (apart from `github.com/go-gl/gl/v3.2-core/gl` for the graphic renderer and  `github.com/go-gl/glfw/v3.3/glfw` for the window), and most of the code is written from scratch.

## 2. High-Level Architecture

The g64 emulator can be conceptually divided into the following major components:

*   **CPU (6510):** Emulates the MOS 6510 CPU, including instruction execution, register management, flag manipulation, and interrupt handling.
*   **VIC-II (Video Interface Chip):** Emulates the VIC-II graphics chip, responsible for generating the video output.
*   **SID (Sound Interface Device):** Emulates the SID sound chip, responsible for generating audio.
*   **CIA 1 & 2 (Complex Interface Adapter):** Emulates the two CIA chips, which handle various I/O tasks, including timers, keyboard input, and serial communication.
*   **PLA (Programmable Logic Array):** Manages the C64's memory map, including RAM, ROM, and I/O regions. Handles bank switching, memory access, I/O ports, cartridge access, color management, and emulator identification.
*   **Memory:** Provides an interface for memory access.
*   **Cartridge:** Emulates different types of cartridges (ROM, RAM expansions, and more complex cartridges like the EasyFlash).
*   **Disk Drive (1541):** Provides a full emulation of the 1541 disk drive, including its own 6502 CPU, VIA chips, and the drive mechanics.
*   **Input:** Handles keyboard and joystick input.
*   **Renderer:** Provides different rendering options (currently ASCII and OpenGL).
*   **Board:** Acts as the "motherboard," connecting all the components.

## 3. Component Details

### 3.1. CPU (src/components/mos6510)

*   **Implementation Approach:** Micro-operation based. Each 6502/6510 instruction is broken down into a sequence of smaller operations, each corresponding (approximately) to a CPU clock cycle.
*   **Dispatch Tables:** Uses dispatch tables (`opcodes.go`) to map opcodes to the corresponding instruction implementation functions and addressing mode functions.
*   **Files:**
    *   `cpu.go`: Main CPU struct, execution loop (`Emulate`), interrupt handling, register access, memory access (via the `IBanks` interface).
    *   `instructions.go`: Declarations of all instruction implementation functions.
    *   `inst_*.go`: Implementation of individual instructions, grouped by category (e.g., `inst_load_store.go`, `inst_arithmetic.go`, etc.).
    *   `opcodes.go`: Dispatch tables (opTable, modeTable).
    *   `stack.go`: Stack-related operations.
    *   `utils.go`: Utility functions.
*   **Interfaces:**
    *   `IBanks`: Used to access memory.
    *   `IPic`: Used to interact with the interrupt controller.

### 3.2. VIC-II (src/components/vic)

*   **Implementation Approach:** Cycle-accurate (or very close to it) emulation of the VIC-II. The `Emulate` function drives the emulation, stepping through each clock cycle of each raster line.
*   **Raster Timing:** Accurate emulation of raster timing, including badlines (where the VIC-II steals memory cycles from the CPU) and raster interrupts.
*   **Graphics Modes:** Supports standard text mode, multicolor text mode, standard bitmap mode, multicolor bitmap mode, and extended background color mode.
*   **Sprites:** Supports all 8 hardware sprites, including collision detection, priority, and expansion.
*   **Files:**
    *   `vic.go`: Main VIC struct, `Emulate` method, register access, interrupt handling.
    *   `graphics.go`: Functions for rendering the different graphics modes.
    *   `sprites.go`: Functions for managing and rendering sprites.
    *   `borders.go`: Functions for rendering the border.
    *   `pal.go`: Constants for PAL timing. Defines a series of functions, one for each cycle of a scanline (`palCycle1` to `palCycle63`), which are called sequentially by `Emulate`.
    *   `tables.go`: Lookup tables.

### 3.3. SID (src/components/sid)

*   **Implementation:** Emulation of the SID sound chip, with support for multiple voices, waveforms, filters, and envelopes.
*   **Files:**
    *   `sid.go`: Main SID struct, methods for managing sound generation and register access.
    *   `envelope.go`: Manages the amplitude of the voices.
    *   `filter.go`: Manages the filter.
    *   `oscillator.go`: Creates the waveforms.
    *   `voice.go`: Structure of the SID voices.
    *   `wave.go`: Implements the different waveforms.

### 3.4. CIA (src/components/cia)

*   **Implementation:** Emulates the two 6526 CIA chips.
*   **Functionality:** Handles timers, time-of-day clock, serial port, and keyboard/joystick input (CIA1).
*   **Files:**
    *   `mos6526.go`: Main CIA struct and methods.
    *   `timer.go`: Timer implementation.
    *   `tod.go`: Time-of-Day clock implementation.

### 3.5. PLA (src/c64/pla)

*   **Functionality:** Manages the C64's memory map, including RAM, ROM, and I/O regions. Handles bank switching, memory access, I/O ports, cartridge access, color management, and emulator identification.
*   **Files:**
    *   `pla.go`: Main PLA struct, memory access methods, port handling.
    *   `memorymap.go`: Defines the different memory configurations (bank layouts) possible on the C64.
    *   `ports.go`: Handles the 6510's I/O ports (primarily for controlling memory banking).
    * `socket.go`: Defines the interfaces for external connection.
    * `emulatorid.go`: Manages the emulator id information.
    * `writetrigger.go`: Manages the triggers to the memory.
* **Interfaces:**
    *   `IExpansionSocket`: Used to access to the cartridges.
    *	`ISocket`: Used to access to the registers.
    * `IRomSocket`: Used to access the system roms.

### 3.6. Memory (src/memory)

*   **Functionality:** Provides the interface `Memory` for accessing memory. The `pla` component handles the memory mapping.
*   **Files:**
    *   `memory.go`: Defines the `Memory` interface, providing a generic way to access memory.
    *   `memory_c64.go`: Implements the `Memory` interface for the C64, including the specific memory map.

### 3.7. Cartridges (src/c64/cartridges)

*   **Functionality:** Manages the loading and switching of different cartridge types.
*   **Files:**
    *   `manager.go`: Manages the cartridges.
    *   `icartridge/icartridge.go`: Defines the `ICartridge` interface, which all cartridge implementations must satisfy.
    *   Subdirectories: Contains implementations for specific cartridge types (e.g., `easyflash/`, `reu/`).

### 3.8. Input (src/c64/inputs)

*   **Keyboard (keyboard.go):**
    *   Buffers keyboard events (key presses and releases) in a FIFO queue.
    *   Converts keyboard input (virtual keys, ASCII characters) into a C64-compatible format.
    *   Provides methods for polling the next keyboard event from the queue.
*   **Joystick (joystick.go):** Manages the joysticks state.

### 3.9. Board (src/c64/board)

*   **board.go:** Represents the C64 motherboard. Connects all the components together (CPU, VIC-II, SID, CIA, memory, expansion port).
*   **Sockets:** Uses "socket" structs (`CPUSocket`, `VicSocket`, `CIA1Socket`, `CIA2Socket`, `SidSocket`) to provide an abstraction layer between the Board and the individual components, this implement the `ISocket` interface.
*   **Clock:** The `Board` uses a `Quartz` component to manage the clock cycles.

### 3.10. Disk Drive 1541 (src/c1541)

*   **Functionality:** Emulates the Commodore 1541 disk drive completely, including its own 6502 CPU, two VIA 6522 chips, RAM, ROM, and the drive mechanics.
*   **Files:**
    *   `c1541.go`: Main file.
    *   `board/board.go`: Main `C1541Board` struct.
    *   `board/cpusocket.go`: The socket for the cpu.
    *   `board/via1socket.go`: The socket for the first VIA.
    *   `board/via2socket.go`: The socket for the second VIA.
    *   `mechanic/mechanic.go`: Manages the drive mechanic.
    *   `mechanic/factory.go`: Manages the disk images.
    *   `disk/void/void.go`: Empty disk.
    *   `disk/gcr/tracks.go`: Implements gcr management.
    *   `disk/disk.go`: Handles disk file.
    *   `disk/track.go`: Handles disk tracks.
    *   `disk/conv.go`: Contains conversion function.
    *   `banks/banks.go`: Manages banks.
    *   `banks/jiffy.go`: Implements jiffy dos.
    *   `banks/loader.go`: Load the disk roms.
    *   `banks/builtin.go`: load the builtin rom.

### 3.11. Renderers (src/asciirender, src/glrender, src/pixels)

*   **asciirender (src/asciirender):** Renderer for text mode.
*   **glrender (src/glrender):** Renderer for openGL graphic mode.
*   **pixels (src/pixels):** used for managing pixels.

### 3.12. Utils

*   **fifo (src/fifo):** implement the fifo queue.
*   **filler (src/filler):** implement the filler functions.
*   **signals (src/signals):** implement signals for inter components management.

### 4. Data Flow (Simplified)

1.  **Initialization:** The `main.go` file initializes the system, creating instances of the `Board`, the renderer, and other core components.
2.  **Main Loop:** The main loop (typically within the renderer) runs continuously.
3.  **Input:** Keyboard and joystick input are captured by the `input` component.
4.  **Event Transmission:** The keyboard events are stored in the `keyboard` component into a FIFO queue.
5.  **Event Read:** The `CIA` reads the events from the `keyboard` component.
6.  **Clock:** the `Board` use `Quartz` component to manage the cycles.
7.  **Emulation Cycle:** The `Board.Emulate()` method is called. This triggers a cascade of calls:
    *   The VIC-II `Emulate` method is called (multiple times, once per cycle). The VIC-II performs its operations for the current cycle, including accessing memory, updating its internal state, and generating video output. The `pal.go` file defines the sequence of operations for each cycle.
    *   The CIA `Emulate` methods are called.
    *   The SID `Emulate` method is called.
    *   The CPU `Emulate` method is called (or, more accurately, the appropriate micro-operation for the current instruction is executed).
    * The `PLA` component handles the memory and ports management.
    * The `Board` component interconnects all the other components.
8.  **Output:** The renderer draws the current frame to the screen, using the data provided by the VIC-II.
9.  **Repeat:** The loop continues, emulating the next frame.
10. **Interrupts**: The PIC manages the interrupts for the various components.
11. **Matrice tastiera:** La matrice della tastiera è gestita dal CIA.

### 5. Key Design Decisions

*   **Micro-operations (CPU):** The CPU emulation is based on micro-operations, which allows for cycle-accurate emulation.
*   **Cycle-Accurate VIC-II:** The VIC-II emulation aims for cycle accuracy, using separate functions for each clock cycle of a scanline.
*   **Interfaces:** Extensive use of interfaces for decoupling.
*   **No External Emulation Libraries:** The emulator is written from scratch, giving the developer complete control over the implementation.
* **PLA:** The PLA component handles the memory map management.

### 6. Error Handling

The emulator handle the error by the function `log.Printf`. In the future an appropriate error handling will be added.

### 7. Configuration

The emulator use the `config.Config` to manage the configuration. The configuration can be changed by cli parameters or by file.

### 8. External Dependencies

*   `github.com/go-gl/gl/v3.2-core/gl`
*   `github.com/go-gl/glfw/v3.3/glfw`