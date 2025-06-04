function boot(config) {
    window.symphonyConfig = config;
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("symphony.wasm"), go.importObject).then((result) => {
        const inst = result.instance;
        go.run(inst);

        setupAudio(inst, 44100);

        setupInputs(inst);

        setupGraphics(inst, 4000);

        setInterval(() => { emulate(); }, 1)
    });
}