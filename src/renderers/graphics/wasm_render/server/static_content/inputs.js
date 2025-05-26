
function setupInputs(wasmInstance) {
    document.addEventListener("keydown", (event) => {
        keyPressed(event.key);
    });
    document.addEventListener("keyup", (event) => {
        keyReleased(event.key);
    });
}