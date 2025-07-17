# Component: MOS 6522 Versatile Interface Adapter (VIA)

## 1. Design Philosophy: A Reusable Hardware Component

This package provides a high-fidelity, component-based implementation of the **MOS Technology 6522 Versatile Interface Adapter (VIA)** chip. It is designed from the ground up to function as a modular, reusable "building block" within the Symphony emulation framework.

The core principle of this implementation is the strict **separation of the generic chip logic from its specific hardware context**. The `VIA` component itself is "agnostic"—it perfectly emulates the internal registers, timers, and logic of a real 6522, but it has no knowledge of what it is connected to.

Its purpose and function are only defined at runtime when it is "plugged into" a specific **Socket**. This architectural choice makes the `VIA` component a true virtual chip, ready to be integrated into any simulated hardware board that requires its functionality.

---
## 2. The Virtual Assembly Process

The `VIA` component is designed to be part of a "virtual assembly" line, managed by the `Board` that hosts it.

1.  **Instantiation**: A generic `VIA` component is created by its factory. At this point, it is a standalone object.
2.  **The Socket (`IMos6522Socket`)**: The `Board` creates a specific "socket" (e.g., a `VIA1Socket` or `VIA2Socket`). This socket implements the `IMos6522Socket` interface and contains all the logic that is specific to that VIA's role on the board—controlling the IEC bus, managing drive mechanics, etc.
3.  **Binding (`Bind`)**: The assembly is completed when the `VIA` component's `Bind` method is called, passing in the socket. This "plugs" the generic chip into its specific role, connecting its internal logic to the outside world through the socket's implementation.

---
## 3. Core Features

The `VIA` component emulates the primary features of the 6522 chip:

* **Parallel I/O Ports**: Full emulation of Port A (PRA) and Port B (PRB), including their Data Direction Registers (DDRA, DDRB). All reads from and writes to the physical "pins" are delegated to the bound socket.
* **16-bit Timers**:
  * **Timer 1**: A 16-bit counter supporting one-shot and continuous (free-running) modes. It can generate interrupts on underflow and is capable of toggling output on PB7 (delegated to the socket).
  * **Timer 2**: A 16-bit counter that can operate in one-shot timed mode or in pulse-counting mode, decrementing on negative edges of the PB6 pin (read via the socket).
* **Shift Register (SR)**: A basic 8-bit shift register with logic for handling serial I/O, clocked by Timer 2 or the system clock.
* **Interrupt Control**: Complete emulation of the Interrupt Flag Register (IFR) and Interrupt Enable Register (IER). The `VIA` sets flags internally based on events (timer underflow, port edges) and relies on the socket's `IRQTrigger()` and `IRQClearTrigger()` methods to signal the system's CPU.
* **Handshake Control**: Edge detection on the `CA1` and `CB1` control lines is fully implemented, allowing these signals to set interrupt flags and latch the input on their respective ports.

---
## 4. The `IMos6522Socket` Interface: The Contract

The `VIA` component is entirely dependent on an external object that implements the `IMos6522Socket` interface. This interface is the "contract" that defines how the virtual chip connects to the main board.

The socket is responsible for:

* **Reading Pin State**: Providing the current state of the physical input pins (e.g., `ReadPortA()`, `ReadPortB()`, `ReadCA1()`).
* **Writing to the "World"**: Handling the logic for when the `VIA` writes to a port (e.g., `SignalPRA()`, `SignalPRB()`). This is where generic register writes are translated into specific actions, like controlling a motor or changing a line on the serial bus.
* **Managing Interrupts**: Forwarding interrupt requests from the `VIA` to the system's interrupt controller (`IRQTrigger()`, `IRQClearTrigger()`).

---
## 5. Key API

### `NewVIA(...) *VIA`
The component factory calls this constructor to create a new, generic `VIA` instance.

### `Bind(socket references.IMos6522Socket)`
This is the crucial assembly method. It connects the generic `VIA` logic to a specific socket implementation, giving the chip its purpose.

### `Reset()`
Resets all internal registers (`PRA`, `PRB`, `DDRA`, `DDRB`, `ACR`, `PCR`, `IFR`, `IER`) to zero and resets its internal `Timer` and `ShiftRegister` components.

### `ReadByte(addr uint16) uint8`
Handles CPU read requests from the VIA's 16 memory-mapped registers. It returns internal register values or delegates to the socket for port reads. Reading from specific timer or flag registers automatically clears the relevant interrupt flags as per hardware behavior.

### `WriteByte(addr uint16, data uint8)`
Handles CPU write requests. It updates internal latches, control registers, and DDRs. Writes to output registers like `PRA` or `PRB` are immediately passed to the socket via methods like `SignalPRB()` to affect the external hardware.

### `Emulate()`
Called on every system clock cycle. This method decrements the timers, checks for underflow conditions, handles the shift register logic, and manages handshake line edge detection. It is the core of the chip's "live" behavior.