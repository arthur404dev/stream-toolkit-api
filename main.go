package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/arthur404dev/404-api/restream"
	"github.com/joho/godotenv"
)

func statusPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API is online and serving")
}

func setupRoutes() {
	http.HandleFunc("/", statusPage)
	http.HandleFunc("/restream/exchange", restream.ExchangeTokens)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	godotenv.Load()
	log.Println("404 api running")
	setupRoutes()
}
