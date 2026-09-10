# Component Framework

This package provides a powerful and flexible framework for building complex, hierarchical, and introspective component-based systems in Go. It is the foundational layer of the Symphony project, designed to create modular, dynamic, and easily debuggable applications like emulators, simulators, or any system composed of interconnected parts.

---

## Core Philosophy

The `component` framework was designed to overcome the limitations of static, monolithic architectures. The core philosophy is built on three pillars:

1.  **Everything is a Component:** Every logical part of the system, from a CPU to a single timer, is a self-contained `IComponent`. This promotes extreme modularity and reusability.
2.  **Introspection is a First-Class Citizen:** A system should be transparent and modifiable *at runtime*. This framework is designed from the ground up to allow any component's state (Properties) and behavior (Commands) to be inspected and invoked dynamically, which is invaluable for debugging and interactive control.
3.  **Convention Over Configuration, with Automation:** Boilerplate code is a source of bugs and slows down development. This framework uses a simple convention (comment tags) and a powerful code generator (`exporter`) to automate the tedious work of creating accessors and registrations, letting developers focus on the core logic of their components.

---

## Key Concepts

The framework is composed of several key parts that work together to create a cohesive system.

### 1. The `IComponent` Interface and `BaseComponent`

At the heart of the framework is the `IComponent` interface and its concrete implementation, `BaseComponent`. By embedding `BaseComponent` into your own structs, you instantly provide them with a rich set of built-in functionalities:

-   **Identity:** A unique, hierarchical ID (e.g., `/board/cpu0`).
-   **Hierarchy:** Parent-child relationships, managed through a `Node` system.
-   **State Exposure:** A system for registering and managing dynamic **Properties**.
-   **Behavior Exposure:** A system for registering and executing dynamic **Commands**.
-   **Shell Integration:** Automatic creation of shell commands for introspection (`dump`) and property access.

### 2. The Component Tree (`node.go`)

Components are not standalone; they exist within a hierarchical tree structure. The `Node` object wraps each `IComponent` and manages the parent-child relationships. This allows for:

-   **Logical Organization:** Representing the natural structure of a complex system (e.g., a `Board` contains a `CPU` and a `VIC`).
-   **Traversal:** Finding any component in the system using a path-based query (e.g., `GetComponentPath("/board/cia1")`).
-   **Scoped Operations:** Executing commands or getting properties on components deep within the hierarchy.

### 3. Properties (`properties.go`)

Properties are how a component exposes its internal state in a controlled manner. A `PropertyInfo` object encapsulates a field's metadata and its access functions (getter and setter).

-   **Dynamic Access:** Properties are accessed by string name, making them ideal for use in an interactive shell or scripting environment.
-   **Type Safety:** The framework uses reflection to ensure that `Set` operations are type-safe. It intelligently handles type conversions from strings or other compatible types.
-   **Read-Only Support:** Properties can be marked as read-only, preventing accidental modification of critical state.

### 4. Commands (`commands.go`)

Commands are how a component exposes its behavior. A `Command` object wraps a Go function, making it dynamically executable.

-   **Dynamic Invocation:** Like properties, commands are executed by their string name.
-   **Automatic Argument Conversion:** When a command is executed from a string-based source (like a shell), the framework automatically parses and converts the string arguments into the types required by the Go function's signature (`int`, `bool`, `string`, etc.).
-   **Strict Signatures:** The system validates function signatures at creation time using `panic`. This is a deliberate design choice: an invalid signature is considered a **programming error** that must be fixed during development, not a runtime error to be handled gracefully.

---

## The Code Generator (`exporter.go`)

The `exporter` is the key to the framework's productivity. It's a powerful tool that automates the creation of all the boilerplate code required to expose properties and commands.

### The "Opt-In" Approach: `// symphony:export`

To expose a field as a property, you simply add a special comment tag above it.

**Example:**
```go
type MyComponent struct {
    *component.BaseComponent
    
    // symphony:export myCounter is a crucial value. readonly
    myCounter int
    
    internalState string // This field will NOT be exposed
}