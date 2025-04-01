# Package mos6522

This package (`src/hardware/mos6522`) provides an emulation of the **MOS Technology 6522 Versatile Interface Adapter (VIA)** chip. The VIA is a common peripheral IC used in many 8-bit systems, including the Commodore 1541 disk drive (where two VIAs are present, typically mapped at `$1800` or `$1C00`). It provides parallel I/O ports, timers, and a shift register.

## Overview

The `mos6522` package implements the core functionalities of the VIA, including:

* **Parallel I/O Ports:** Emulation of Port A (PRA) and Port B (PRB), including their corresponding Data Direction Registers (DDRA, DDRB). Interaction with external hardware connected to these ports is **delegated** to an `IVIASocket` interface provided during initialization.
* **Timers:**
  * **Timer 1 (T1):** 16-bit timer operating in one-shot or continuous (free-running) mode, capable of generating interrupts on underflow and potentially controlling output on PB7.
  * **Timer 2 (T2):** 16-bit timer operating in one-shot mode (timed interval) or pulse counting mode (counting pulses on PB6). Capable of generating interrupts on underflow.
* **Shift Register (SR):** Basic register implementation (full serial shifting logic might be simplified or depend on external clocking not shown in `Emulate`).
* **Control Registers:** Emulation of the Auxiliary Control Register (ACR) and Peripheral Control Register (PCR) which configure timer modes, shift register operation, and handshake lines (CA1, CA2, CB1, CB2 - *Note: Handshake line logic might be simplified or delegated*).
* **Interrupts:** Emulation of the Interrupt Flag Register (IFR) and Interrupt Enable Register (IER). Interrupt generation is triggered based on timer underflows (T1, T2) and potentially other sources (SR, handshake lines - *needs verification*), and the final IRQ signal is managed via the `IVIASocket`.

## `VIA` Struct

The core emulation logic is encapsulated in the `VIA` struct:

```go
type VIA struct {
    *component.BaseComponent // Embeds base component features
    pra    uint8             // Port A Data Register (ORA/IRA)
    ddra   uint8             // Port A Data Direction Register (DDRA)
    prb    uint8             // Port B Data Register (ORB/IRB)
    ddrb   uint8             // Port B Data Direction Register (DDRB)
    t1c    uint16            // Timer 1 Counter (read)
    t1l    uint16            // Timer 1 Latch (write)
    t2c    uint16            // Timer 2 Counter (read)
    t2l    uint16            // Timer 2 Latch (write)
    sr     uint8             // Shift Register
    acr    uint8             // Auxiliary Control Register
    pcr    uint8             // Peripheral Control Register
    ifr    uint8             // Interrupt Flag Register (pending IRQs)
    ier    uint8             // Interrupt Enable Register (IRQ mask)
    socket references.IVIASocket // Interface to the specific hardware connection/board
}
```
## Key Methods

