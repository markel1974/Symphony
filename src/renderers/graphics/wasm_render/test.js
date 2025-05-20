// Carica il modulo WASM.
const go = new Go();
WebAssembly.instantiateStreaming(fetch("symphony.wasm"), go.importObject).then((result) => {
    go.run(result.instance);

    // Ottieni un riferimento al canvas.
    const canvas = document.getElementById("canvas"); // Assumi che ci sia un elemento <canvas> con id="canvas"
    const ctx = canvas.getContext("2d");

    // Funzione per disegnare un singolo carattere.
    function drawChar(char, fgColor, bgColor, x, y) {
        // ... (implementazione: usa ctx.fillStyle, ctx.fillRect, ctx.fillText, ecc.) ...
    }

    // Ciclo principale (eseguito a ogni frame).
    function renderFrame() {
        // Esegui un frame di emulazione.
        go.exports.emulateFrame();

        // Ottieni il display buffer.
        const buffer = go.exports.getDisplayBuffer();

        //INTEGRAZIONE PER IL DISEGNO DEL BUFFER
        // Disegna il display buffer.

        // Richiedi il prossimo frame.
        requestAnimationFrame(renderFrame);
    }

    // Gestisci gli eventi di tastiera.
    document.addEventListener("keydown", (event) => {
        go.exports.keyPressed(event.keyCode); // Invia il codice del tasto al modulo WASM
    });

    // Avvia il ciclo principale.
    renderFrame();
});