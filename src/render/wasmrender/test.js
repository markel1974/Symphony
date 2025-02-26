// Carica il modulo WASM.
const go = new Go();
WebAssembly.instantiateStreaming(fetch("g64.wasm"), go.importObject).then((result) => {
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
        const uint8array = new Uint8Array(buffer);
        // Itera sui dati nel buffer e disegna i caratteri sul canvas.
        for (let y = 0; y < 25; y++) {
            for (let x = 0; x < 40; x++) {
                const charIndex = (y * 40 + x) * 2;
                const charCode = uint8array[charIndex]
                const colorCode = uint8array[charIndex+1]

                const char = String.fromCharCode(charCode); // Dovresti usare una tabella di conversione PETSCII -> Unicode!
                const fgColor = convertC64ColorToCSS(colorCode & 0x0F); // Funzione di conversione da implementare
                const bgColor = convertC64ColorToCSS((colorCode >> 4) & 0x0F);

                //TODO
                //document.getElementById("canvas").getContext("2d").putImageData(new
                //ImageData(Uint8ClampedArray.from(x), 1, 1), 1, 1);

                drawChar(char, fgColor, bgColor, x * 8, y * 8); // Assumendo una griglia di 8x8 pixel per carattere
            }
        }

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