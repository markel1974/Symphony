# Symphony - A Configurable Emulation Platform

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Overview

`g64` is a Commodore 64 computer emulator written entirely in Go. The project was developed as a personal challenge; no external libraries were used for the emulation itself. The emulator aims to faithfully reproduce the behavior of the original hardware, with a focus on cycle-accurate emulation of the 6510 CPU and the VIC-II graphics chip.

**Project Status:**

`g64` is a personal project under active development. It is *functional* and capable of running the *vast majority* of C64 software, including complex games and demos.  The project currently *lacks* unit and integration tests. The documentation is *incomplete*.

**DO NOT USE IN PRODUCTION OR WITH UNTRUSTED SOFTWARE.**

## Features

*   Accurate emulation of the MOS 6510 CPU (a 6502 variant).
*   Accurate emulation of the VIC-II graphics chip (text mode, bitmap, multicolor, sprites, raster interrupts, badlines, etc.).
*   Emulation of the SID sound chip (waveforms, filters, envelopes).
*   Emulation of the CIA 6526 chips (timers, TOD, I/O).
*   Memory management (RAM, ROM, banking).
*   Support for cartridges (via a modular expansion system):
    *   EasyFlash
    *   REU
    *   ... (others, to be listed)
*   Emulation of the 1541 disk drive (complete, including the CPU).
*   Keyboard and joystick input handling.
*   ASCII renderer (textual output in the terminal).
*   OpenGL renderer.
*   Copy and paste support.

## Architecture

`g64` is designed with a modular architecture, with separate components for the different parts of the emulated system.

**Main Components:**

*   **`cpu/`:**  6510 CPU emulator.
    *   `cpu.go`:  Main CPU logic, execution cycle, register/flag/interrupt handling.
    *   `instructions.go`: Declarations of functions for instructions.
    *   `inst_*.go`: Implementation of 6502/6510 instructions (divided by category).
    *   `opcodes.go`: Dispatch tables (opcode -> function).
    *   `stack.go`: Stack management.
    *   `utils.go`: Utility functions.
*   **`memory/`:** Memory management.
    *   `memory.go`:  `Memory` interface.
    *   `memory_c64.go`: C64-specific implementation.
*   **`banks/`:**  Memory banking and component access (RAM, ROM, I/O) management.
    *   `banks.go`: Main banking logic.
    *   `memorymap.go`: Definitions of memory configurations.
    *   `ports.go`: Management of 6510 I/O ports.
    *   `observer.go`:  Memory observation/modification functionality.
*   **`cartridges/`:** Cartridge management.
    *   `manager.go`: Cartridge manager.
    *   `icartridge/`: `ICartridge` interface.
    *   Subdirectories for specific implementations (e.g., `easyflash/`).
*   **`components/`**:
    *   `vic/`: VIC-II emulator.
        *    `vic.go`: Main logic.
        *   `character_mode.go`: Management of the character mode.
        *   `graphic_mode.go`: Management of the graphics mode.
        *   `irq.go`: Interrupt management.
        *   `memory.go`: Memory Management.
        *   `raster.go`: Raster management.
        *   `registers.go`: Register Management.
        *   `sprite.go`: Sprite management
    *   `sid/`: SID emulator.
    *   `cia/`: CIA chip emulators.
*   **`pixels`:**  *Internal* library for text-based rendering.
*   **`board`:** C64 "motherboard". Connects the various components.
*   **`main.go`:** Starts the program.
*   **`asciirender`:**  ASCII renderer.
*   **`glrender`:** OpenGL renderer.

**Interfaces:**

`g64` makes extensive use of Go interfaces to decouple components. Key interfaces include:

*   `memory.Memory`:  For memory access.
*   `mos6510.IBanks`: For accessing memory banks.
*   `mos6510.IPic`:  For interrupt management.
*   `board.ISocket`: For interacting with chips (CPU, VIC-II, SID, CIA).
*   `icartridge.IExpansion`:  For interacting with cartridges.
* `icartridge.ICartridge`: For cartridge management.
*   `interfaces.ITerminal`

**External Dependencies:**

*  `github.com/go-gl/gl/v3.2-core/gl`
*  `github.com/go-gl/glfw/v3.3/glfw`

**[TODO: Add an architecture diagram, if possible.]**

## Installation

```bash
go get [github.com/markel1974/g64](https://github.com/markel1974/g64)