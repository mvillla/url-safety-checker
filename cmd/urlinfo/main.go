package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mvillla/url-safety-checker/internal/httpapi"
	"github.com/mvillla/url-safety-checker/internal/lookup"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	malwareURLsFile := os.Getenv("MALWARE_URLS_FILE")
	if malwareURLsFile == "" {
		malwareURLsFile = "data/malware_urls.txt"
	}

	urls, err := lookup.LoadURLsFile(malwareURLsFile)
	if err != nil {
		log.Fatalf("load malware URLs file: %v", err)
	}

	store := lookup.NewMemoryStore(urls)
	lookupService := lookup.NewService(store)

	addr := ":" + port
	handler := httpapi.NewHandler(lookupService).Routes()

	log.Printf("loaded %d malware URLs from %s", len(urls), malwareURLsFile)
	log.Printf("starting url safety checker on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
