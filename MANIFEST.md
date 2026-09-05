
# Symphony: The WYSIWYG Emulation Framework

> Symphony is not just another emulator. It is a high-fidelity, interactive hardware laboratory designed from the ground up to do more than just run old software—it's built to deconstruct, understand, and compose virtual machines.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/markel1974/symphony)](https://goreportcard.com/report/github.com/markel1974/symphony)

---

## The Guiding Vision: A Hardware Sandbox

**What if, instead of just using an emulator, you could *build* a vintage computer by dragging and dropping components onto a virtual motherboard? What if you could connect a CPU, a video chip, RAM, and then press 'Play' to watch it come to life?**

This is the fantasy that has guided Symphony's development from day one.

Every architectural decision in this project is a direct answer to that "What if?". The result is a framework that is not merely an emulator, but a powerful, system-agnostic **emulation chassis** designed for the ultimate goal: a web-based, WYSIWYG hardware builder.

## From Vision to Architecture

To realize this dream, a traditional design was not an option. Symphony's architecture is the necessary foundation for its ambitious goal:

* **Modular Components (`IComponent`):** To be able to "drag" a chip, every piece of hardware had to be a self-contained, interchangeable software object.
* **Sockets & Wiring (`ISocket`):** To "connect" the chips, the framework needed an abstraction for the physical connectors on a motherboard, defining the "contract" for how components interact.
* **Snapshots as Blueprints:** To "save" and "load" a visually constructed machine, a format was needed to describe the entire hardware topology and its state. The snapshot system is precisely this: the blueprint of the user's creation.
* **Headless & WASM Core:** To "bring the machine to life" anywhere, especially on the web, the simulation engine had to be completely decoupled from its graphical representation.

## Key Features

* **Deep Dynamic Introspection:** A powerful, built-in **SSH console** allows you to inspect and modify the state of *any* component at runtime.
* **Snapshot-Driven:** A single, human-readable file defines the entire hardware configuration *and* its runtime state. Build, save, load, and share complex machines atomically.
* **Physically Accurate Simulation:** The framework prioritizes physical fidelity over high-level emulation, modeling complex behaviors like the C1541 drive mechanics or the exact bus arbitration of the C64 expansion port.
* **Truly Modular & Extensible:** Built entirely around interfaces, a component factory, and a central registry. Adding new hardware (like a Z80 CPU or a Yamaha sound chip) means simply creating a new, self-contained package—without modifying the core framework.
* **Pure Go Core:** Ensures type and memory safety, excellent performance, and effortless cross-compilation, with **no CGo dependencies**.

## Current Implementation: The Commodore 64

As its first and primary proof-of-concept, Symphony includes a fully **cycle-exact** Commodore 64 and an independent 1541 Floppy Drive implementation. This is not just a demo; it is a testament to the framework's power, capable of handling the most demanding software and hardware interactions with a level of fidelity rarely seen in other emulators.

The framework's component factory already registers a massive, production-ready library of topological building blocks, including:
* **Microprocessors & I/O:** MOS6510, MOS6522 (VIA), MOS6526 (CIA), External CPU
* **Custom Chips:** MOS6581 (SID), MOS6569 (VIC-II)
* **Logic & Motherboards:** C64 Board, C64 PLA, C1541 PLA, Quartz, Dynamic Throttle
* **Memory & Media:** C64 RAM, C64 Color RAM, C64 ROMs, C1541 RAM, C1541 ROMs, Generic Media
* **Peripherals & Buses:** C64 Keyboard, C64 Joystick, IEC Bus
* **RAM Expansions (REU):** 128K, 256K, 512K, 1MB, 2MB, 4MB, 8MB, and 16MB variants
* **Cartridge Formats:** Standard C64 Cartridges, Magic Desk, Ocean, EasyFlash, Final Cartridge III

## The Future: An Invitation to Build

The ultimate goal of Symphony is to create a platform where a global community of developers and hardware enthusiasts can contribute. The framework is the foundation. The next step is to expand the library of available "building blocks": more CPUs, more video and sound chips, and more systems.

This project is an invitation to anyone fascinated by the idea of not just using an emulator, but helping to build a universe of virtual machines.

## Getting Started

*(This section to be completed with detailed build and run instructions.)*

```bash
# Clone, build, and run the default C64 configuration
git clone [repository-url]
cd symphony
go build .
./symphony -r gl -m c64