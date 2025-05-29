let _audioCtx; // Il tuo AudioContext, assicurati sia inizializzato e 'running'
//const audioQueue = []; // La coda dove Go aggiunge gli AudioBuffer
//let nextScheduleTime = 0; // L'AudioContext time in cui il prossimo buffer dovrebbe iniziare
//let isSchedulerActive = false; // Flag per sapere se il loop dello scheduler è attivo
//let schedulerTimerId = null; // ID del timer per setTimeout
let _audioHasUserInteracted = false;
let _audioSampleRate = 0;
let _audioNextStartTime = 0;

//let _audioQueue = [];
//let _audioIsPlaying = false;

// --- Parametri configurabili per lo scheduler ---
//const SCHEDULE_AHEAD_TIME = 0.1; // Quanto avanti schedulare (es. 100ms)
//const SCHEDULER_INTERVAL_MS = 10; // Ogni quanto far girare lo scheduler (es. 25ms)
//const MIN_START_OFFSET = 0.020;   // Piccolo offset per il primissimo avvio (20ms)
//const SAFETY_OFFSET_ON_LAG = 0.005; // Piccolo offset se siamo in ritardo (5ms)

function setupAudio(wasmInstance, sampleRate) {
    try {
        _audioCtx = new (window.AudioContext || window.webkitAudioContext)({ sampleRate: sampleRate })
        //audioCtx = new (window.AudioContext || window.webkitAudioContext)();
        _audioSampleRate = _audioCtx.sampleRate;
        wasmAudioInit(_audioCtx, _audioSampleRate, samplesReadyCallback);

        document.body.addEventListener('click', async () => {
            if (!_audioHasUserInteracted && _audioCtx && _audioCtx.state === 'suspended') {
                try {
                    await _audioCtx.resume();
                    console.log("JS: AudioContext resumed by user gesture. State:", _audioCtx.state);
                    wasmAudioFlush();
                    _audioHasUserInteracted = true;
                    //const resumedTime = audioCtx.currentTime;
                    _audioNextStartTime = 0 //resumedTime + 0.020;
                } catch (e) { console.error("JS: Error resuming AudioContext on gesture", e); }
            }
        }, { once: true });
        //document.body.addEventListener('click', onUserAudioInteraction, { once: true });
    } catch (e) {
        console.error("Web Audio API is not supported in this browser", e);
        alert("Web Audio API is not supported in this browser");
    }
}

function samplesReadyCallback(jsFloat32Array)  {
    if (!_audioCtx) {
        console.error("JS samplesReadyCallback: AudioContext non inizializzato.");
        return;
    }
    if (_audioCtx.state !== 'running') {
        console.warn("JS samplesReadyCallback: AudioContext non 'running'. Stato:", _audioCtx.state, ". Buffer scartato.");
        return;
    }
    if (!jsFloat32Array || jsFloat32Array.length === 0) {
        return;
    }
    const audioBuffer = _audioCtx.createBuffer(1, jsFloat32Array.length, _audioSampleRate);
    audioBuffer.getChannelData(0).set(jsFloat32Array);
    const source = _audioCtx.createBufferSource();
    source.buffer = audioBuffer;
    source.connect(_audioCtx.destination);
    const currentTime = _audioCtx.currentTime;
    let startTime = _audioNextStartTime;
    if (startTime === 0) {
        //console.log(`JS: First ever buffer, scheduling at currentTime + 0.010s`);
        console.log(`JS: First ever buffer, scheduling at currentTime`);
        startTime = currentTime //+ 0.010;
    } else if (startTime < currentTime - 0.1) {
        // Se siamo MOLTO in ritardo (100ms)
        console.warn(`JS: Major lag detected. Resetting start time. Ideal: ${startTime.toFixed(4)}, Current: ${currentTime.toFixed(4)}. Starting in 10ms.`);
        _audioNextStartTime = 0
        return;
        //startTime = currentTime + 0.010;
    } else if (startTime < currentTime) {
        // Lag minore, schedula al tempo ideale passato
        // const lag = currentTime - startTime;
        // console.log(`JS: Scheduling with slightly past time. Ideal: ${startTime.toFixed(4)}, Current: ${currentTime.toFixed(4)}, Lag: ${lag.toFixed(4)}s`);
        // Non c'è bisogno di modificare startTime, Web Audio lo gestirà.
    }
    source.start(startTime);
    _audioNextStartTime = startTime + audioBuffer.duration;
}