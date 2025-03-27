# Symphony: IComponent Lifecycle and Tree Interaction

This document describes the lifecycle of components implementing the `references.IComponent` interface within the Symphony emulation framework and how they interact within the hierarchical component tree.

## Core Concepts

* **`IComponent`:** The fundamental interface for all managed entities within Symphony (CPU, VIC, SID, CIA, Timers, Boards, etc.). It defines common operations for identification, hierarchy navigation, state management (properties, snapshots), command execution, and basic hardware control (`Reset`, `Emulate`).
* **`BaseComponent`:** A struct embedded by almost all concrete component implementations. It provides default implementations for many `IComponent` methods, particularly those related to the component tree, properties, and commands.
* **`Node` (`references.INode`):** Represents a node within the hierarchical component tree. Each `Node` holds a reference to its associated `IComponent`, its parent `Node`, and a map of its child `Node`s (keyed by the child's `HardwareId`).
* **Component Tree:** The entire emulated system is represented as a tree of `Node`s, with the main `Board` component typically at the root. This structure defines the configuration and hierarchy.
* **Factory (`references.IComponentFactory` & Implementations):** Responsible for creating instances of *specific component types* based on a string identifier (e.g., "mos6510", "cia", "vic") and an instance number. Only "top-level" or independently creatable components have factories registered in the central registry.
* **Snapshot (`map[string]interface{}`):** A hierarchical map representing both the **configuration** (which components exist and their relationships via the tree structure) and the **state** (values of properties) of the emulator. It is the single source of truth for initializing or restoring the system.
* **Sockets (`Connector` Interface & Specific Implementations):** Lightweight structs (e.g., `CPUSocket`, `VICSocket`) usually held by the `Board`. They serve two main purposes:
    1.  Act as intermediaries holding references to the actual components needed for connections.
    2.  Implement the `Connector` interface (`Setup`, `Connect`) to manage the dependency resolution and connection logic in distinct phases.
* **`HardwareId`:** A unique identifier for a component *within its parent's context*, typically `name:instance` (e.g., "cia:0", "cpu:0"). Used as the key in the `Node`'s children map and within the `components` map during setup.
* **`GetId()`:** Returns the *full path* identifier of a component (e.g., "c64.cia:0:ICIA").

## Component Lifecycle Phases

The lifecycle is primarily orchestrated during the emulator's startup, driven by the `component.RestoreAll` function and the `Board.Setup` method.

**Phase 0: Tree Creation & State Restoration (`component.RestoreAll`)**

* **Trigger:** Called once at the very beginning (e.g., in `main.go`), passing `nil` for the parent/component and the snapshot/configuration map (`state`).
* **Mechanism:** Uses the recursive helper function `_restore`.
* **Process (`_restore`):**
    1.  **Check Existence:** Is the `component` parameter `nil`?
        * **Yes (Component needs creation):**
            * Extract the component's full ID (`currentIdKey`) from the `state` map keys.
            * Extract the component `type` (e.g., "cia", "mos6510") and `instance` number from the `details` section within the `state` map corresponding to `currentIdKey`.
            * Call `factory.Create(parentComponent, type, instance)` to instantiate the component.
            * The component's constructor (`New...`) is called:
                * It initializes the component's internal fields.
                * It calls `component.NewBaseComponent(...)`, passing the `factory`, `parentComponent`, `name`, `instance`, the component itself (`m`), and its specific interface type (`references.IdICIA(m, instance)`).
                * `NewBaseComponent` sets the `id`, `name`, `instance`, `kind`, and `factory` fields. It then calls `component.Register`.
            * `component.Register`:
                * Creates a new `Node` using `newNode(parentComponent.GetNode(), component)`.
                * Sets the component's internal `node` reference using `component.SetNode()`.
                * Adds the new component's node to the parent's `children` map using `parentComponent.GetNode().AddComponent()`. This builds the tree structure.
        * **No (Component already exists):**
            * Uses the component passed as a parameter.
            * Retrieves the component's state sub-map from the main `state` map using `component.GetId()`.
    2.  **Restore Properties:** Finds the `"properties"` section in the component's state map and calls `component.Restore(propertiesMap)` to set the component's internal state.
    3.  **Recurse on Children:** Finds the `"children"` section in the component's state map. It iterates through the **keys (child IDs)** and **values (child state maps)** listed *in the snapshot*. For each child entry:
        * It looks up the *existing* child component in the current component's node using `component.GetChild(childId)`. This *will find* "private" children created by the parent's constructor, and *might find* other children if `RestoreAll` is re-run on an existing tree.
        * It calls `_restore` recursively, passing the `factory`, the *current* component as the parent, the *found child* (or `nil` if not found yet), and the specific state map for that child.

* **Outcome:** A complete component tree reflecting the snapshot structure is created in memory, and the state of each component (its properties) is restored. No inter-component connections (via sockets) are established yet. `RestoreAll` returns the root `IComponent` (the `Board`).

**Phase 1: Socket Setup (`Board.Setup` -> `Connector.Setup`)**

* **Trigger:** Called within `Board.Setup` *after* `RestoreAll` has completed.
* **Mechanism:**
    1.  `Board.Setup` creates all necessary *socket* instances (e.g., `NewCPUSocket()`).
    2.  `Board.Setup` creates a temporary map (`components`) mapping `HardwareId` to `IComponent` by traversing the newly built tree (e.g., using `s.GetNode().GetComponents()`).
    3.  `Board.Setup` iterates through a list (`connections`) containing all its created sockets (each implementing the `Connector` interface).
    4.  For each `Connector`, it calls `connector.Setup(components, cfg)`.
* **Process (`Connector.Setup` - e.g., `CPUSocket.Setup`):**
    * Receives the global `components` map.
    * Uses the map to *find* the specific components it depends on by their `HardwareId` (e.g., `"cpu:0:I6510"`, `"pic:0:IPIC6510"`, `"pla:0:IPLAc64"`).
    * Performs a *type assertion* to the required *interface* (e.g., `references.I6510`).
    * Stores these interface references in its *own internal fields*.
    * Performs any configuration specific to the socket itself using `cfg`.
* **Outcome:** All sockets now hold valid references (as interfaces) to the actual component instances they need to interact with.

**Phase 2: Component Connection (`Board.Setup` -> `Connector.Connect`)**

* **Trigger:** Called within `Board.Setup` *after* all sockets have run their `Setup` phase.
* **Mechanism:**
    1.  `Board.Setup` iterates through the same list (`connections`) of `Connector`s.
    2.  For each `Connector`, it calls `connector.Connect()`.
* **Process (`Connector.Connect` - e.g., `CPUSocket.Connect`):**
    * Uses the component interface references stored during its `Setup` phase.
    * Calls the *actual component's* `Setup` method (e.g., `cpu.Setup(socket)`), passing the socket itself or other necessary interfaces obtained during the socket's `Setup`. This establishes the crucial link allowing the component to call back to the socket (and thus indirectly to other components).
* **Outcome:** All components are now fully connected and initialized, ready for emulation. Dependencies are resolved.

**Phase 3: Build Emulation Sequence (`Board.Setup` -> `rebuildEmulation`)**

* **Trigger:** Called within `Board.Setup` *after* the `Connect` phase.
* **Mechanism:**
    1.  `rebuildEmulation` iterates through a predefined sequence (`_hardwareSequence`).
    2.  For each component ID in the sequence, it finds the corresponding component in the `components` map.
    3.  It checks `component.NeedsEmulation()`.
    4.  If true, it appends the component's `Emulate` method (as a function pointer) to the `Board`'s `emulation` slice.
* **Outcome:** The `Board.emulation` slice contains an ordered list of functions to call for each emulation step, optimized to include only active components.

**Phase 4: Initial Reset (`Board.Setup` -> `board.reset`)**

* **Trigger:** Called at the end of `Board.Setup`.
* **Mechanism:** Calls the `Reset` method on key sockets (which typically delegate to the connected component's `Reset` method).
* **Outcome:** Ensures the emulated system starts in a known, clean state.

**Phase 5: Runtime Emulation (`Board.Emulate`)**

* **Trigger:** Called repeatedly by the main application loop (e.g., managed by the renderer).
* **Mechanism:** Simply iterates through the pre-built `emulation` slice and calls each function pointer (`f()`).
* **Outcome:** Executes one step (e.g., one clock cycle or one phase) of the emulation for all active components in the correct order.

**Interaction within the Tree (Post-Setup):**

* **Via Sockets:** The `Board` interacts with major components via the methods defined in the *socket structs* (which embed the component interfaces).
* **Component-to-Component:** Components interact with each other *indirectly* via the interfaces passed to them during their `Setup` phase (which were originally retrieved by the sockets). For example, the CPU calls `banksRead` which points to the `PLA`'s `Read` method.
* **Introspection:** Any external tool (like the console) can use `FindNode` or `GetComponentPath` to locate any component and then use `GetProperty`/`SetProperty`/`CommandExec` to interact with it, traversing the tree as needed.

This multi-phase initialization ensures maximum decoupling and flexibility, allowing the system structure to be defined entirely by the snapshot while resolving dependencies correctly before emulation begins.