# Symphony: A Configurable and Introspectable Emulation Framework

> Symphony is an open-source, highly configurable, and deeply introspectable emulation framework written entirely in Go, designed for accuracy, extensibility, and the dynamic exploration of computer systems.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/markel1974/symphony)](https://goreportcard.com/report/github.com/markel1974/symphony) ---

## Beyond Emulation: A Framework for Creation and Exploration

Symphony is fundamentally different from most emulators. It's an **emulation framework** designed with **modularity, runtime configuration, and deep introspection** as first-class citizens. If you've ever wanted to not just *run* classic systems, but truly *understand*, *modify*, and *experiment* with their architecture in real-time, Symphony provides the platform.

Imagine dynamically connecting virtual hardware components, live-patching memory, inspecting CPU registers or chip state mid-execution, saving the *entire machine definition and state* into a single file, and controlling everything through a powerful SSH console with its own windowing system. This is the core vision realized in Symphony.

Currently featuring a highly accurate Commodore 64 and 1541 disk drive implementation, Symphony's core architecture is system-agnostic and ready to be extended to other platforms.

## What Makes Symphony Unique?

* **Deep Dynamic Introspection:** Interact with the running emulation at an unprecedented level. Access and modify the state of *any* component (CPU registers, memory banks, I/O chip state, custom properties) at **runtime** using a simple, hierarchical path system (e.g., `get c64.cpu.pc`, `set c64.vic.border_color 1`) via the integrated console or potentially external tools. Components can also register custom commands.
* **Snapshot = Configuration = State:** Symphony eliminates traditional config files. A single **snapshot** (a `map[string]interface{}`, typically stored as JSON) defines the **entire hardware configuration** (component types, tree structure) *and* the **runtime state**. Loading a snapshot *builds* (or restores) the complete emulated machine. This allows for extreme flexibility in defining, saving, loading, and sharing machine configurations.
* **Truly Modular Architecture:** Built entirely around the `references.IComponent` interface and `component.BaseComponent` embedding. Every element is a component in a hierarchical tree. Components interact *only* through Go interfaces (including specialized `Socket` interfaces), ensuring maximum **decoupling, testability, and extensibility**. Swap CPU models, add custom hardware, or build entirely new systems by implementing the required interfaces.
* **Clean 3-Phase Initialization:** Component interdependencies are managed elegantly:
    1.  **Create/Restore (`component.RestoreAll`):** Builds the entire component tree from the snapshot, creating components via the factory and restoring their properties.
    2.  **Socket Mounting (`Board` -> `ISocket.Mount`):** Each component's "socket" (mediator object) finds its required dependencies by navigating the tree from its parent.
    3.  **Component Binding/Setup (`ISocket.Mount` -> `IHardware.Bind`/`Setup`):** Sockets establish the final connections, typically by calling a `Bind` or `Setup` method on the component itself.
* **Dynamic Component Factory:** A central factory (`hardware.Factory`) uses a **dynamic registry** (`registry` package) populated via package `init()` functions. Adding support for new component types requires only creating the component and its factory and registering it – no changes needed to the central factory. The factory also serves required external services like `IDisplayBuffer` and `IAudioRender`.
* **Advanced Integrated Console:** A powerful, **text-based console environment** accessible via **SSH**:
    * Multi-Window & Multi-Process support within the terminal session.
    * Full access to the introspection system (`get`, `set`).
    * Execution of global and component-specific custom commands.
    * Built-in tools like performance graphs (e.g., GC visualization).
    * Standard shell features (`cd`, `ls`, `pwd`, `history`, process management).
* **Fully Headless Operation:** The core Symphony emulation engine runs **completely headless**. Graphics (`IDisplayBuffer`) and audio (`IAudioRender`) are optional external dependencies injected during setup, making Symphony ideal for automated testing, server tasks, or custom frontends.
* **Accuracy-Focused Core:** Emulation components (like the C64 VIC-II) are designed with cycle-level accuracy in mind where appropriate.
* **Pure Go Core:** Ensures **portability**, **memory safety**, and easier development with **no CGo dependencies** for the core emulation logic.

## Architecture Highlights

* **Interfaces:** `references.IComponent` (combined `IHardware`, `INavigate`, `ICommand`), `references.IFactory`, `references.ISocket`, plus specific hardware interfaces (`I6510`, `IVIC`, `ISID`, `ICIA`, `IVIA`, `IPlaC64`, `IIecDevice`, `IIecProtocolDevice`, `ICartridge`, etc.).
* **Core Components:** `component.BaseComponent`, `component.Node`, `hardware.Factory`, `registry`, `iec.Dispatcher`, `iec.Protocol`, `quartz.Quartz`.
* **Key Principle:** Strong decoupling and dynamic configuration managed through the component tree and snapshot data.

## Current Implementation: Commodore 64 / 1541 Example

Symphony includes a highly accurate implementation serving as a proof-of-concept:

* **Commodore 64:** 6510 CPU (micro-op level), cycle-accurate 6569 VIC-II (PAL), 6581 SID, 2x 6526 CIA, C64 PLA, KERNAL/BASIC/Char ROMs (including JiffyDOS option), Keyboard, Joysticks.
* **Commodore 1541 Disk Drive:** Full drive emulation including 6502 CPU, 2x 6522 VIA, Drive Mechanics emulation, 1541 PLA, DOS ROM (standard and JiffyDOS options), GCR support.
* **`MediaDrive` Component:** A virtual IEC device using `iec.Protocol` and an `adapters.IAdapter` interface to map to:
    * Host Filesystem Directories (`adapters.Directory`)
    * ZIP Archives (`adapters.Zip`)
    * Single Host Files (`adapters.File` - for testing/shortcuts)
* **Cartridges:** Support for standard formats, EasyFlash, REU (various sizes), Magic Desk, Ocean Type 1. Includes CRT file loader.

This implementation successfully runs complex software like GEOS and demanding games/demos, showcasing the framework's capabilities.

## Getting Started

[**TODO:** Add detailed build instructions (Go version, dependencies for optional renderers like GL) and basic run commands.]

```bash
# Example: Build
go build -o symphony ./cmd/symphony # Assuming main is in cmd/symphony

# Example: Run default C64 config using RestoreAll from embedded default snapshot
# (Requires main.go to be adapted to use RestoreAll instead of manual creation)
# ./symphony -r gl # Use OpenGL renderer (needs GLFW deps)

# Example: Run C64 with a specific PRG using the file adapter shortcut
# ./symphony -r gl -m c64 -d 8:/path/to/your/game.prg

# Example: Run C64 with a specific EasyFlash CRT
# ./symphony -r gl -m c64 -c easyflash:/path/to/your/cart.crt

# Example: Run headless
# ./symphony -r none -a none -m c64