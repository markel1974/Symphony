package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

var embeddedFiles embed.FS

func main() {
	subFS, err := fs.Sub(embeddedFiles, "static_content")
	if err != nil {
		log.Fatal("Errore nel creare il sub-filesystem:", err)
	}
	fileServer := http.FileServer(http.FS(subFS))
	http.Handle("/", fileServer)
	log.Println("Server avviato su http://localhost:8080")
	log.Println("Prova ad accedere a http://localhost:8080/index.html o http://localhost:8080/css/style.css")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Errore ListenAndServe: ", err)
	}
}
