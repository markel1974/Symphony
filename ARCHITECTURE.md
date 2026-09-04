# Symphony - Architectural Overview

[![License](https://img.shields.io/badge/License-LGPL_v2.1-blue.svg)](https://www.gnu.org/licenses/old-licenses/lgpl-2.1.en.html)

## 1. Introduction

Symphony is a highly configurable, deeply introspectable emulation framework written entirely in Go. Initially designed as a highly accurate Commodore 64 and 1541 emulator, it has evolved into a fully-fledged, system-agnostic framework for the dynamic exploration of computer systems.

It relies on a micro-kernel architecture, separating the core execution engine (Virtual Machine) from the hardware implementations (Component Tree) and using a message-based system for asynchronous operations. The framework is designed for deep runtime introspection and dynamic state management without CGo dependencies for its core logic.

## 2. High-Level Architecture

The Symphony ecosystem is built upon a modular separation of concerns. The primary domains are:

*   **Virtual Machine (VM):** An abstract execution platform capable of interpreting a Pluggable Instruction Set Architecture (ISA). It acts as the execution core for both hardware simulation and high-level scripting logic.
*   **Symphony Kernel:** Acts as the host system, managing processes, asynchronous inter-component messaging, process life-cycles, and terminal interfaces.
*   **Component Tree (Hardware):** A deeply hierarchical representation of simulated hardware. Everything is an `IComponent`. State and configuration are unified through dynamic tree loading and Snapshot restoration.
*   **Compilers:** A Go-native compiler infrastructure to translate a subset of the Go language into the VM's bytecode for dynamic scripting and advanced user-space processes.
*   **Renderers:** Independent modules responsible for the display output (e.g., OpenGL, text-mode).

## 3. Directory Structure and Modules (`src/`)

### 3.1. Kernel (`src/kernel`)

The kernel coordinates the execution environment:
*   **Processes (`process`):** Manages the VM instances as isolated userspace processes.
*   **Messages (`messages`):** Implements an asynchronous messaging API used for system calls, hardware interrupts, and inter-component signaling.
*   **Terminal (`terminal`):** Provides a built-in, ssh-accessible console to introspect the running state of the machine.

### 3.2. Virtual Machine (`src/vm`)

Symphony's execution engine is built around architectural symmetry:
*   **Bytecode & Opcodes (`bytecode`, `opcodes`):** Each instruction is a self-contained entity (`IOpExecutor`) containing definition, interpretation, and compilation logic.
*   **Sequencers (`sequencers`):** Pluggable execution cores (e.g. 6502/6510 instructions, high-level VM bytecode) that drive the state of the machine forward.

### 3.3. Compilers (`src/compilers`)

Contains tools and AST parsing to translate a significant subset of Go (interfaces, closures, structs) into the VM's native bytecode, enabling complex runtime scripting.

### 3.4. Hardware Components (`src/hardware`)

This directory houses the actual implementations of specific machines and architectures. All components conform to the `IComponent` interface and are instantiated via a central registry.

*   **`mos6510_rev1` (CPU):** Micro-operation based emulation of the 6510 CPU. Uses dispatch tables for cycle-accurate execution mapping.
*   **`mos6569_vic_rev1` (VIC-II):** Cycle-accurate emulation of the video chip, mapping raster timing, badlines, and sprites precisely.
*   **`mos6581_sid_rev1` (SID):** Sound chip emulation with waveforms, envelopes, and filters.
*   **`mos6526_cia_rev1` (CIA):** Handles I/O, timers, TOD, and keyboard/joystick scanning.
*   **`mos6522_via_rev1` (VIA):** Used in the 1541 disk drive.
*   **`c64_pla_rev1` (PLA):** Manages the complex memory mapping, bank switching, and cartridge access for the C64.
*   **Boards (`c64_board_rev1`, `c1541_board_rev1`):** The "motherboards" that orchestrate the connection between components (sockets mapping).
*   **Inputs & Storage:** Implementations for keyboards (`c64_keyboard_rev1`), joysticks (`c64_joystick_rev1`), disks and roms.

### 3.5. Registry (`src/registry`)

Provides a dynamic component factory. Components self-register via `init()` functions, allowing the framework to dynamically construct an entire machine from a JSON-like configuration snapshot.

### 3.6. Renderers (`src/renderers`)

Handles video output:
*   `gl`: OpenGL renderer for standard graphics.
*   `none`: Headless renderer for background tasks, server deployment, or testing.
*   `ascii`: Text-based fallback renderer.

## 4. Initialization and Data Flow

1.  **Snapshot Loading (`RestoreAll`):** The initialization starts by feeding a snapshot (`map[string]interface{}`) to the static factory. This recursively builds the Component Tree and restores the specific state/configuration of every component.
2.  **3-Phase Connection:**
    *   **Phase 1:** `RestoreAll` builds the components and their tree hierarchy.
    *   **Phase 2:** `Board.Setup` orchestrates `socket.Mount`, establishing structural dependencies dynamically across the tree.
    *   **Phase 3:** Final interface bindings (`Wire`) to establish communication lines between dependent modules.
3.  **Kernel Execution:** The kernel schedules VM processes. Processes run sequentially or concurrently depending on their sequencers.
4.  **Emulation Cycle:** For hardware like the C64, a clock-driven `Emulate()` cascade happens:
    *   VIC-II, CIA, SID, and CPU progress cycle-by-cycle based on micro-operations.
    *   Memory accesses go through the PLA mapping interface.
    *   Asynchronous kernel messages dispatch interrupts and hardware events.
5.  **Output:** Frame buffers are sent to the active renderer for display.

## 5. Key Design Decisions

*   **Snapshot = Configuration = State:** A single file determines everything.
*   **Decoupled Execution (Microkernel):** Emulated processors run as isolated processes making "syscalls" to the kernel.
*   **Pure Go (No CGo):** Ensures portability, simplicity of cross-compilation, and memory safety.
*   **Deep Introspection:** Built-in SSH server provides a window into the runtime state of any `IComponent`.