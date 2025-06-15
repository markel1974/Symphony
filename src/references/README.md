# References & Interfaces Package

## Architectural Role

This package is the **architectural heart** of the Symphony framework. It contains no concrete logic, only Go `interface` definitions. These interfaces serve as the "contracts" or "schematics" that define how every component in a virtual machine must interact.

The existence of this package is a deliberate design choice based on a core software engineering principle: **Dependency Inversion**. In Symphony, components do not depend on each other directly; they depend on the abstract contracts defined here.

For example, the `CPU` component does not know what a `PLA` is. It only knows that it is connected to a component that fulfills the `IC64Pla` interface contract. This is the key to the framework's modularity and extensibility.

## Key Benefits of This Approach

1.  **Total Decoupling:** Components are completely decoupled. The `CPU`'s logic can be developed and tested independently of the `PLA`, as long as both adhere to the interface contract that binds them.

2.  **Maximum Extensibility:** This architecture makes adding new or alternative hardware incredibly simple. For example, to add a new type of sound chip, one only needs to create a new component that implements the `ISID` interface. The rest of the system can use it without any changes, because it only ever communicates with the `ISID` contract, not the concrete implementation.

3.  **Enhanced Testability:** It allows for the creation of "mock" components for unit testing. A component like the `CPU` can be tested in isolation by providing it with a mock `IC64Pla` that simulates memory access, without needing to initialize the entire C64 system.

4.  **Self-Documenting Architecture:** The interfaces in this package serve as the ultimate technical documentation. By reading an interface like `IC64Expansion`, a developer can understand exactly what capabilities a cartridge has and how it is allowed to interact with the rest of the system, because the interface mirrors the physical pins of the real hardware's expansion port.

## Key Interfaces

This package defines the contracts for all major subsystems, including but not limited to:

* **`IComponent`:** The foundational interface that every single object in the component tree must implement. It provides core functionalities for identity, navigation, configuration, and state management.
* **`IC64Board`:** The contract for a "motherboard" component, which is responsible for orchestrating the lifecycle (setup, connection, emulation ticks) of all other components.
* **`IMos6510` (`i6510.go`):** Defines the contract for a CPU, including methods for executing cycles and handling interrupts.
* **`IC64Pla` (`ic64_pla.go`):** Defines the contract for the C64's memory mapping unit, exposing methods for reading and writing to the bus.
* **`IExpansionC64` (`ic64_expansion.go`):** A masterclass in design, this interface acts as a virtual "expansion port connector," defining every valid hardware signal a cartridge can send or receive.
* **`IFactory` (`ifactory.go`):** The contract for the component factories, ensuring the central registry can delegate the creation of any component in a uniform way.

In summary, this package is the constitution of the Symphony universe. It enforces the rules, defines the relationships, and makes the entire framework robust, clean, and a pleasure to extend.