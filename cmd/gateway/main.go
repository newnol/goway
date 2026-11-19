package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/newnol/goway/internal/httpapi"
)

func main() {
	apiKey := strings.TrimSpace(os.Getenv("GOWAY_API_KEY"))
	if apiKey == "" {
		log.Fatal("GOWAY_API_KEY is required")
	}
	addr := os.Getenv("GOWAY_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	server := &http.Server{
		Addr:              addr,
		Handler:           httpapi.New(apiKey),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("goway listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
