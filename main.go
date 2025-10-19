package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bienvenue sur mon site !"))
	})

	log.Println("Serveur démarré sur : http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
