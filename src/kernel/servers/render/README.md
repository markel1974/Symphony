# The Symphony Render Server: A Compositing Window Manager for the Terminal

This module (`kernel/servers/render`) implements a fully asynchronous, compositing window manager designed specifically for Text User Interfaces (TUI) and terminal environments. 

Instead of allowing user processes to write directly to the standard output (which causes race conditions, visual artifacts, and prevents windowing), Symphony forces all graphical output through this dedicated rendering server via the Microkernel's message bus.

## Architectural Highlights

### 1. Asynchronous Compositing
The Render server behaves similarly to modern graphical compositors (like Wayland or X11), but adapted for text grids. 
User processes (like the `xsh` shell, the `xvi` text editor, or the C64 emulator) do not draw to the screen. Instead, they send rendering requests (e.g., `MessageTypePaintRequest`, `MessageTypeWriteColor`) over the asynchronous kernel bus. The Render server buffers these requests into independent virtual `Surface` objects.

### 2. Z-Index and Window Management
Because the TUI is composited, the Render server supports true window management:
- **Z-Indexing**: Every running process is assigned a `zIndexCounter`. Background processes continue to execute and update their virtual surfaces, but only the active foreground surfaces are painted on top.
- **Offsets and Scaling**: Surfaces support `offsetX`, `offsetY`, and `scale`. A process can draw a 80x24 interface, but the Render server can scale it down or move it to a specific quadrant of the terminal, allowing for split-screen layouts or picture-in-picture debugging.
- **Window Selection**: The server includes a native `WindowSelector` mode, allowing the user to seamlessly switch focus between running processes (similar to Alt-Tab), bringing different virtual surfaces to the foreground.

### 3. Decoupled Display Driver
The `Render` struct relies on an `interfaces.IDisplayDriver`. This means the compositing engine itself has no idea *how* the pixels or text characters are actually put on the screen. It can render to a local raw TUI, a remote SSH terminal, a WebSockets bridge, or a graphical SDL/OpenGL context, simply by injecting a different driver implementation.

## Why use a Render Server?
In a multitasking OS, multiple processes running in the background might panic, log errors, or try to update progress bars simultaneously. A compositing render server ensures that:
1. The foreground application (e.g., a text editor) is never visually corrupted by a background application.
2. Background applications don't block waiting for I/O; they simply update their hidden virtual surfaces at maximum speed.
3. The host terminal is only refreshed at a controlled frame rate (using an adaptive ticker), eliminating flicker and reducing bandwidth overhead over SSH.
