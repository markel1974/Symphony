const _graphicsVsSource = `
      attribute vec4 aVertexPosition;
      attribute vec2 aTextureCoord;
      varying highp vec2 vTextureCoord;
      void main(void) {
        gl_Position = aVertexPosition; // Il quad coprirà l'intero clip space
        vTextureCoord = aTextureCoord;
      }
    `;

const _graphicsFsSource = `
      varying highp vec2 vTextureCoord;
      uniform sampler2D uSampler;
      void main(void) {
        gl_FragColor = texture2D(uSampler, vTextureCoord);
      }
    `;


function initShaderProgram(gl, vsSource, fsSource) {
    const vertexShader = loadShader(gl, gl.VERTEX_SHADER, vsSource);
    const fragmentShader = loadShader(gl, gl.FRAGMENT_SHADER, fsSource);
    if (!vertexShader || !fragmentShader) return null;

    const shaderProgram = gl.createProgram();
    gl.attachShader(shaderProgram, vertexShader);
    gl.attachShader(shaderProgram, fragmentShader);
    gl.linkProgram(shaderProgram);

    if (!gl.getProgramParameter(shaderProgram, gl.LINK_STATUS)) {
        console.error('can\'t initialize shader: ' + gl.getProgramInfoLog(shaderProgram));
        return null;
    }
    return shaderProgram;
}

function loadShader(gl, type, source) {
    const shader = gl.createShader(type);
    gl.shaderSource(shader, source);
    gl.compileShader(shader);

    if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
        console.error('error while compiling shader (' + (type === gl.VERTEX_SHADER ? "Vertex" : "Fragment") + '): ' + gl.getShaderInfoLog(shader));
        gl.deleteShader(shader);
        return null;
    }
    return shader;
}

function createAndSetupTexture(gl) {
    const texture = gl.createTexture();
    gl.bindTexture(gl.TEXTURE_2D, texture);

    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
    gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
    gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true);

    gl.bindTexture(gl.TEXTURE_2D, null);
    return texture;
}

function initBuffers(gl) {
    const positions = new Float32Array([
        -1.0, -1.0, 1.0, -1.0, -1.0,  1.0, // Triangle 1
        -1.0,  1.0, 1.0, -1.0, 1.0,  1.0,  // Triangle 2
    ]);
    const positionBuffer = gl.createBuffer();
    gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, positions, gl.STATIC_DRAW);
    const textureCoordinates = new Float32Array([
        0.0,  0.0, // Triangle 1 Basso sx texture
        1.0,  0.0, // Triangle 1 Basso dx texture
        0.0,  1.0, // Triangle 1 Alto sx texture
        0.0,  1.0, // Triangle 2 Alto sx texture
        1.0,  0.0, // Triangle 2 Basso dx texture
        1.0,  1.0, // Triangle 2 Alto dx texture
    ]);
    const textureCoordBuffer = gl.createBuffer();
    gl.bindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW);
    return {
        position: positionBuffer,
        textureCoord: textureCoordBuffer,
        vertexCount: 6, // 2 triangles = 6 vertex
    };
}


function setupScene(gl, programInfo, buffers) {
    // Queste chiamate settano lo stato di WebGL per il rendering
    // e non cambiano a meno che tu non cambi il tipo di oggetto, shader, ecc.

    gl.useProgram(programInfo.program); // Specifica lo shader program da usare (solitamente una volta)

    // Attributi dei vertici (posizione)
    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.position);
    gl.vertexAttribPointer(programInfo.attribLocations.vertexPosition, 2, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(programInfo.attribLocations.vertexPosition);

    // Attributi delle coordinate texture
    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.textureCoord);
    gl.vertexAttribPointer(programInfo.attribLocations.textureCoord, 2, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(programInfo.attribLocations.textureCoord);

    // Queste impostazioni sono generali e di solito fisse per il rendering di frame
    gl.clearColor(0.0, 0.0, 0.0, 1.0); // Una volta all'inizio o solo se il colore di clear cambia
    gl.clearDepth(1.0); // Idem
    // gl.enable(gl.DEPTH_TEST); // Solo se stai facendo 3D e non hai un singolo quad 2D
    // gl.depthFunc(gl.LEQUAL); // Idem
}

function setupGraphics(wasmInstance, frameInterval) {
    const canvas = document.getElementById("canvas");
    const graphicsGL = canvas.getContext("webgl");

    if (!graphicsGL) {
        alert("WebGL not supported or disabled");
        return;
    }

    const shaderProgram = initShaderProgram(graphicsGL, _graphicsVsSource, _graphicsFsSource);
    if (!shaderProgram) return;

    const graphicsShaderProgramInfo = {
        program: shaderProgram,
        attribLocations: {
            vertexPosition: graphicsGL.getAttribLocation(shaderProgram, 'aVertexPosition'),
            textureCoord: graphicsGL.getAttribLocation(shaderProgram, 'aTextureCoord'),
        },
        uniformLocations: {
            uSampler: graphicsGL.getUniformLocation(shaderProgram, 'uSampler'),
        },
    };

    const graphicsBuffers = initBuffers(graphicsGL);
    const graphicsTexture = createAndSetupTexture(graphicsGL);
    const graphicsImageWidth = getDisplayWidth();
    const graphicsImageHeight = getDisplayHeight();
    const initialSurfacePtr = getSurfacePointer();
    const initialSurfaceLen = getSurfaceLen();
    let surfaceViewFromWasm = new Uint8Array(wasmInstance.exports.mem.buffer, initialSurfacePtr, initialSurfaceLen);

    setupScene(graphicsGL, graphicsShaderProgramInfo, graphicsBuffers)

    function onVBlank() {
        const currentWasmMemoryBuffer = wasmInstance.exports.mem.buffer;
        if (surfaceViewFromWasm.buffer !== currentWasmMemoryBuffer || surfaceViewFromWasm.byteLength === 0) {
            surfaceViewFromWasm = new Uint8Array(currentWasmMemoryBuffer, initialSurfacePtr, initialSurfaceLen);
        }
        //const jsFrameBuffer = getDisplayBuffer();
        graphicsGL.bindTexture(graphicsGL.TEXTURE_2D, graphicsTexture);
        graphicsGL.texImage2D(graphicsGL.TEXTURE_2D, 0, graphicsGL.RGBA, graphicsImageWidth, graphicsImageHeight, 0, graphicsGL.RGBA, graphicsGL.UNSIGNED_BYTE, surfaceViewFromWasm);
        requestAnimationFrame(() => {
            //drawWebGLScene(graphicsGL, graphicsShaderProgramInfo, graphicsBuffers, graphicsTexture);
            graphicsGL.activeTexture(graphicsGL.TEXTURE0); // Attiva la texture unit 0
            graphicsGL.bindTexture(graphicsGL.TEXTURE_2D, graphicsTexture); // Collega la tua texture a quella unit
            graphicsGL.uniform1i(graphicsShaderProgramInfo.uniformLocations.uSampler, 0); // Comunica allo shader che uSampler è sulla texture unit 0
            graphicsGL.drawArrays(graphicsGL.TRIANGLES, 0, 6); // Assumendo un quad fatto da 6 vertici (2 triangoli)
            //drawFrame(graphicsGL, graphicsShaderProgramInfo, graphicsTexture);
        });
    }

    initRender(frameInterval, onVBlank)
    //canvas.width = imageWidth * 3;
    //canvas.height = imageHeight * 3;
    //gl.viewport(0, 0, gl.canvas.width, gl.canvas.height);
}