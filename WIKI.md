# g64: An Accurate and Innovative Commodore 64 Emulator

Welcome to the official wiki for `g64`, a high-fidelity, cycle-accurate emulator for the Commodore 64 (C64) and VIC-20 computers. This wiki serves as a comprehensive guide to the emulator's architecture, components, features, and usage.

## Table of Contents

1.  [Introduction](#introduction)
2.  [Key Features](#key-features)
3.  [Architecture Overview](#architecture-overview)
    *   [Cycle-Accurate Emulation](#cycle-accurate-emulation)
    *   [Signal-Based Decoupling](#signal-based-decoupling)
    *   [Socket-Based Interfacing](#socket-based-interfacing)
    *   [Clock Management](#clock-management)
    *	[Memory Management](#memory-management)
    * [Component Management](#component-management)
    * [Peripherals Management](#peripherals-management)
    * [Interrupt Management](#interrupt-management)
    * [Input Management](#input-management)
    * [Cartridge Management](#cartridge-management)
    * [Files Management](#files-management)
    * [Configuration Management](#configuration-management)
    * [Error Management](#error-management)
    * [Three Renderers](#three-renderers)
4.  [Components](#components)
    *   [CPU (MOS 6510)](#cpu-mos-6510)
    *   [VIC-II (MOS 6569)](#vic-ii-mos-6569)
    *   [SID (MOS 6581)](#sid-mos-6581)
    *   [CIA (MOS 6526)](#cia-mos-6526)
    *   [PLA](#pla)
    *   [Quartz](#quartz)
    *   [PIC (Programmable Interrupt Controller)](#pic-programmable-interrupt-controller)
    *   [IEC (Serial Bus)](#iec-serial-bus)
    * [1541 (Disk Drive)](#1541-disk-drive)
    * [REU (RAM Expansion Unit)](#reu-ram-expansion-unit)
5. [Cartridges](#cartridges)
    * [Generic Cartridges](#generic-cartridges)
	* [EasyFlash](#easyflash)
	* [Ocean](#ocean)
	* [MagicDesk](#magicdesk)
	* [REU](#reu)
6.  [Renderers](#renderers)
    *   [OpenGL Renderer (`glrender`)](#opengl-renderer-glrender)
    *   [WebAssembly Renderer (`wasmrender`)](#webassembly-renderer-wasmrender)
    *   [ASCII Renderer (`asciirender`)](#ascii-renderer-asciirender)
7.  [Inputs](#inputs)
    *   [Keyboard](#keyboard)
    *   [Joystick](#joystick)
    *   [Mouse](#mouse)
8. [Usage](#usage)
	* [Command Line Options](#command-line-options)
9.  [Building from Source](#building-from-source)
10. [Contributing](#contributing)
11. [License](#license)
12. [Acknowledgments](#acknowledgments)

## Introduction

`g64` is a meticulously crafted emulator designed to faithfully replicate the experience of using a Commodore 64 and VIC-20 computer. Built entirely in Go, `g64` prioritizes cycle-accurate emulation, modularity, and a high degree of accuracy.

## Key Features

*   **Cycle-Accurate Emulation:** Provides a highly accurate emulation by synchronizing all components to a central clock.
*   **Modular Architecture:** Uses interfaces, sockets, and a unique signal system for component decoupling and flexibility.
*   **No External Dependencies:** Implemented from scratch in Go, ensuring complete control and minimal bloat.
*   **WebAssembly Support:** Offers a proof-of-concept WebAssembly renderer for browser-based emulation.
*   **Multiple Renderers:** Supports OpenGL, WebAssembly, and ASCII renderers.
*   **Comprehensive Hardware Emulation:** Emulates many components of the original hardware, including:
    *   MOS 6510 CPU
    *   MOS 6569 VIC-II (Video Interface Chip)
    *   MOS 6581 SID (Sound Interface Device)
    *   MOS 6526 CIA (Complex Interface Adapter)
    *   PLA (Programmable Logic Array)
    *   1541 Disk Drive (Mechanics and GCR)
    * IEC bus
    * Cartridges (EasyFlash and custom)
    * REU (Ram Expansion Unit)
*   **Cartridge and Disk Image Support:** Runs software from various formats (D64, CRT, BIN).
*   **Efficient:** The emulator is efficient and fast.
* **Correct management:** the project provides a correct and efficient management of:
	* Memory
	* Types
	* Components
	* Communication
	* Configuration
	* Errors
	* Inputs
	* Interrupts
* **Clipboard:** the project supports the clipboard.
* **Joystick and keyboard:** the project manage correctly the joystick and the keyboard.
* **Complete:** the project is very complete and allow to run most of the original software.

## Architecture Overview

`g64` features an innovative architecture that emphasizes accuracy, modularity, and maintainability.

### Cycle-Accurate Emulation

The core of `g64` is driven by a central clock (`src/components/quartz`) that synchronizes all components. Each component advances based on clock cycles, enabling precise cycle-accurate emulation without explicit timing management.

### Signal-Based Decoupling

`g64` uses a custom signal system (`src/signals`) for communication between components. This method enables significant decoupling, improving modularity and simplifying future extensions.

### Socket-Based Interfacing

Components in `g64` interact through socket interfaces (e.g., `CPUSocket`, `VicSocket`, `SidSocket`). These sockets define clear communication paths and encapsulate each component's internal logic. This method enhances modularity and maintainability.
### Clock Management
The clock management is handled by the `src/components/quartz` package.
### Memory Management
The project provide a correct memory management.
### Component Management
The project provide a correct management of all components.
### Peripherals Management
The project provide a correct management of all peripherals.
### Interrupt Management
The project provide a correct management of all interrupt.
### Input Management
The project provide a correct management of all input.
### Cartridge Management
The project provide a correct management of all cartrides.
### Files Management
The project provide a correct management of the files.
### Configuration Management
The project provide a correct management of the configuration.
### Error Management
The project has the capability to manage errors.
### Three Renderers
The project provide the ability to use three different types of renderers.

## Components

### CPU (MOS 6510)

*   **Responsibility:** The central processing unit, responsible for executing instructions.
*   **Location:** `src/components/6510`
*   **Socket:** `CPUSocket` (defined in `src/c64/board/cpusocket.go` and `src/c1541/board/cpusocket.go`)

### VIC-II (MOS 6569)

*   **Responsibility:** The Video Interface Chip, responsible for generating video output.
*   **Location:** `src/components/vic`
*   **Socket:** `VicSocket` (defined in `src/c64/board/vicsocket.go`)

### SID (MOS 6581)

*   **Responsibility:** The Sound Interface Device, responsible for generating audio output.
*   **Location:** `src/components/sid`
*   **Socket:** `SidSocket` (defined in `src/c64/board/sidsocket.go`)

### CIA (MOS 6526)

*   **Responsibility:** The Complex Interface Adapter, responsible for handling I/O operations, timers, and interrupts.
*   **Location:** `src/components/cia`
* **Sockets:**
	* `CIA1Socket` (defined in `src/c64/board/cia1socket.go`)
	* `CIA2Socket` (defined in `src/c64/board/cia2socket.go`)

### PLA

* **Responsibility:** The Programmable Logic Array, responsible for handling memory mapping and chip select logic.
* **Location:** `src/c64/pla`

### Quartz

*   **Responsibility:** The central clock, responsible for synchronizing all components.
*   **Location:** `src/components/quartz`

### PIC (Programmable Interrupt Controller)

*   **Responsibility:** The programmable interrupt controller, responsible for managing interrupts.
*   **Location:** `src/components/6510`

### IEC (Serial Bus)

*   **Responsibility:** The IEC (Commodore Serial Bus) is responsible for communication with peripheral devices.
*   **Location:** `src/components/iec/virtualdrive`

### 1541 (Disk Drive)

* **Responsibility:** The disk drive, responsible for disk operations.
* **Location:** `src/c1541`
* **Sockets:**
	* `Via1Socket` (defined in `src/c1541/board/via1socket.go`)
	* `Via2Socket` (defined in `src/c1541/board/via2socket.go`)

### REU (RAM Expansion Unit)

* **Responsibility:** The RAM Expansion Unit (REU) is an optional memory expansion for the C64.
* **Location:** `src/c64/cartridges/reu`

## Cartridges

`g64` supports a variety of cartridge types.

### Generic Cartridges

* **Description:** Generic cartridges are ROM cartridges without any special bank switching or other features.
* **Location:** `src/c64/cartridges/generic`

### EasyFlash

* **Description:** EasyFlash is a popular cartridge that allows you to run multiple cartridges from a single cartridge.
* **Location:** `src/c64/cartridges/easyflash`

### Ocean

* **Description:** Ocean cartridges are cartridges with special bank switching.
* **Location:** `src/c64/cartridges/ocean`

### MagicDesk

* **Description:** Magic Desk cartridges are cartridges with special bank switching.
* **Location:** `src/c64/cartridges/magicdesk`

### REU

* **Description:** REU (RAM Expansion Unit) is a memory expansion for the C64.
* **Location:** `src/c64/cartridges/reu`

## Renderers

`g64` supports three different rendering methods.

### OpenGL Renderer (`glrender`)

*   **Responsibility:** Provides high-performance graphical output using OpenGL.
*   **Location:** `src/render/glrender`

### WebAssembly Renderer (`wasmrender`)

*   **Responsibility:** Enables running the emulator in a web browser using WebAssembly.
*   **Location:** `src/render/wasmrender`

### ASCII Renderer (`asciirender`)

*   **Responsibility:** Outputs a text-based representation of the C64 screen in the console.
*   **Location:** `src/render/asciirender`

## Inputs

`g64` supports keyboard, joystick, and mouse inputs.

### Keyboard

*   **Responsibility:** Handles keyboard input.
* **Location:** `src/c64/inputs` and `src/vic20/inputs`

### Joystick

*   **Responsibility:** Handles joystick input.
* **Location:** `src/c64/inputs/joystick.go` and `src/vic20/inputs/joystick.go`

### Mouse

*   **Responsibility:** Handles mouse input.
* **Location:** `src/c64/board/board.go` and `src/vic20/board/board.go`

## Usage

### Command Line Options

*   `-c`: Path to a cartridge file (CRT, BIN). You can use `Key:Value` to set more cartridges.
    ```bash
        -c "/path/to/my/cartridge.crt"
        -c "SCPU:;REU16M:/path/to/my/doom.reu"
    ```
*   `-d`: Path to a drive file (D64). You can use `Key:Value` to set more drives.
    ```bash
        -d "/path/to/my/disk.d64"
        -d "8:/path/to/my/disk.d64;9:/path/to/my/disk2.d64"
    ```
*   `-f`: Path to one or more disk images. You can use `;` to set more disk files. The files will be inserted into the selected drive.
     ```bash
        -f "/path/to/my/disk.d64;/path/to/my/disk2.d64"
        -f "/path/to/my/disk.d64"
    ```
*   `-m`: Set the hardware mode (vic20 or c64).
    ```bash
        -m "vic20"
    ```
*   `-p`: Path to a prg file
    ```bash
        -p "/path/to/my/program.prg"
    ```
*   `-j`: Disable JiffyDOS speed-up.
*   `-a`: Use ASCII render.
*   `-h`: Show help.
*   `-v`: Show version.

## Building from Source

1.  Clone the repository:

    ```bash
    git clone https://github.com/markel1974/c64emu.git
    cd c64emu
    ```

2.  Build the project:

    ```bash
    go build -o g64 main.go
    ```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests on GitHub.
Please read the CONTRIBUTING.md for more details.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

*   The Commodore 64 community for their continued enthusiasm and support.
*   The developers of VICE and other emulators, whose work has been an inspiration.