* **`NewVIA(...) *VIA`:** Constructor. Initializes the VIA struct with default register values (mostly 0) and registers it within the Symphony component tree using `BaseComponent.Register`.
* **`Setup() error`:** Part of the `IHardware` interface. Handles component-specific setup after all components exist but before connections.
* **`Bind(socket references.IVIASocket) error`:** Called during the connection phase (likely by the socket's `Mount` method). Receives and stores the `IVIASocket` interface, providing the link for the VIA to interact with the external world (read/write ports, signal IRQs).
* **`Connect() error`:** Part of the `IHardware` interface. Handles final connection steps after `Bind`.
* **`Reset()`:** Resets all internal VIA registers to 0.
* **`ReadByte(addr uint16) uint8`:** Handles reads from the VIA's 16 memory-mapped registers (offset `addr & 0x0f`).
  * **Ports A/B (`$x0`, `$x1`, `$xF`):** Delegates reading the *effective* port state (considering input pins and DDR) to `socket.ReadPRA/B()`.
  * **DDRs (`$x2`, `$x3`):** Returns the current value of `ddra` or `ddrb`.
  * **Timers (`$x4`-`$x9`):** Returns the low/high bytes of timer counters or latches. Reading T1 counter low (`$x4`) or T2 counter (`$x8`) also clears their respective interrupt flags in the IFR.
  * **SR (`$xA`):** Returns the Shift Register value.
  * **ACR/PCR (`$xB`, `$xC`):** Returns the control register values.
  * **IFR (`$xD`):** Returns the Interrupt Flag Register. If any enabled interrupt is pending, sets the top bit (`0x80`). Reading IFR *clears* all pending flags in `ifr` and signals `socket.IRQClear()`.
  * **IER (`$xE`):** Returns the Interrupt Enable Register (with bit 7 always set).
* **`WriteByte(addr uint16, data uint8)`:** Handles writes to the VIA's registers.
  * **Ports A/B (`$x0`, `$x1`, `$xF`):** Updates the internal data register (`pra`/`prb`) and calls `socket.WritePRA/B()` to allow the socket to handle the actual output based on DDR.
  * **DDRs (`$x2`, `$x3`):** Updates `ddra`/`ddrb` and calls `socket.WriteDDRA/B()`.
  * **Timers (`$x4`-`$x9`):** Writes to the low/high bytes of the timer latches (`t1l`, `t2l`). Writing to the high byte of T1 (`$x5`) also transfers the latch to the counter (`t1c`) and clears the T1 interrupt flag. Writing to T2 high byte (`$x9`) does similarly for T2.
  * **SR (`$xA`):** Writes to the Shift Register.
  * **ACR/PCR (`$xB`, `$xC`):** Writes to the control registers.
  * **IFR (`$xD`):** Writing clears bits in the IFR corresponding to the bits set in `data`.
  * **IER (`$xE`):** Writing with bit 7 set *enables* interrupts specified by bits 0-6. Writing with bit 7 clear *disables* interrupts specified by bits 0-6. Calls `irqTrigger()` after update.
* **`Emulate()`:** Simulates one clock cycle for the VIA.
  * **Timer 1:** Decrements `t1c`. Handles underflow (sets IFR bit 6), optionally reloads from latch (if ACR bit 6 is set), and triggers IRQ check.
  * **Timer 2:** Decrements `t2c` (if not in pulse counting mode - ACR bit 5 is 0). Handles underflow (sets IFR bit 5) and triggers IRQ check.
* **`EmulationRequired()`:** Returns `true` because the timers need to be clocked.
* **`SignalPRA()` / `SignalPRB()`:** Methods to allow the VIA to signal its socket about changes to Port A/B output state (likely called from `WriteByte` for registers `$x0`, `$x1`, `$xF`).
* **`ByteReady()`:** Checks PCR bits to determine serial port readiness.
* **`irqTrigger()`:** Checks IFR against IER and calls `socket.IRQTrigger()` if an enabled interrupt is pending.
* **`irqUpdateMask(data uint8)`:** Implements the logic for updating the IER based on writes to register `$xE`.

## `IVIASocket` Interface

The VIA component relies heavily on an object implementing the `references.IVIASocket` interface, provided via the `Bind` method. This socket acts as the bridge between the generic VIA logic and the specific hardware context it's connected to (e.g., the C1541 board). The socket is responsible for:

* Reading the actual state of physical input pins connected to Port A and Port B (`ReadPRA`, `ReadPRB`).
* Handling the output to physical pins connected to Port A and Port B based on the VIA's internal state (`WritePRA`, `WritePRB`, `WriteDDRA`, `WriteDDRB`). This includes platform-specific logic (like IEC line driving in the C1541).
* Receiving interrupt requests from the VIA and forwarding them to the system's interrupt controller (`IRQTrigger`, `IRQClear`).

## Dependencies

* `github.com/markel1974/c64emu/src/component`
* `github.com/markel1974/c64emu/src/references`

## Integration

The `VIA` component is typically created by the `ComponentFactory` (using `mos6522.NewFactory`). Its `Bind` method is called during the Symphony initialization sequence to link it with its corresponding socket implementation (`VIA1Socket` or `VIA2Socket` in the C1541 context). It participates in the main emulation loop via its `Emulate` method.

## Limitations / TODOs

* Full functionality of the Shift Register (SR), including different modes and external clocking, may require further implementation or verification.
* Handshake line control (CA1, CA2, CB1, CB2) via the Peripheral Control Register (PCR) might be simplified or require specific socket implementation logic.
* Precise cycle-level timing for register access and timer operations might need refinement for maximum accuracy.
* Add comprehensive unit tests.
* Enhance code comments, especially regarding specific register bits, timer modes, and ACR/PCR interactions.