# Hardware Component Factory

## Overview

This package provides the central factory responsible for instantiating all hardware components within the Symphony emulation framework. It acts as a master constructor that can build any registered hardware component based on a type identifier string, effectively decoupling the system's assembly logic from the concrete implementations of the components themselves.

This factory is a cornerstone of Symphony's extensible and snapshot-driven architecture.

## Architectural Role

The primary role of this `Factory` is to enable a fully **data-driven approach to machine construction**. Instead of the main application having to know about every possible chip and peripheral, it simply asks this factory: "please create a new component of type 'cpu_6510'".

This provides several key advantages:
* **Extensibility:** New hardware components can be added to the framework without ever modifying the core factory or system assembly logic.
* **Decoupling:** The configuration of a machine (defined in a snapshot file) is completely decoupled from the component implementation code.
* **Centralized Service Provision:** The factory also acts as a provider for global services that components might need, such as an `IDisplayBuffer` or an `IAudioRender`, ensuring that all components have access to the same shared resources.

## Mechanism of Operation: A Registry-Based Factory

Symphony uses a powerful and idiomatic Go pattern to make its factory system dynamic and decentralized.

1.  **Central Registry (`/registry`):** A global registry holds a map that associates a component type `string` (e.g., `"c64_pla"`) with a specific component factory object that knows how to build it.

2.  **Self-Registration:** Each hardware component package (e.g., `/hardware/vic`, `/hardware/sid`) is responsible for its own factory. Inside each package, a special Go `init()` function is used to automatically register its factory with the central registry when the application starts.

    *Example from a component's `factory.go` file:*
    ```go
    func init() {
        // The VIC factory registers itself with the central registry
        registry.RegisterComponentFactory(NewFactory())
    }
    ```

3.  **Delegation:** When the main `hardware.Factory`'s `Create` method is called, it simply performs a lookup in the central registry using the provided type string. It finds the correct sub-factory and delegates the creation task to it.

This self-registration mechanism means the core framework doesn't need a list of all possible components; it discovers them dynamically at startup.

## How to Add a New Component to the Framework

Thanks to this architecture, adding a new hardware component to Symphony is a clean and straightforward process:

1.  **Create the New Component Package:** For example, `hardware/z80/`.
2.  **Implement the Component:** Write your `z80.go` file, ensuring it fulfills the `IComponent` interface contract from the `/references` package.
3.  **Create a Component-Specific Factory:** In `hardware/z80/factory.go`, create a `Factory` struct that implements the `IFactory` interface. Its `Create` method will simply call `NewZ80(...)`.
4.  **Register the Factory:** Add an `init()` function to `hardware/z80/factory.go`:
    ```go
    func init() {
        registry.RegisterComponentFactory(NewFactory())
    }
    ```
5.  **Update Configuration:** You can now instantiate your new component from a snapshot file by using its registered identifier string (e.g., `"cpu_z80"`).

The framework will automatically know how to build your new component without any changes to the core emulation or assembly logic.