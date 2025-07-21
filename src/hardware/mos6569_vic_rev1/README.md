# Package mos6569

This package (`src/components/vic`) provides an emulation of the MOS Technology 6569 (PAL) / 8565 (PAL-B) Video Interface Chip II (VIC-II), the graphics controller chip used in the Commodore 64 home computer. This emulator aims for cycle-accurate emulation of the VIC-II's behavior.

## Overview

The `mos6569` package emulates all the core features of the VIC-II, including:

* **Raster Scanning:** Cycle-accurate emulation of the raster beam, including horizontal and vertical blanking intervals, border generation, and "badlines" (where the VIC-II steals memory cycles from the CPU).
* **Graphics Modes:**
    * Standard Character Mode (Text Mode) - with and without multicolor
    * Extended Background Color Mode (ECM) - text mode
    * Standard Bitmap Mode
    * Multicolor Bitmap Mode
* **Sprites:** Emulation of all eight hardware sprites.
* **Scrolling:** Smooth scrolling (both horizontal and vertical fine scrolling).
* **Interrupts:** Generation of raster interrupts.
* **Color Palette:** Support for the C64's 16-color palette.
* **Memory access**: Precise handling of video memory access.

## Implementation Philosophy: Structure Precedes Function

The effectiveness of this implementation stems from a deliberate design discipline: **structure precedes function**. Instead of writing a program that simply simulates the final *output* of the VIC-II, this is a software model that replicates **how the chip is built**, based directly on its original block diagrams.

This approach is a form of **Register-Transfer Level (RTL) Modeling**, a term from electronic engineering. The correct visual output is not the direct goal, but an **emergent consequence** of a faithfully replicated hardware structure being set in motion.

The pillars of this philosophy are:

### 1. Architecture as a Mirror of the Hardware

Each software object in the code is a 1:1 translation of a functional block from the MOS 6569 chip:
* **`VIC`**: The entire integrated circuit, connecting all other blocks.
* **`Sequencer`**: The **clock signal** and control logic that orchestrates every cycle.
* **`Memory`**: The **memory bus interface**, handling VIC's access to RAM and ROM.
* **`Graphics`**, **`SpriteHandler`**, **`Borders`**: The discrete logic blocks for graphics, sprites, and borders, each with its own precise structural responsibility.

### 2. Function as an Emergent Behavior

Only after defining this rigid architecture is the function implemented:
* The main emulation loop (`Emulate`) simply sends a **clock pulse** to the `Sequencer`.
* The phase functions (`phase...`) describe what each hardware block does during that single clock cycle.

This ensures that the output is not only correct, but correct for the right reasons, replicating the timing and side-effects of the real hardware with a fidelity that would otherwise be unattainable.

### 3. Optimization Through Design (Near Zero-Branch)

Performance is not a result of "hacks" but a product of the design itself. The use of **dispatch tables** and **function pointers** instead of `switch`/`if` constructs in critical paths is a natural consequence of this RTL model, leading to extremely fast and predictable execution.

## Package Structure

The `vic` package is organized into the following files:

* `vic.go`: The main file. Contains the `VIC` struct and the `Emulate` method. Also includes functions for reading and writing VIC registers.
* `sequencer.go`: Contains constants and functions related to the video timing, one for each cycle of a scanline.
* `graphics.go`: Contains the functions for rendering the different graphics modes.
* `sprite_handler.go`: Contains the high-level logic for handling sprites.
* `sprite.go`: Contains the logic for handling a single sprite.
* `borders.go`: Contains the logic for rendering the borders.
* `collisions.go`: Implements the collision detection logic.
* `memory.go`: Implements memory access.
* `tables.go`: Contains pre-calculated lookup tables used for optimization.
* `factory.go`: Implements the factory pattern for integration into the system.

## Dependencies

* `github.com/markel1974/c64emu/src/component`
* `github.com/markel1974/c64emu/src/references`
* `github.com/markel1974/c64emu/src/config`
* `github.com/markel1974/c64emu/src/registry`

## Usage

The `mos6569` package is tightly integrated with the other components of the emulator (CPU, memory, etc.).

## TODO

* Investigate and implement any missing undocumented VIC-II behaviors.

## License

This project is released under the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).