package oto_render

/*
// ContinuousReader is a type for continuous audio data streaming and playback management.
// It uses a mutex for thread-safe operations on audio chunks and handles data writing via a custom write function.
// The type integrates with an audio player to support playback functionality.
type ContinuousReader struct {
	lock             sync.Mutex
	lastChunk        []uint32
	lastChunkSamples int
	player           oto.Player
	bytes            int
	writeFn          writeFn
}

// NewContinuousReader initializes a new ContinuousReader for streaming audio based on the specified sample rate.
// Returns a pointer to the ContinuousReader instance and an error if the initialization fails.
func NewContinuousReader() *ContinuousReader {
	r := &ContinuousReader{}
	return r
}

// Setup initializes the ContinuousReader with a given sample rate, channel count, and audio format. Returns an error if failed.
func (r *ContinuousReader) Setup(sampleRate int, channelCount int, fo string) error {
	format, ok := _formats[fo]
	if !ok {
		return fmt.Errorf("audio format not found")
	}
	ctx, ready, err := oto.NewContext(sampleRate, channelCount, format.Format)
	if err != nil {
		return fmt.Errorf("failed to create oto context: %w", err)
	}
	<-ready
	r.bytes = format.Bytes
	r.writeFn = format.Func
	r.player = ctx.NewPlayer(r)
	r.player.SetVolume(1.0)
	return nil
}

// Play starts or resumes audio playback using the underlying oto.Player.
func (r *ContinuousReader) Play() {
	r.player.Play()
}

// Err returns the current error state of the audio player associated with the ContinuousReader.
func (r *ContinuousReader) Err() error {
	return r.player.Err()
}

// AddChunk appends a new chunk of audio data to the buffer and updates the sample count for playback synchronization.
func (r *ContinuousReader) AddChunk(chunk []uint32, samples int) {
	chunkLen := len(chunk)
	if chunkLen == 0 {
		return
	}
	lastChunkSamples := samples / 2
	r.lock.Lock()
	defer r.lock.Unlock()
	if chunkLen != len(r.lastChunk) {
		r.player.(oto.BufferSizeSetter).SetBufferSize(((chunkLen / 2) * r.bytes) + 1)
		r.lastChunk = make([]uint32, chunkLen)
	}
	copy(r.lastChunk, chunk)
	r.lastChunkSamples = lastChunkSamples
}

// Read reads data from the last available chunk into the provided buffer and returns the number of bytes written.
func (r *ContinuousReader) Read(buf []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.lastChunk == nil || r.lastChunkSamples == 0 {
		return 0, nil
	}
	written := 0
	target := r.lastChunkSamples
	if max := target * r.bytes; max >= len(buf) {
		target = len(buf) / r.bytes
	}
	for x := 0; x < target; x++ {
		written += r.writeFn(buf, r.lastChunk[x], written)
	}
	return written, nil
}


*/
