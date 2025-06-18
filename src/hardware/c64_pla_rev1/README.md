
# C64 PLA (Programmable Logic Array) Emulator in Go

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/markel1974/symphony)](https://goreportcard.com/report/github.com/markel1974/symphony)

This project is a high-performance, high-fidelity emulation of the Programmable Logic Array (PLA) chip used in the Commodore 64, written entirely in Go. The PLA is the central "traffic cop" of the C64's memory map, arbitrating access between the CPU, RAM, and various ROM and I/O chips. This implementation is designed for accuracy, performance, and modular integration into emulation frameworks like Symphony.

## Key Features

This emulator implements the complex logic of the C64's memory management system with a focus on fidelity and performance.

* **Accurate Bank Switching**: Fully implements all 32 memory configurations of the C64, correctly mapping RAM, KERNAL ROM, BASIC ROM, Character ROM, Cartridge ROM, and I/O regions based on the system's state.
* **Dynamic Reconfiguration**: The memory map is reconfigured at runtime in response to changes in the 6510 CPU's on-chip I/O port (`$0001`), precisely modeling how the `/LORAM`, `/HIRAM`, and `/CHAREN` lines control the system's memory layout.
* **Cartridge Support Integration**: Tightly integrates with a cartridge management system, respecting the `/GAME` and `/EXROM` lines to correctly handle various cartridge memory mapping schemes, including `ROL` (ROM Low) and `ROH` (ROM High).
* **Faithful I/O Decoding**: Accurately simulates the behavior of the `74LS138` decoder (U15) for the I/O block at `$D000-$DFFF`, correctly routing memory accesses to the VIC-II, SID, and CIA chips.
* **`WriteTriggers` Mechanism**: Includes a powerful and flexible system for hooking into memory write operations at specific addresses, enabling advanced debugging, introspection, or the implementation of complex hardware features without modifying the core PLA logic.

## Architecture and Design

The PLA emulator is built on the same core principles that define the Symphony framework, prioritizing architectural clarity and performance.

* **High-Fidelity Behavioral Model**: The approach is to replicate the observable *results* of the PLA's logic with perfect accuracy. The configuration is driven by a comprehensive "truth table" that defines the output for every possible state.
* **Dispatch Table Pattern**: The core of the memory access mechanism is a pair of dispatch tables (`bankRead` and `bankWrite`). These are arrays of function pointers. Instead of using a `switch` statement in the critical `Read`/`Write` path, the PLA performs a direct lookup and an indirect call to the correct handler function for the targeted memory bank. This is a significant performance optimization that eliminates branching.
* **Data-Driven Configuration**: The `memorymap.go` file contains the static `_memoryMap` array, which acts as the definitive source of truth for all 32 memory configurations. The `applyMemoryConfig` function uses this data to "compile" the `bankRead`/`bankWrite` dispatch tables at runtime, only when a configuration change is necessary.
* **Component-Based ("Headless") Design**: As shown by `factory.go`, this package is designed as a self-contained, "headless" component. Its logic is entirely decoupled from other parts of the system, with which it communicates through well-defined interfaces and sockets.

## File Structure

* `pla.go`: Contains the core `PLA` component, the dispatch tables, and the main `Read`/`Write` logic.
* `ports.go`: Implements the emulation of the 6510 CPU's on-chip I/O port at memory locations `$0000` and `$0001`, which directly controls the PLA's state.
* `memorymap.go`: Defines the static data for all 32 C64 memory configurations.
* `writetrigger.go`: Implements the optional write-hooking system.
* `factory.go`: Implements the factory pattern for integration into the Symphony framework.

## How to Use (Integration)

Within the Symphony framework, the PLA component has a clear lifecycle:
1.  **Creation**: An instance is created by the `ComponentFactory` using `NewPLA()`.
2.  **Binding**: The `Board` component calls the PLA's `Bind` method, supplying it with interfaces to all the other necessary components (RAM, ROMs, Cartridge Manager, I/O chips). The PLA stores these as function pointers for its dispatch tables.
3.  **Execution**: The emulated CPU calls the PLA's `Read(address)` and `Write(address, value)` methods for every memory access. These methods perform the fast dispatch to the correct memory handler.
4.  **Reconfiguration**: When the CPU writes to `$0001`, the PLA's `update` and `RebuildMemoryConfig` methods are triggered to re-calculate the memory map and update the `bankRead`/`bankWrite` dispatch tables accordingly.

## Dependencies

* `github.com/markel1974/symphony/src/component`
* `github.com/markel1974/symphony/src/config`
* `github.com/markel1974/symphony/src/references`
* `github.com/markel1974/symphony/src/registry`