package main

import (
	handler "groupie/server"
	"log"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:              ":8081",          //adresse du server (le port choisi est à titre d'exemple) // listes des handlers
		ReadHeaderTimeout: 10 * time.Second, // temps autorisé pour lire les headers
		WriteTimeout:      10 * time.Second, // temps maximum d'écriture de la réponse
		IdleTimeout:       60 * time.Second, // temps maximum entre deux rêquetes
		MaxHeaderBytes:    1 << 20,          // 1 MB // maxinmum de bytes que le serveur va lire
	}
	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/artists", handler.ArtistHandler)
	http.HandleFunc("/filters", handler.FiltersHandler)
	http.HandleFunc("/search", handler.SearchHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
