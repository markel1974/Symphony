# VIC-II Sprite Rendering: A Pipelined Approach

This document details the cycle-accurate, pipelined sprite rendering model of the MOS 6569 (VIC-II) chip and its specific implementation within this Go package. This architecture is a direct result of the "Structure Precedes Function" philosophy outlined in the main `README.md`, ensuring that complex hardware behaviors emerge naturally from a faithful software model.

## The Hardware Pipeline: A One-Line Delay

A common simplification in emulators is to assume that the data for a sprite on scanline `N` is fetched, prepared, and drawn entirely within the lifecycle of scanline `N`. On real hardware, this is not the case. The VIC-II is a highly efficient pipeline processor that works ahead.

The key principle is that the **preparation of sprite data for scanline `N+1` occurs during the execution of scanline `N`**. This creates a one-line delay that is crucial for correctly emulating timing-sensitive graphical effects.

The process is divided into two logical phases that overlap in time:

### 1. Preparation Phase (Scanline `N`)

While the raster beam is drawing the visible pixels of scanline `N`, the VIC-II is already using the horizontal and vertical blanking periods to prepare for the *next* line, `N+1`. This involves:

* **Y-Coordinate Check:** At the end of scanline `N`, the VIC-II checks which sprites' Y-coordinates match scanline `N+1`. This determines which sprites will be active and require data fetching.
* **DMA & Latching:** During the horizontal blanking period of scanline `N+1`, the VIC-II performs DMA for each of the 8 sprites. This is a distributed process where each sprite has its own dedicated time slot to fetch its data pointer and graphical patterns from RAM. Crucially, at the beginning of each sprite's DMA slot, its rendering attributes (X-position, color, priority) are "latched"—captured from the registers and frozen for the upcoming line.

### 2. Drawing Phase (Scanline `N+1`)

When the raster beam draws the visible portion of scanline `N+1`, the VIC-II is no longer reading from the main registers to render sprites. It uses the data that was **prepared and frozen during the previous cycle (scanline `N`)**. The video mixer combines the background graphics with the pixels from the prepared sprite data in real-time.

## The Software Implementation

This emulator faithfully replicates the hardware pipeline using a **double-buffered state** for sprites, orchestrated by the `Sequencer`.

### Double Buffering with `spriteToggle`

The core of the solution is a toggle (`seq.spriteToggle`) that flips every scanline. This creates two distinct sets of buffers for sprite data and attributes within each sprite object: an `odd` set and an `even` set.

This allows one buffer set to be safely **written to** (the preparation phase for the next line) while the other is being **read from** (the drawing phase for the current line).

* The functions `FetchPhase1` and `FetchPhase2` are called from the `Sequencer` at the precise, hardware-correct cycles for each sprite.
* They use the state of `spriteToggle` to write the fetched data and latched attributes into the "preparation" buffer set.

### Batched Drawing as a Serialization Step

The `Draw(odd bool)` function is called once at the end of the scanline's sequencer logic. While this appears to be a "batch" operation, it is simply the final serialization of the emulated state. The function uses the `odd` parameter to intelligently select and draw from the "display" buffer set—the one that was fully prepared during the previous scanline cycle.

This architecture correctly models the hardware's one-line delay, ensuring that even the most subtle timing-dependent graphical effects are rendered with high fidelity.