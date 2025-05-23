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
		log.Printf("found in embeddedFiles: %s, Directory: %t\n", path, d.IsDir())
		return nil
	})
	if err != nil {
		log.Fatalf("error while listing embedded files: %v", err)
	}

	subFS, err := fs.Sub(embeddedFiles, "static_content")
	if err != nil {
		log.Fatalf("error while creating 'static_content': %v", err)
	}

	err = fs.WalkDir(subFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		log.Printf("found (server root): %s, directory: %t\n", path, d.IsDir())
		return nil
	})
	if err != nil {
		log.Fatalf("error whili listing files: %v", err)
	}

	fileServer := http.FileServer(http.FS(subFS))
	http.Handle("/", fileServer)

	addr := ":8080"
	log.Printf("\nstarting server on %s\n", addr)
	if err = http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
