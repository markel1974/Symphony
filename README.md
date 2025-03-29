# Symphony: A Configurable Emulation Framework

> Symphony is one of the first open-source, highly configurable, and deeply introspectable emulation frameworks written entirely in Go, designed for accuracy, extensibility, and the dynamic exploration of computer systems.

## Beyond Emulation: A Framework for Creation

Symphony isn't just another emulator designed to run software for a single machine. It's a foundational **framework** built from the ground up with **modularity, dynamic configuration, and deep introspection** at its core. Born from a desire to overcome the limitations of traditional, static emulator designs, Symphony provides a unique platform for accurately simulating hardware, experimenting with different configurations, debugging complex interactions, and learning about computer architecture in an interactive way.

While it currently features a highly accurate Commodore 64 and 1541 disk drive implementation as a proof-of-concept, Symphony's architecture is inherently **system-agnostic**, ready to be extended to emulate a wide variety of machines.

## What Makes Symphony Unique?

* **Dynamic Introspection:** This is Symphony's cornerstone. Access and modify the state of *any* component (CPU registers, memory, chip state, custom properties) at **runtime** using a simple path-based system (e.g., `get c64.cpu.pc`, `set c64.vic.border_color 1`). This is facilitated by a powerful, built-in interactive console.
* **Snapshot-Driven Configuration = State:** Forget separate config files. Symphony treats a **snapshot** as the **single source of truth** for both the **hardware configuration** and the **runtime state**. Loading a snapshot *builds* the entire component tree and restores its state. Modifying the snapshot *is* modifying the configuration. This enables incredible flexibility and easy sharing of complex setups.
* **Truly Modular Architecture:** Built entirely around the `IComponent` interface and `BaseComponent` embedding. Every element, from the main board to internal timers, is a component within a hierarchical tree. Components communicate *only* through well-defined interfaces (often via specialized "Sockets").
* **Clean Dependency Management:** A unique 3-phase initialization process (Create/Restore -> Socket Setup -> Component Connect), orchestrated by the framework, ensures that component interdependencies are resolved gracefully and automatically, eliminating complex manual wiring or order dependencies in the configuration.
* **Advanced Debug Console:** Go beyond simple command lines. Symphony features an integrated console accessible via **SSH**, supporting **multiple windows**, **multiple concurrent processes** (run debuggers alongside games!), custom commands per component, and even basic **text-mode graphing** for real-time metrics.
* **Pure Go Implementation:** Written entirely in Go for performance, safety, concurrency potential, and portability. The core emulation has **no CGo dependencies**.
* **Extensibility First:** Designed from the start as a framework. Adding support for new hardware components (different CPUs, video chips, peripherals) or entire new systems primarily involves implementing the `IComponent` interface, creating a specific component factory, and registering it.

## Architecture Highlights

Symphony's power stems from its carefully designed architecture:

1.  **Component Tree (`IComponent`, `BaseComponent`, `Node`):** The entire emulated system is a tree of components. `IComponent` is the universal interface, `BaseComponent` provides common functionality (ID, node access, properties, commands, snapshotting), and `Node` manages the tree structure.
2.  **Snapshot/Configuration (`RestoreAll`/`_restore`):** The static `component.RestoreAll` function recursively builds the component tree *and* restores component state directly from a snapshot map (`map[string]interface{}`). It uses the `ComponentFactory` to instantiate components based on "type" information within the snapshot, but only if the component isn't marked "internal" (allowing parent components like CIA to create their own children like Timers).
3.  **Component Factory (`hardware.Factory`, `references.IFactory`, `registry`):** A central factory uses a dynamic registry (populated via `init()` functions in component packages) to find and delegate creation to specific component factories (`IFactory` implementations), enabling easy addition of new component types. The factory also acts as a provider for global services like `IDisplayBuffer` and `IAudioRender`.
4.  **Sockets (`Connector` Interface):** Lightweight structs associated with the `Board` that manage the *connection* between components. They implement the `Connector` interface with two key methods called during initialization *after* the tree is built:
    * `Setup(components map[string]IComponent, cfg *config.Config)`: Finds required component dependencies within the provided map (using `HardwareId`) and stores their interfaces internally. Also gets services like display/audio from the factory via the component it manages.
    * `Connect()`: Uses the stored references to make the final connection (often by calling the actual component's specific `Setup` method, passing the socket or required interfaces).
5.  **Optimized Emulation Loop (`Board.Emulate`, `rebuildEmulation`):** The `Board` builds an optimized slice of `Emulate` function pointers during initialization (`rebuildEmulation`), ensuring the main loop only calls active components directly, minimizing overhead.

## Current Implementation: Commodore 64 Example

Symphony includes a highly accurate implementation of the Commodore 64 and the 1541 disk drive, demonstrating the framework's capabilities:

* **CPU (6510):** Micro-operation level emulation.
* **VIC-II (6569 PAL):** Cycle-accurate emulation, handling raster effects, sprites, collisions, and various graphics modes.
* **SID (6581):** Sound chip emulation.
* **CIA (6526 x2):** Emulation of both CIA chips, including timers and TOD clock.
* **PLA:** Accurate memory mapping logic.
* **1541 Drive:** Full emulation including the 6502 CPU, 2x 6522 VIAs, ROM/RAM, drive mechanics, and low-level IEC protocol handling.
* **Cartridges:** Support for various cartridge types (including EasyFlash, REU).
* **FSDrive:** Virtual IEC drive mapping to the host filesystem (uses `iecprotocol.Protocol`).

## Getting Started

[**TODO:** Add detailed build and run instructions]

```bash
# Example build
go build -o symphony ./src

# Example run (OpenGL, C64 default)
./symphony -r gl -m c64