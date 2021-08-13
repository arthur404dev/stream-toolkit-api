package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arthur404dev/404-api/restream"
	"github.com/joho/godotenv"
)

func statusPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API is online and serving")
}

func setupRoutes(port string) {
	http.HandleFunc("/", statusPage)
	http.HandleFunc("/restream/exchange", restream.ExchangeTokens)
	log.Printf("API Listening on port:%+v and serving...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	godotenv.Load()
	p := os.Getenv("PORT")
	setupRoutes(p)
}
