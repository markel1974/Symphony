# g64 - Commodore 64 Emulator in Go

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

g64 - Commodore 64 Emulator - Architectural Overview

1. Introduction

g64 is a Commodore 64 emulator written in Go. It aims for high accuracy, including cycle-accurate emulation of the VIC-II video chip and the 6510 CPU. The emulator is designed to be modular, with distinct components responsible for different aspects of the C64 hardware. It does not rely on external emulation libraries (apart from termbox-go for the ASCII renderer and Pixel for the graphic renderer), and most of the code is written from scratch.

2. High-Level Architecture

The g64 emulator can be conceptually divided into the following major components:

CPU (6510): Emulates the MOS 6510 CPU, including instruction execution, register management, flag manipulation, and interrupt handling.
VIC-II (Video Interface Chip): Emulates the VIC-II graphics chip, responsible for generating the video output.
SID (Sound Interface Device): Emulates the SID sound chip, responsible for generating audio.
CIA 1 & 2 (Complex Interface Adapter): Emulates the two CIA chips, which handle various I/O tasks, including timers, keyboard input, and serial communication.
Memory: Manages the C64's memory map, including RAM, ROM, and I/O regions. Handles bank switching.
Cartridge: Emulates different types of cartridges (ROM, RAM expansions, and more complex cartridges like the EasyFlash).
Disk Drive (1541): Provides a full emulation of the 1541 disk drive, including its own 6502 CPU, VIA chips, and the drive mechanics.
Input: Handles keyboard and joystick input.
Renderer: Provides different rendering options (currently ASCII and OpenGL).
Board: Acts as the "motherboard," connecting all the components.
3. Component Details

3.1. CPU (src/components/mos6510)

Implementation Approach: Micro-operation based. Each 6502/6510 instruction is broken down into a sequence of smaller operations, each corresponding (approximately) to a CPU clock cycle.
Dispatch Tables: Uses dispatch tables (tables.go) to map opcodes to the corresponding instruction implementation functions and addressing mode functions.
Files:
cpu.go: Main CPU struct, execution loop (Run), interrupt handling, register access, memory access (via the IBanks interface).
instructions.go: Declarations of all instruction implementation functions.
inst_*.go: Implementation of individual instructions, grouped by category (e.g., inst_load_store.go, inst_arithmetic.go, etc.).
tables.go: Dispatch tables (opTable, modeTable).
stack.go: Stack-related operations.
utils.go: Utility functions.
Interfaces:
IBanks: Used to access memory.
IPic: Used to interact with the interrupt controller.
3.2. VIC-II (src/components/vic)

Implementation Approach: Cycle-accurate (or very close to it) emulation of the VIC-II. The Emulate function drives the emulation, stepping through each clock cycle of each raster line.
Raster Timing: Accurate emulation of raster timing, including badlines (where the VIC-II steals memory cycles from the CPU) and raster interrupts.
Graphics Modes: Supports standard text mode, multicolor text mode, standard bitmap mode, multicolor bitmap mode, and extended background color mode.
Sprites: Supports all 8 hardware sprites, including collision detection, priority, and expansion.
Files:
vic.go: Main VIC struct, Emulate method, register access, interrupt handling.
graphics.go: Functions for rendering the different graphics modes.
sprites.go: Functions for managing and rendering sprites.
borders.go: Functions for rendering the border.
pal.go: Constants for PAL timing. Defines a series of functions, one for each cycle of a scanline (palCycle1 to palCycle63), which are called sequentially by Emulate.
tables.go: Lookup tables.
3.3. SID (src/components/sid)

Files:
sid.go:
envelope.go
filter.go
oscillator.go
voice.go
wave.go
(Further details on the SID implementation would be added here, after analysis.)

3.4. CIA (src/components/cia)

Implementation: Emulates the two 6526 CIA chips.
Functionality: Handles timers, time-of-day clock, serial port, and keyboard/joystick input (CIA1).
Files:
mos6526.go: Main CIA struct and methods.
timer.go: Timer implementation.
tod.go: Time-of-Day clock implementation.
3.5. Memory (memory/ and banks/)

memory/memory.go: Defines the Memory interface, providing a generic way to access memory.
memory/memory_c64.go: Implements the Memory interface for the C64, including the specific memory map.
banks/banks.go: Implements the memory banking logic, allowing different ROMs and RAM areas to be mapped into the 6510's address space.
banks/memorymap.go: Defines the different memory configurations (bank layouts) possible on the C64.
banks/ports.go: Handles the 6510's I/O ports (primarily for controlling memory banking).
banks/observer.go: Provides functionality for inspecting and modifying memory (used for debugging and program loading).
3.6. Cartridges (src/c64/cartridges)

manager.go: Manages the loading and switching of different cartridge types.
icartridge/icartridge.go: Defines the ICartridge interface, which all cartridge implementations must satisfy.
Subdirectories: Contains implementations for specific cartridge types (e.g., easyflash, reu).
3.7. Input (src/c64/inputs)

Manages the keyboard.
Manages the joysticks.
3.8. Board (src/c64/board)

board.go: Represents the C64 motherboard. Connects all the components together (CPU, VIC-II, SID, CIA, memory, expansion port).
Sockets: Uses "socket" structs (CPUSocket, VicSocket, CIA1Socket, CIA2Socket, SidSocket) to provide an abstraction layer between the Board and the individual components.
3.9. Disk Drive 1541 (src/c1541)

Emulates the Commodore 1541 disk drive completely, including its own 6502 CPU, two VIA 6522 chips, RAM, ROM, and the drive mechanics.
Files:
c1541.go
board/board.go:
board/cpusocket.go:
board/via1socket.go:
board/via2socket.go:
mechanic/mechanic.go *mechanic/factory.go
disk/void/void.go
disk/gcr/tracks.go, disk/disk.go, disk/track.go, disk/conv.go *banks/banks.go, banks/jiffy.go, banks/loader.go, banks/builtin.go
3.10. Renderers

asciirender: Renderer for text mode.
glrender: Renderer for openGL graphic mode.
4.  Data Flow (Simplified)

Initialization: The main.go file initializes the system, creating instances of the Board, the renderer, and other core components.
Main Loop: The main loop (typically within the renderer) runs continuously.
Input: Keyboard and joystick input are captured.
Emulation Cycle: The Board.Emulate() method is called. This triggers a cascade of calls:
The VIC-II Emulate method is called (multiple times, once per cycle). The VIC-II performs its operations for the current cycle, including accessing memory, updating its internal state, and generating video output. The pal.go file defines the sequence of operations for each cycle.
The CIA Emulate methods are called.
The CPU Emulate method is called (or, more accurately, the appropriate micro-operation for the current instruction is executed).
Output: The renderer draws the current frame to the screen, using the data provided by the VIC-II.
Repeat: The loop continues, emulating the next frame.
5. Key Design Decisions

Micro-operations (CPU): The CPU emulation is based on micro-operations, which allows for cycle-accurate emulation.
Cycle-Accurate VIC-II: The VIC-II emulation aims for cycle accuracy, using separate functions for each clock cycle of a scanline.
Interfaces: Extensive use of interfaces for decoupling.
No External Emulation Libraries: The emulator is written from scratch, giving the developer complete control over the implementation.
6.  Missing Pieces (in this high-level description):

Detailed description of the SID emulation.
Detailed description of the CIA emulation.
Detailed description of the 1541 emulation.
Error handling strategy.
Configuration options.

This document provides a high-level overview of the g64 emulator's architecture. 
Further details can be found by examining the source code of the individual components.