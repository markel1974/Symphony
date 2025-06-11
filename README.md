# Symphony: A Configurable Emulation Framework

> Symphony is one of the first open-source, highly configurable, and deeply introspectable emulation frameworks written entirely in Go, designed for accuracy, extensibility, and the dynamic exploration of computer systems.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/markel1974/symphony)](https://goreportcard.com/report/github.com/markel1974/symphony) ---

**Symphony is not just another emulator. It's an exploration.**

Born from decades of low-level development experience and a frustration with the limitations of traditional, static emulator designs, Symphony was built with a different philosophy. It's a **framework** designed for those who don't just want to *run* old software, but want to *understand*, *modify*, and *experiment* with the underlying hardware architecture itself.

Imagine dynamically rewiring components, inspecting CPU registers mid-flight, saving the *entire* state and *configuration* of a complex machine into a single file, and doing it all through a powerful, scriptable console accessible from anywhere via SSH. **That's the core idea of Symphony.**

While it currently boasts a highly accurate Commodore 64 and 1541 disk drive implementation as a proof-of-concept, Symphony's modular architecture is inherently **system-agnostic**, ready to be extended to emulate a wide variety of machines.

## What Makes Symphony Unique?

* **Deep Dynamic Introspection:** Inspect and modify the state of *any* component at **runtime** via a powerful built-in console (accessible via SSH!). Change registers, memory, or custom properties on the fly.
* **Snapshot = Configuration = State:** A single snapshot defines the **entire hardware configuration** *and* runtime state. Build, save, load, and share complex machine setups easily.
* **Truly Modular Architecture:** Built entirely around the `IComponent` interface. Components communicate *only* through well-defined interfaces. Easily add new hardware, swap implementations (e.g., different SID versions), or even build entirely new systems.
* **Advanced Integrated Console:** Access via **SSH**, featuring **multiple windows**, **concurrent processes**, custom commands per component, and even basic **text-mode graphing**. It's a complete debugging and experimentation environment.
* **Fully Headless Operation:** The core Symphony emulation engine runs **completely headless**, without requiring any graphical or audio output. Ideal for automated testing, server-side emulation tasks, or integration with custom frontends.
* **Accuracy-Focused:** Where implemented (like the C64 VIC-II), aims for **cycle-level accuracy**.
* **Pure Go Core:** Ensures **portability** and **memory safety** with **no CGo dependencies** for the core logic.
* **Extensibility First:** Designed as a framework from the ground up.

## Architecture Highlights

Symphony's power stems from its carefully designed architecture:

1.  **The Component Tree (`IComponent`, `BaseComponent`, `Node`):** Everything is a component implementing `IComponent` and embedding `BaseComponent`. These are organized hierarchically in a tree managed by `Node` objects, providing universal methods for identification, navigation, properties, commands, and state management.
2.  **Snapshot-Driven Initialization (`RestoreAll`/`_restore`):** The static `component.RestoreAll` function recursively builds the component tree *and* restores component state directly from a snapshot map (`map[string]interface{}`). It uses the `ComponentFactory` to instantiate components based on "type" information within the snapshot (unless marked "internal").
3.  **Component Factory (`hardware.Factory`, `IFactory`, `registry`):** A central factory uses a dynamic registry (populated via package `init()` functions) to delegate creation to specific component factories (`IFactory` implementations). The factory also serves as a provider for global services (`IDisplayBuffer`, `IAudioRender`).
4.  **Sockets (`ISocket` & Implementations):** Lightweight structs that act as typed intermediaries, primarily used during the connection phase. They implement `ISocket`'s `Mount` method.
5.  **3-Phase Connection (`RestoreAll` -> `Board.Setup` -> `Board.Connect` -> Socket `Mount`):** Initialization occurs in distinct phases after `RestoreAll` builds the tree: `Board.Setup` orchestrates the `socket.Mount()` call for all relevant sockets. Inside `Mount`, the socket finds its dependencies by navigating the tree via its parent (`IComponent`) reference, stores the required interfaces, and calls the component's final `Setup` (`IHardware.Setup`) or `Bind` method to establish the connection. *(Self-correction on phases needed here)* -> *Correction*: **Revised 3-Phase:** 1. `RestoreAll` (Create Tree + Restore State). 2. `Board.Setup` orchestrates `socket.Mount` (which finds dependencies AND likely calls component `Setup`/`Bind`). 3. `Board.Connect` (if still needed, maybe just calls component `Connect`?). [**TODO:** Re-verify exact final sequence of Setup/Connect/Mount calls based on latest code].

## Current Implementation: Commodore 64 / 1541

[... as before ...]

## Getting Started

[**TODO:** Add detailed build and run instructions. Mention headless execution.]

## Live Evaluation
https://markel1974.itch.io/symphony

```bash
# Example: Run default C64 config with OpenGL renderer
./symphony -r gl -m c64

# Example: Run headless (no graphics/audio)
./symphony -r none -a none -m c64