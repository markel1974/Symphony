# C64 PLA (Programmable Logic Array) Emulator in Go

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/markel1974/symphony)](https://goreportcard.com/report/github.com/markel1974/symphony)

This project provides a high-performance, high-fidelity emulation of the Programmable Logic Array (PLA) chip logic for the Commodore 64, written entirely in Go. The PLA acts as the central "traffic cop" of the C64's memory map, arbitrating access between the CPU, RAM, ROM, and I/O chips. This implementation is designed for accuracy, performance, and modular integration into emulation frameworks like Symphony.

## Key Features

This emulator implements the complex logic of the C64's memory management system with a strong focus on both performance and fidelity.

* **Complete Memory Mapping**: Fully implements all 32 primary memory configurations of the C64, correctly mapping RAM, KERNAL ROM, BASIC ROM, Character ROM, Cartridge ROM, and I/O regions.
* **Dynamic Reconfiguration**: The memory map is reconfigured at runtime in response to changes in the 6510 CPU's on-chip I/O port (`$0001`), precisely modeling how the `/LORAM`, `/HIRAM`, and `/CHAREN` lines control the system's memory layout.
* **Accurate Cartridge Handling**: Tightly integrates with a cartridge management system, respecting the `/GAME` and `/EXROM` lines to handle various memory mapping schemes, including standard 8K, 16K, and Ultimax modes.
* **Faithful "Write-Through" Behavior**: Correctly models the C64's behavior where writes to addresses occupied by ROM are passed through to the underlying RAM.
* **Accurate Open Bus Reads**: Reads from unmapped memory regions (`UND`) correctly return the last byte read by the VIC-II chip, a critical detail for advanced demo and game compatibility.
* **`WriteTriggers` Mechanism**: Includes a powerful and flexible system for hooking into memory write operations at specific addresses, enabling advanced debugging, introspection, or the implementation of complex hardware features.

## Architecture and Design

The PLA emulator is built on modern software design principles to ensure clarity, maintainability, and peak performance.

* **Dispatch Table Pattern**: The core of the memory access mechanism is a pair of dispatch tables (`bankRead` and `bankWrite`). These are arrays of function pointers. Instead of using `switch` statements in the critical `Read`/`Write` path, the PLA performs a direct array lookup and an indirect function call. This is a significant performance optimization that eliminates conditional branching in the hot path.
* **Data-Driven Configuration**: The `memorymap.go` file contains a static "truth table" (`_memoryMap`) for all 32 memory configurations. The `applyMemoryConfig` function reads this data and "compiles" it into the live `bankRead`/`bankWrite` dispatch tables at runtime, but only when a configuration change is necessary.
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
2.  **Binding**: The `Board` component calls the PLA's `Bind` method, supplying it with interfaces to all the other necessary components (RAM, ROMs, Cartridge Manager, VIC-II signals, I/O chips). The PLA stores these as function pointers for its dispatch tables.
3.  **Execution**: The emulated CPU and VIC-II call the PLA's `Read(address)` and `Write(address, value)` methods for every memory access.
4.  **Reconfiguration**: When the CPU writes to `$0001` or a cartridge changes its banking state, the `applyMemoryConfig` method is triggered to rebuild the dispatch tables based on the new system-wide state.

## Dependencies

This component is designed to be part of the Symphony framework and relies on its core packages:
* `github.com/markel1974/symphony/src/component`
* `github.com/markel1974/symphony/src/config`
* `github.com/markel1974/symphony/src/references`
* `github.com/markel1974/symphony/src/registry`