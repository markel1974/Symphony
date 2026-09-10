# The Symphony Virtual File System (VFS) Server

The `kernel/servers/file_system` module provides a powerful Virtual File System designed to map the internal architecture and state of the OS into a navigable, hierarchical structure. 

In traditional UNIX systems, "everything is a file" (like `/proc` or `/sys`). Symphony takes this philosophy and applies it directly to the live Go components and commands running inside the Microkernel.

## Architectural Highlights

### 1. Hierarchical Command Traversal
The VFS manages a tree of `ICommand` nodes. Instead of physical files on a disk, a "directory" in this file system is actually a live Component or a Command Group within the OS.
- **Current Working Directory (`cwd`)**: Just like in a bash shell, the VFS tracks the user's `cwd`. When a user types `cd /sys/emulator/c64`, the VFS resolves the path and updates the state.
- **Dynamic Resolution**: The VFS handles commands like `ls` (`MessageTypeFileSystemCWDDirectoryListingRequest`) by interrogating the live `ICommand` interface of the current node, dynamically listing all exposed properties, methods, and sub-components.

### 2. Autocompletion Engine
Because the file system is built on top of the live reflection/exporter system of Symphony's components, it provides rich, native autocompletion (`MessageTypeFileSystemSuggestionRequest`). 
When a user presses `<TAB>` in the `xsh` shell, the request is routed asynchronously to the VFS server. The VFS analyzes the `cwd`, looks at the exposed methods of the target component, and returns precise autocomplete suggestions (including properties and executable commands) back to the terminal.

### 3. Asynchronous Server Model
Like the `render` server, the VFS is an isolated, event-loop-driven server. It communicates strictly via the Microkernel's message bus. 
If a user requests a deep search across the component tree (`MessageTypeFileSystemFindRequest`), the VFS handles the search asynchronously. This ensures that parsing complex paths or traversing deeply nested hardware component structures never blocks the main kernel loop or the user's rendering thread.

## Why this matters
By treating internal component states and executable Go methods as navigable paths in a Virtual File System, Symphony achieves true "transparency." An engineer doesn't need to write custom debug endpoints or HTTP handlers to inspect a service; they simply `cd` into the running component and execute commands just as if they were navigating a hard drive.
