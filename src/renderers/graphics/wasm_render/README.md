
# Symphony WASM Renderer (`wasm_render`)

This package provides the WebAssembly rendering and input handling capabilities for the Symphony emulation framework. It allows Symphony to run within a web browser, rendering graphics via WebGL (through direct memory access to the WASM module) and capturing user input from JavaScript.

## Overview

The `wasm_render` package consists of three main parts:

1.  **`render.go`:** The core Go/WASM bridge. It exports functions callable by JavaScript to drive the emulation loop, retrieve display data, and pass input. It also implements callbacks for the Symphony `Board` (like VBlank).
2.  **`displaybuffer.go`:** A specialized display buffer optimized for WASM. It converts C64 color indices into an RGBA surface (pixel buffer) that can be efficiently read by JavaScript and passed to WebGL for texture updates. It precalculates color palettes for fast pixel setting.
3.  **`inputs.go`:** Handles the mapping of JavaScript `event.key` strings to Symphony's internal virtual key codes and joystick events, allowing browser keyboard input to control the emulated system.

## Features

* **Browser-Based Emulation:** Enables Symphony (and its emulated systems like the C64) to run directly in modern web browsers.
* **WebGL Rendering:** Utilizes WebGL for rendering the emulated display by directly accessing the frame buffer in WASM memory (via `unsafe.Pointer`).
* **Frame Synchronization:** Employs a `vBlank` flag mechanism to synchronize the Go emulation loop with the browser's `requestAnimationFrame` for smooth rendering.
* **Efficient Pixel Data Transfer:** Exposes a pointer to the raw pixel data (`surface`) in the `DisplayBuffer`, allowing JavaScript to create a `Uint8Array` view directly into WASM memory for optimal performance when updating WebGL textures.
* **Comprehensive Input Mapping:** Translates JavaScript keyboard event strings into appropriate Symphony virtual key codes or joystick actions. Includes a toggle (`joyKeys`) to map arrow keys and Control to either keyboard or joystick 1 input.
* **Optimized Palette Conversion:** The `DisplayBuffer` precalculates C64 color indices to full RGBA values (including an 8-pixel wide version `colors8` for `SetMulti8`) to accelerate pixel setting operations.

## Core Components

### 1. `Render` Struct (`render.go`)

This is the main struct managing the WASM rendering lifecycle.

* **Fields:**
    * `board references.IBoard`: Reference to the main Symphony board.
    * `displayBuffer *DisplayBuffer`: The specialized RGBA surface.
    * `vBlank bool`: Synchronization flag with JavaScript's render loop.
    * `w, h int`: Dimensions of the emulated display.
    * `input *Inputs`: The input handler instance.
* **Key Methods:**
    * `NewRender() *Render`: Constructor.
    * `Setup(board references.IBoard, cfg *config.Config) error`: Initializes the renderer, stores board/config references, and crucially calls `board.Mount(g)` allowing the `Board` to call back methods on this `Render` instance (like `VBlank`, `LedActivity`).
    * `CreateDisplayBuffer(w int, h int) (references.IDisplayBuffer, error)`: Creates and returns the specialized `DisplayBuffer`.
    * `Start() error`: The main Go entry point for the WASM module. It exports Go functions to JavaScript (`js.Global().Set(...)`) and then blocks using an unbuffered channel (`<-c`) to keep the WASM instance alive and responsive to JavaScript calls.
    * **Exported to JS:**
        * `emulateFrame()`: Called by JS's `requestAnimationFrame`. Runs `g.board.Emulate()` until `g.vBlank` is set by the `VBlank` callback.
        * `getSurfacePointer() unsafe.Pointer`: Returns a raw pointer to the pixel data in WASM memory.
        * `getSurfaceLen() int`: Returns the length of the pixel data.
        * `getDisplayBuffer() js.Value`: (Alternative) Returns a copy of the pixel data as a JS `Uint8Array`.
        * `getDisplayWidth() int`, `getDisplayHeight() int`: Return display dimensions.
        * `keyPressed(key string, pressed bool)`, `keyReleased(key string, pressed bool)`: Pass keyboard events to `g.input.Key()`.
    * `VBlank()`: Callback method, called by the `Board` during VBlank. Sets `g.vBlank = true` to release the `emulateFrame` loop.
    * `LedActivity(...)`: Callback for LED state changes (currently a placeholder).

