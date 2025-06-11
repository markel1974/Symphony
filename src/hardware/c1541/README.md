# Commodore 1541 Drive Emulation Package

## Overview

This package provides a complete, high-fidelity, and physically-modeled simulation of the Commodore 1541 disk drive, designed for integration within the Symphony emulation framework.

The primary goal of this package is not just to achieve logical compatibility (i.e., reading and writing `.d64` sectors), but to emulate the **physical and temporal behavior** of the real hardware at a cycle-accurate level. This ensures maximum compatibility with software that relies on precise timing, undocumented drive behavior, and advanced copy-protection schemes.

## Architectural Philosophy

Following the core philosophy of the Symphony framework, the C1541 drive is not a monolithic block. It is modeled as a system of distinct, modular components with a strict separation of concerns, mirroring the subsystems of the actual hardware.

The main components are:

1.  **The Board (`/board`):** The top-level component representing the drive's main logic board. It orchestrates all other internal components.
2.  **The Mechanic (`/mechanic`):** The physics simulation engine. It models the behavior and timing of the drive's electromechanical parts (the motor and the read/write head).
3.  **The Disk (`/disk`):** The representation of the magnetic disk media itself. It handles data at the physical track level, including GCR encoding and media effects.

## Component Breakdown

### 1. The Board (`board`)

The `Board` component acts as the central hub for a complete drive instance.
* **Orchestration:** It contains and connects the drive's internal `6502 CPU`, the two `6522 VIA` chips, the `PLA` for memory mapping, the `ROMs` (CBM DOS), and the `Mechanic`.
* **Emulation Cycle:** Its `Emulate()` method propagates the system clock "tick" to all its children, ensuring the entire drive system advances in lock-step.
* **IEC Bus Interface:** It serves as the endpoint for commands coming from the main computer's IEC bus.

### 2. The Mechanic (`mechanic`)

This is the heart of the physical simulation. It is responsible for modeling the complex, imperfect, and analog nature of the drive's hardware. The high-fidelity `Async` implementation includes:

* **Cycle-Accurate Timing:** Every operation is measured in 1MHz clock cycles, matching the real hardware's timing source.
* **Detailed Head Seek Model:** The movement of the read/write head is not instantaneous. The model simulates:
* **Asymmetric Seek Delays:** Different timings for moving inward (towards higher tracks) versus outward.
* **Cumulative Damping:** Consecutive steps in the same direction increase the settling time required due to accumulated kinetic energy.
* **Directional Backlash:** A one-time penalty is applied on the first step after a direction change to simulate mechanical slack and motor polarity reversal.
* **Vibration Factor:** Frequent, rapid direction changes ("Z-pattern" seeks) increase a "vibration factor" that temporarily lengthens all settling times, modeling mechanical stress.
* **Natural Damping:** The vibration factor slowly decays over time when the head is idle, allowing the virtual mechanism to "settle down."
* **Realistic Motor:** Simulates the `spin-up` delay required for the motor to reach its stable 300 RPM speed.

### 3. The Disk (`disk`)

This component represents the floppy disk media.

* **GCR (Group Code Recording):** Manages the conversion of standard bytes to and from the GCR format used to physically store data on the disk surface.
* **Physical Track Layout:** Accurately models the track structure based on hardware constants, including the precise size of inter-sector gaps and the final tail gap for each of the four speed zones.
* **Half-Track and Noise Simulation:** To support advanced copy protections, the disk model implements half-tracks. When the head is positioned on an "odd" half-track (e.g., 18.5), the `ApplyNoise()` method is used to introduce realistic, deterministic bit-level corruption to the track data, simulating the signal interference from adjacent physical tracks.

## How It Works: A Read Operation Flow

1.  A command from the emulated C64 arrives via the IEC bus to the C1541 `Board`.
2.  The `Board` routes the command to the drive's internal `CPU`.
3.  The drive's `CPU` executes its ROM code (CBM DOS) to handle the request.
4.  The ROM code issues commands to the `Mechanic` (e.g., "move head to track 18").
5.  The `Mechanic`'s `Head.Move()` state machine simulates the physical movement, step by step, calculating the precise cycle delay for each step based on its complex timing model.
6.  Once the head is stable, the ROM commands a read.
7.  The `Mechanic`'s `ReadWrite()` method requests bytes from the `Disk` object.
8.  The `Disk` object, using the master clock to simulate rotation, provides a continuous stream of GCR bytes from the correct track (applying noise if on a half-track).
9.  The data is processed by the drive's ROM and sent back to the C64.