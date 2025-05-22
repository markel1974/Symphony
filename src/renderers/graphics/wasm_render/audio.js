// ... nel tuo <script> dopo aver caricato symphony.wasm ...

let audioCtx;
let nextStartTime = 0;     // Tiene traccia di quando il prossimo buffer dovrebbe iniziare a suonare
let audioQueue = [];       // Coda per i buffer audio ricevuti da Go
let isPlaying = false;
let minSamplesForPlayback = 1024; // Inizia a suonare quando abbiamo almeno questi campioni in coda

function setupAudio() {
    try {
        audioCtx = new (window.AudioContext || window.webkitAudioContext)();
        const sampleRate = audioCtx.sampleRate;

        // Funzione callback che Go chiamerà con i campioni pronti
        const samplesReadyCallback = (jsFloat32Array) => {
            // jsFloat32Array è un Float32Array creato in Go e passato a JS
            if (jsFloat32Array && jsFloat32Array.length > 0) {
                const audioBuffer = audioCtx.createBuffer(1, jsFloat32Array.length, sampleRate);
                audioBuffer.getChannelData(0).set(jsFloat32Array);
                audioQueue.push(audioBuffer);
                schedulePlayback();
            }
        };

        initAudioContextAndGetCallback(audioCtx, sampleRate, samplesReadyCallback);

        // Prova a sbloccare l'AudioContext (richiesto da alcuni browser dopo interazione utente)
        //document.body.addEventListener('click', () => {
        //    if (audioCtx.state === 'suspended') {
        //        audioCtx.resume().then(() => {
        //            console.log('AudioContext resumed on user interaction.');
        //            // Potresti voler chiamare wasmAudioPlay() qui se l'emulazione era in attesa
        //        });
        //    }
        //}, { once: true });
    } catch (e) {
        console.error("Web Audio API is not supported in this browser", e);
        alert("Web Audio API is not supported in this browser");
    }
}

function schedulePlayback() {
    if (audioCtx.state === 'suspended') {
        return;
    }
    let totalQueuedSamples = 0;
    for(const buffer of audioQueue) {
        totalQueuedSamples += buffer.length;
    }
    if (isPlaying && audioQueue.length > 0 || (!isPlaying && totalQueuedSamples >= minSamplesForPlayback)) {
        isPlaying = true;
        while (audioQueue.length > 0) {
            const audioBuffer = audioQueue.shift();
            const source = audioCtx.createBufferSource();
            source.buffer = audioBuffer;
            source.connect(audioCtx.destination);

            const currentTime = audioCtx.currentTime;
            if (nextStartTime < currentTime) {
                nextStartTime = currentTime;
            }
            source.start(nextStartTime);
            nextStartTime += audioBuffer.duration;
        }
    } else if (audioQueue.length === 0) {
        isPlaying = false;
    }
}

// Chiama setupAudio dopo che il WASM è stato inizializzato e go.run() è stato chiamato.
// All'interno del .then() di instantiateStreaming, dopo go.run(wasmInstance);

// --- Esempio nell'HTML ---
// WebAssembly.instantiateStreaming(fetch("symphony.wasm"), go.importObject).then((result) => {
//     const wasmInstance = result.instance;
//     go.run(wasmInstance); // Fa partire il main di Go, che esporterà le funzioni
//
//     // ... Inizializzazione WebGL ...
//
//     setupAudio(); // Chiama questa dopo che il WASM è attivo e le func Go sono esportate
//
//     // ... Loop renderFrame ...
// });