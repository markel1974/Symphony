# Package mos6569

This package provides an emulation of the MOS Technology 6569/8565/8566(PAL)—6567/8562/8564(NTSC) Video Interface Chip II (VIC-II), the graphics controller chip used in the Commodore 64 home computer. This emulator aims for cycle-accurate emulation of the VIC-II's behavior.

## Overview

The `mos6569` package emulates all the core features of the VIC-II, including:

* **Raster Scanning:** Cycle-accurate emulation of the raster beam, including horizontal and vertical blanking intervals, border generation, and "bad lines" (where the VIC-II steals memory cycles from the CPU).
* **Graphics Modes:**
  * Standard Character Mode (Text Mode) - with and without multicolor
  * Extended Background Color Mode (ECM)
  * Standard Bitmap Mode
  * Multicolor Bitmap Mode
* **Sprites:** Emulation of all eight hardware sprites, including a faithful replication of the one-line rendering pipeline delay.
* **Scrolling:** Smooth scrolling (both horizontal and vertical fine scrolling).
* **Interrupts:** Generation of raster interrupts.
* **Color Palette:** Support for the C64's 16-color palette.
* **Memory Access**: Precise handling of video memory access.

## Implementation Philosophy: Structure Precedes Function

The effectiveness of this implementation stems from a deliberate design discipline: **structure precedes function**. Instead of writing a program that simply simulates the final *output* of the VIC-II, this is a software model that replicates **how the chip is built**, based directly on its original block diagrams.

This approach is a form of **Register-Transfer Level (RTL) Modeling**, a term from electronic engineering. The correct visual output is not the direct goal, but an **emergent consequence** of a faithfully replicated hardware structure being set in motion.

The pillars of this philosophy are:

### 1. Architecture as a Mirror of the Hardware

Each software object in the code is a 1:1 translation of a functional block from the MOS 6569 chip:
* **`VIC`**: The entire integrated circuit, connecting all other blocks.
* **`Sequencer`**: The **clock signal** and control logic that orchestrates every cycle.
* **`MemoryUnit`**: The **memory bus interface**, handling VIC's access to RAM and ROM.
* **`GraphicsUnit`**, **`SpritesUnit`**, **`BordersUnit`**: The discrete logic blocks for graphics, sprites, and borders, each with its own precise structural responsibility.

**Perfect Architectural Isolation**: Each component maintains complete encapsulation with no direct access to the internal fields of other components. Communication occurs exclusively through well-defined public methods and interfaces, mirroring the physical pin-based connections of the original hardware blocks.

### 2. Architectural Optimization (Compiler-Friendly Design)

Performance in this emulator is not achieved through post-design "hacks," but is an intrinsic property of its architecture, which is deliberately designed to be compiler-friendly.

* **Pervasive Inlining:** The core principle is to enable aggressive **function inlining** by the Go compiler across the entire codebase. Nearly every method within critical components—from the `Sequencer`'s phase functions to memory accessors in `MemoryUnit`, rendering logic in `GraphicsUnit`, and state updates in `SpritesUnit`—is intentionally kept small, direct, and often marked with the `//go:nosplit` directive. This design minimizes the overhead of function calls, which would otherwise accumulate significantly in the tight, cycle-by-cycle emulation loop.

* **Near Zero-Branch Execution:** The use of **dispatch tables** and **function pointers** instead of `switch` or `if` statements is another key aspect of this philosophy. This allows the emulator to transition between states and modes with minimal branching, resulting in highly predictable and fast execution paths.

The result is an architecture where the high-level, readable Go code is "unrolled" by the compiler into a linear and highly optimized machine code representation, achieving cycle-accurate emulation with maximum performance.

### 3. Function as an Emergent Behavior

Only after defining this rigid architecture is the function implemented. The main emulation loop (`Emulate`) simply sends a **clock pulse** to the `Sequencer`, which propagates through the interconnected components. This ensures that the output is not only correct, but correct for the right reasons, replicating the timing and side-effects of the real hardware with a fidelity that would otherwise be unattainable.

## Package Structure

The `vic` package is organized into the following files:

* `vic.go`: The main file. Contains the `VIC` struct and the `Emulate` method. Also includes functions for reading and writing VIC registers.
* `sequencer.go`: Contains constants and functions related to the video timing, one for each cycle of a scanline.
* `graphics_unit.go`: Contains the functions for rendering the different graphics modes.
* `sprites_unit.go`: Contains the high-level logic for handling sprites.
* `sprite.go`: Contains the logic for handling a single sprite.
* `borders_unit.go`: Contains the logic for rendering the borders.
* `collisions_unit.go`: Implements the collision detection logic.
* `memory_unit.go`: Implements memory access.
* `interrupts.go`: Manages all VIC-II interrupt sources (Raster, Sprite-Sprite, Sprite-Background, Light Pen).
* `light_pen.go`: Handles the light pen logic, capturing coordinates and triggering interrupts.
* `tables.go`: Contains pre-calculated lookup tables used for optimization.
* `factory.go`: Implements the factory pattern for integration into the system.

## Dependencies

* `github.com/markel1974/c64emu/src/component`
* `github.com/markel1974/c64emu/src/references`
* `github.com/markel1974/c64emu/src/config`
* `github.com/markel1974/c64emu/src/registry`

## License

This project is released under the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).