### 2. `DisplayBuffer` Struct (`displaybuffer.go`)

Implements `references.IDisplayBuffer` and prepares C64 pixel data in RGBA format.

* **Fields:**
    * `colors [][4]uint8`: Lookup table mapping 16 C64 colors (expanded to 256 indices) to `[R,G,B,A]` values.
    * `colors8 [][32]uint8`: Optimized lookup table where each entry contains 8 repetitions of an RGBA color (for `SetMulti8`).
    * `surface []byte`: The raw RGBA pixel buffer.
    * `maxLen int`: Length of `surface`.
* **`NewDisplayBuffer(w, h int)`:** Constructor. Initializes color lookup tables and the `surface` buffer.
* **Key Methods:**
    * `GetSurfacePointer() unsafe.Pointer`: Returns a pointer to the `surface` for direct memory access from JS.
    * `Set(idx int, data uint8)`: Sets a single pixel at `idx` by copying an RGBA value from `colors` LUT.
    * `Set8(idx int, data [8]uint8)`: Sets 8 distinct pixels.
    * `SetMulti8(idx int, data uint8)`: Sets 8 pixels to the *same* color using the optimized `colors8` LUT.

### 3. `Inputs` Struct (`inputs.go`)

Handles mapping of JavaScript key events to Symphony's internal input system.

* **Fields:**
    * `board references.IBoard`: Reference to the board to send input events.
    * `keyMapper map[string]func(bool)`: Maps JavaScript `event.key` strings to functions that call `board.KeyboardSetKey()` or `board.Joy1SetKey()`.
    * `joyKeys bool`: Toggles arrow/control keys between keyboard and joystick 1.
* **Key Methods:**
    * `Setup(...)`: Populates the `keyMapper` with extensive mappings for letters, numbers, special keys, and arrow/control keys (with `joyKeys` logic).
    * `Key(keyString string, pressed bool)`: Called by JS. Looks up `keyString` in `keyMapper` and executes the associated function.

## Integration with JavaScript (Example)

The `render.go` module is designed to be called from a JavaScript frontend (like the `test.js` or the HTML file you provided earlier). The typical flow is:

1.  JS loads `symphony.wasm` and `wasm_exec.js`.
2.  JS calls `go.run(wasmInstance)`. This executes `main()` in Go, which eventually calls `Render.Start()`.
3.  `Render.Start()` exports Go functions to JS and then blocks.
4.  JS rendering loop (using `requestAnimationFrame`):
    a.  Calls the exported `emulateFrame()` Go function.
    b.  The Go `emulateFrame` runs `board.Emulate()` until `VBlank()` is called by the board.
    c.  JS retrieves the updated frame buffer pointer/length using exported Go functions (`getSurfacePointer`, `getSurfaceLen`).
    d.  JS creates a `Uint8Array` view into WASM memory: `new Uint8Array(wasmInstance.exports.mem.buffer, surfacePtr, surfaceLen)`.
    e.  JS updates a WebGL texture using `gl.texImage2D` with this `Uint8Array`.
    f.  JS draws the textured quad.
5.  JS keyboard event listeners call the exported `keyPressed` and `keyReleased` Go functions.

## Usage

This renderer is automatically selected and used when Symphony is compiled for the `js/wasm` target. The main application needs to provide an HTML host page with a canvas element and JavaScript glue code to load and run the WASM module, similar to the example HTML you provided.

## Dependencies

* `syscall/js`: For Go-JavaScript interoperation.
* `unsafe`: For `GetSurfacePointer()`.
* Other Symphony core packages (`component`, `config`, `references`).

## Limitations & TODOs

* Mouse input is not currently handled in the `VBlank` or `Inputs` struct.
* Clipboard functionality is commented out.
* LED activity is only logged, not visually represented in WASM.
* Full joystick/gamepad support via browser Gamepad API is a potential future enhancement.