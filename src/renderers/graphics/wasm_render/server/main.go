package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static_content
var embeddedFiles embed.FS

func main() {
	err := fs.WalkDir(embeddedFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Printf("Trovato in embeddedFiles: %s, Directory: %t\n", path, d.IsDir())
		return nil
	})
	if err != nil {
		log.Fatalf("Errore nel listare i file embeddati: %v", err)
	}

	subFS, err := fs.Sub(embeddedFiles, "static_content")
	if err != nil {
		log.Fatalf("Errore nel creare il sub-filesystem da 'static_content': %v", err)
	}

	log.Println("\n--- Contenuto accessibile tramite 'subFS' (radice del server web) ---")
	err = fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Il path "." qui è la radice di subFS, cioè "static_content"
		log.Printf("Trovato in subFS (server root): %s, Directory: %t\n", path, d.IsDir())
		return nil
	})
	if err != nil {
		log.Fatalf("Errore nel listare i file in subFS: %v", err)
	}

	fileServer := http.FileServer(http.FS(subFS))
	http.Handle("/", fileServer)

	log.Println("\nServer avviato su http://localhost:8080")
	log.Println("Se 'wasm.html' è listato sopra in 'subFS (server root)', dovrebbe essere accessibile a http://localhost:8080/wasm.html")
	log.Println("Se 'symphony.wasm' è listato sopra in 'subFS (server root)', dovrebbe essere accessibile a http://localhost:8080/symphony.wasm")

	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Errore ListenAndServe: ", err)
	}
}
