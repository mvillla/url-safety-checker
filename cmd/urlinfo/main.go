package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mvillla/url-safety-checker/internal/httpapi"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	handler := httpapi.NewHandler().Routes()

	log.Printf("starting url safety checker on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
