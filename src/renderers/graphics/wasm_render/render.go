package wasm_render

// GOOS=js GOARCH=wasm go build -o g64.wasm ./src/render/wasmrender

/*

import (
"embed"
"net/http"
)

//go:embed g64.wasm
var wasmFS embed.FS

func main() {
	http.Handle("/", http.FileServer(http.FS(wasmFS)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

*/

/*
import (
"syscall/js" // Per interagire con JavaScript
)

var (
	c64 *C64 // Ipotetico oggetto che rappresenta l'emulatore C64
	displayBuffer []byte // Buffer per l'output testuale (da riempire)
)

func emulateFrame(this js.Value, args []js.Value) interface{} {
	// Esegui un frame di emulazione.
	c64.EmulateFrame() // Ipotetica funzione

	// Aggiorna il display buffer.
	// (Questo è un esempio *molto semplificato*.  Dovresti adattarlo
	// al modo in cui il VIC-II gestisce la memoria video.)
	// for i := 0; i < 1000; i++ { // 40x25 caratteri
	//     displayBuffer[i*2] = c64.Memory.Read(0x0400 + uint16(i)) // Carattere
	//     displayBuffer[i*2+1] = c64.Memory.Read(0xD800 + uint16(i)) // Colore
	// }
	// Invia il buffer
	return nil
}
func GetDisplayBuffer(this js.Value, args []js.Value) interface{} {
	//copy(displayBuffer, ...);
	return nil
}

func keyPressed(this js.Value, args []js.Value) interface{} {
	// Gestisci la pressione di un tasto (esempio).
	keyCode := args[0].Int()
	// ... (converti keyCode in un codice tasto del C64) ...
	// ... (invia il codice tasto al C64 emulato) ...
	return nil
}

//TODO WASM
// https://garciat.com/posts/go-wasm/
// https://github.com/seqsense/webgl-go/tree/master

func main() {
	c := make(chan struct{}, 0)
	// Inizializza l'emulatore.
	c64 = NewC64() // Ipotetica funzione

	// Esponi le funzioni Go a JavaScript.
	js.Global().Set("emulateFrame", js.FuncOf(emulateFrame))
	js.Global().Set("keyPressed", js.FuncOf(keyPressed))
	js.Global().Set("getDisplayBuffer", js.FuncOf(GetDisplayBuffer))

	<-c
}
*/
