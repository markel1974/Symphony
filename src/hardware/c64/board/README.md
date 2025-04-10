# Package `board` in `src/c64/`

The `board` package in `src/c64/` is the core component responsible for orchestrating the emulation of the Commodore 64's main hardware. It acts as the central hub connecting and managing all critical components. It provide the management of the various peripherals, the management of the memory and the management of the bus.

## Responsibilities

The primary responsibilities of the `board` package include:

*   **Component Management:** Creating, initializing, and managing all key C64 hardware components:
    *   CPU (MOS 6510)
    *   VIC-II (Video Interface Chip)
    *   SID (Sound Interface Device)
    *   CIA (Complex Interface Adapter)
    *   PLA (Programmable Logic Array)
    *   Interrupt Controller
    * Memory
    * Cartridges
    * IEC
* **Disk Drive Management:** manage the disk drive.
*   **Emulation Loop:** Implementing the main emulation loop (`Emulate()`) that drives the simulation of all components.
*   **Interrupt Handling:** Managing interrupt requests from various components.
*   **Memory Management:** managing the memory.
* **Bus management:** managing the bus.
*   **Input Handling:** Processing keyboard, joystick, and mouse inputs.
*   **Reset:** Implementing the reset of the C64.
* **Throttling:** manage the throttling.
*   **Configuration:** Manage the configuration and the options.
* **Clipboard:** manage the clipboard.
* **Clock management:** manage the clock.
* **Signals management:** manage the various signals.
* **Components cycle of life:** manage the cycle of life of the various components.

## Key Files and Components

### `board.go`

*   **`Board` Struct:** Represents the central hardware board of the C64, containing instances of all the key components (CPU, VIC, SID, CIA, PLA, etc.).
*   **`NewBoard()`:** Initializes a new `Board` instance.
*   **`Setup()`:** Initializes all C64 components, including the CPU, VIC, SID, CIA, PLA, and the interrupt controller. It also handles the configuration and the setup of the system.
*   **`reset()`:** Resets all components to their initial states.
* **`AsyncReset()`:** Esegue un reset asincrono.
*   **`Emulate()`:** Executes a single emulation cycle, updating all the components and the clock.
* **`Throttle()`:** Return the throttling object.
* **`GetText()`:** return the text from the VIC.
* **`GetScreenSize()`:** return the size of the screen.
* **`DiskChange()`:** handles the disk change.
* **`KeyboardPaste()`, `KeyboardSetCommand()`, `KeyboardNumLockToggle()`, `KeyboardCapitalToggle()`, `KeyboardSetVirtualKey()`:** handle the keyboard inputs.
* **`SetMouse()`:** handle the mouse input.
* **`Joy1SetKey()`, `Joy2SetKey()`:** handle the joystick inputs.
* **`JoySwap()`:** manage the joystick swap.
* **`ExtRamRead()`, `ExtRamWrite()`:** manage the reu ram read/write.
* **`Joystick1Move`, `Joystick2Move`:** manage the joystick movement.
* **`dmaLowSlot`, `rdyLowSlot`, `aecLowSlot`:** manage the bus signals.
* **`irqTriggerSlot`, `irqClearSlot`:** manage the irq signals.
* **`nmiTriggerSlot`, `nmiClearSlot`:** manage the nmi signals.
* **`vicLastCycleSLot`:** manage the vic last cycle signal.
* **`vicVBlankSlot`:** manage the vic vblank signal.
* **`ledStateChangedSlot`:** manage the led state changed signal.

### Socket Interfaces

The `board` package uses socket interfaces to connect with individual components, providing a clean and modular design. Each socket defines a contract for communication with its respective component.

*   **`CPUSocket` (`cpusocket.go`):**
    *   **Responsibility:** Manages communication with the CPU.
    * **Methods:** `NewCPUSocket()`, `Setup()`, `Reset()`, `GetPic()`, `GetBanks()`.
*   **`VicSocket` (`vicsocket.go`):**
    *   **Responsibility:** Manages communication with the VIC-II.
    *   **Methods:** `NewVicSocket()`, `Setup()`, `Reset()`, `Cycle()`, `GetDisplayBuffer()`, `GetBanks()`, `IRQTrigger()`, `SetBanks()`.
*   **`SidSocket` (`sidsocket.go`):**
    *   **Responsibility:** Manages communication with the SID.
    *   **Methods:** `NewSidSocket()`, `Setup()`, `Reset()`, `SetPotXY()`, `Prepare()`, `Update()`, `GetPlayer()`.
*   **`CIA1Socket` (`cia1socket.go`):**
    *   **Responsibility:** Manages communication with CIA1.
    *   **Methods:** `NewCIA1Socket()`, `Setup()`, `Reset()`, `Update()`, `ReadPortA()`, `ReadPortB()`, `WritePortA()`, `WritePortB()`, `WriteDdrA()`, `WriteDdrB()`, `IRQTrigger()`.
*   **`CIA2Socket` (`cia2socket.go`):**
    *   **Responsibility:** Manages communication with CIA2.
    *   **Methods:** `NewCIA2Socket()`, `Setup()`, `Reset()`, `Update()`, `ReadPortA()`, `ReadPortB()`, `WritePortA()`, `WritePortB()`, `WriteDdrA()`, `WriteDdrB()`, `IRQTrigger()`.
*   **`Expansion` (`expansion.go`):**
    * **Responsibility:** manages the expansion.
    * **Methods:** `NewExpansion()`, `Reset()`, `Read()`, `Write()`, `GameExRomConfigChanged()`, `NMITrigger()`, `SetDMALow()`, `ResetTrigger()`, `IRQTrigger()`, `IRQClear()`, `IRQTriggerBind()`, `IRQClearBind()`, `BusAvailable()`, `AECAvailable()`, `Cycle()`, `CycleAlarm()`, `RamSetWriteTrigger()`, `RamRemoveWriteTrigger()`, `RmwFlags()`.

### PLA Package (`src/c64/pla`)

*   **`MemoryMap` (`memorymap.go`):**
    * **Responsibility:** Manages the memory map configurations of the C64, providing a way to switch between different memory layouts.
    * **Methods:** `NewMemoryMap()`, `Get()`.
* **`EmulatorId` (`emulatorid.go`):**
    * **Responsibility:** Manages the emulator id.
    * **Methods:** `NewEmulatorId()`, `Read()`.

## Usage

The `board` package is used internally by the main emulation loop to drive the simulation. The renderers and input components interact with the `Board` through the interface `components/board/iboard.go`.

## Key Concepts

*   **Cycle Accuracy:** The package is designed to support cycle-accurate emulation by synchronizing component operations to a central clock.
*   **Modularity:** The use of sockets enhances modularity, making it easier to add or modify components.
* **Clock:** the package use a clock management system.
* **Interrupt:** the package use an interrupt management system.
* **Memory:** the package use a memory management system.
* **Components:** the package manage all the components.
* **Peripherals:** the package manage all the peripherals.
* **Bus:** the package manage the bus.
* **Cartridges:** the package manage the cartridges.
* **Disk drive:** the package manage the disk drive.

## Further Development

*   **Error Handling:** Enhance error handling in all functions to provide better feedback and stability.
* **More comments:** add more comments to the code.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests on GitHub.

## License

This project is licensed under the [MIT License](LICENSE).