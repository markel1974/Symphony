# Symphony: A Configurable Emulation Framework

Symphony is a powerful and flexible emulation framework written in Go. It's designed to be a platform for building accurate and highly configurable emulators of various computer systems. While currently focused on emulating the Commodore 64 (with full 1541 disk drive support), the architecture is inherently extensible to other platforms.  Symphony emphasizes deep, dynamic introspection, allowing users and developers to inspect and modify the emulated system's state in real time.

## Key Features

*   **Modular Architecture:** Built upon a robust component-based architecture, with well-defined interfaces (`IComponent`, and specific socket interfaces) for maximum flexibility and extensibility.  Components are organized in a hierarchical tree structure.
*   **Dynamic Introspection:**  A unique and powerful introspection system allows you to examine and modify *any* property of *any* component at runtime. This is accessible through a built-in, interactive console.
*   **Snapshot-Driven Configuration:**  The entire system configuration (which components are present, their properties, and their interconnections) is defined by a snapshot.  Loading a snapshot effectively *creates* or *reconfigures* the emulated system.  This enables easy switching between different hardware setups (e.g., C64 vs. VIC-20) and facilitates a potential future WYSIWYG configuration UI.
*   **Cycle-Accurate VIC-II Emulation:** Symphony aims for cycle-accurate emulation of the Commodore 64's VIC-II video chip, ensuring accurate rendering of complex graphical effects.
*   **Complete 1541 Emulation:** Includes a full emulation of the 1541 disk drive, including its 6502 CPU, VIA chips, and drive mechanics. This allows for accurate emulation of software that relies on 1541-specific features (fast loaders, copy protection, etc.).
*   **Integrated Debugging Console:** A powerful, text-based console (accessible locally or remotely via SSH) provides access to the introspection system, allows for command execution, and even supports multiple simultaneous processes and windows within the console itself. It includes features like a simple text-mode graphing capability for visualizing performance metrics.
*   **Go Language:** Written entirely in Go, leveraging the language's strengths in concurrency, safety, and performance.
*   **No External Dependencies:**  The core emulation logic has no external dependencies (beyond the Go standard library), simplifying building and distribution.  (Note: Specific renderers, like the OpenGL renderer, *do* have external dependencies.)
*   **Open Source (LGPL v2.1):**  Symphony is intended to be released under the LGPL v2.1 license, encouraging community contributions and use in other projects.

## Architecture Overview

Symphony's architecture revolves around the concept of interconnected *components*.  Every emulated element, from the CPU to a single timer within a CIA chip, is represented as an `IComponent`.

*   **`IComponent` Interface:** The foundation of the system.  All components implement `IComponent`, providing a consistent interface for:
    *   `GetId()`:  Retrieving a unique identifier (e.g., "c64.cpu", "c64.iec.c1541.8.mos6502").  The ID is a hierarchical path reflecting the component's position in the tree.
    *   `GetNode()`: Accessing the component's position within the hierarchy.
    *   `GetProperty(string)`/`SetProperty(string, interface{})`:  Getting and setting component properties (registers, flags, memory locations, etc.) *by name*. This is the core of the dynamic introspection system.
    *   `Dump()`/`Restore(map[string]interface{})`: Saving and restoring the component's state (used for snapshots and configuration).
    *   `CommandAdd(...)`, `CommandExec(...)`:  Adding and executing custom commands specific to the component.
    *    `Kind() string`: returns the type.
*   **`BaseComponent`:** A struct that provides a default implementation of `IComponent`, handling common tasks like ID management, node connections, and property introspection.  All components embed `BaseComponent`.
*   **`Node`:** Represents a node in the component tree. Handles parent/child relationships and provides methods for traversing the tree (e.g., `FindNode`).
*   **`ComponentFactory`:** A centralized factory responsible for creating component instances based on a string identifier (e.g., "mos6510", "mos6569"). This allows the system configuration to be data-driven.  The factory uses a map of specific factory functions for each component type.
*   **Sockets:** Specialized interfaces (e.g., `CPUSocket`, `VICSocket`, `CIASocket`) define the *specific* interactions between the `Board` and the core components.  This decouples the `Board` from the concrete component types.  These sockets *embed* the relevant component interface (e.g., `CPUSocket` embeds `I6510`).
*   **`Properties` and `PropertyInfo`:** These structures, used by `BaseComponent`, provide the dynamic introspection capabilities.  They allow access to component properties by name (string), while still maintaining type safety at runtime.
* **Dispatcher:**

**Initialization and the Role of `RestoreAll`:**

Symphony uses a unique approach to initialization.  Instead of creating all components directly in the `Board`'s `Setup` method, the system uses the `component.RestoreAll` function:

1.  **Snapshot Loading:**  `RestoreAll` takes a map[string]interface{} as input.  This map represents either:
    *   A previously saved snapshot of the emulator's state.
    *   A *minimal initial state* defining the desired system configuration (e.g., a map containing just the "c64" root node, if no snapshot is provided).

2.  **Recursive Component Creation:** `RestoreAll` recursively traverses the provided map, which mirrors the desired component tree structure.  For *each* component description found in the map:
    *   If a component with the corresponding ID already exists in the tree, its state is restored from the map.
    *   If a component with the corresponding ID does *not* exist, `RestoreAll` uses the `ComponentFactory` to *create* a new instance of the appropriate component type.  The factory is consulted based on the "type" field within the snapshot data.

3.  **Component Connection:** *After* `RestoreAll` has created (or restored) the entire component tree, the `Board`'s `Setup` method is called.  The `Setup` method then:
    *   Creates the *socket* objects (which are lightweight interfaces).
    *   Uses `GetNode().GetChild` to *locate* the key components within the tree (e.g., CPU, VIC, etc.).
    *   Uses the socket's `Connect` method to establish the connections between the `Board` and the components, passing the specific component interfaces.

This approach means that the *same mechanism* is used for both initial system creation *and* for restoring a previously saved state.  The snapshot effectively *defines* the system configuration.

## Getting Started

[**TODO:** Add detailed instructions on how to build and run the emulator. Include:]

*   **Prerequisites:**  Go version, any necessary tools.
*   **Cloning the repository:** `git clone ...`
*   **Building:** `go build ...`
*   **Running:** `./symphony [options]`
*   **Basic usage examples:**  Loading a PRG file, loading a disk image, etc.

## Console Usage

[**TODO:** Provide a detailed guide to using the console, including:]

*   Connecting to the console (SSH).
*   Basic navigation (`cd`, `ls`, `pwd`).
*   Inspecting properties (`get`).
*   Modifying properties (`set`).
*   Running commands.
*   Using the built-in help (`help`).
*   Examples of common debugging tasks.
*   Window.

## Contributing

[**TODO:** Add contribution guidelines. This should cover:]

*   How to report bugs (using GitHub Issues, ideally).
*   How to suggest new features (GitHub Issues or a dedicated forum/discussion board).
*   How to submit code changes (pull requests).
*   Coding style guidelines (mention `go fmt`, linters, etc.).
*   Testing guidelines (how to write and run tests).
*   Any specific requirements for commit messages.

## License

This project is licensed under the LGPL v2.1 License - see the [LICENSE](LICENSE) file for details.