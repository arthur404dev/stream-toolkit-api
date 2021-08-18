package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arthur404dev/404-api/chat"
	"github.com/arthur404dev/404-api/restream"
	"github.com/joho/godotenv"
)

func statusPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API is online and serving")
}

func liveChat(w http.ResponseWriter, r *http.Request) {
	ws, err := chat.Handler(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	go chat.Writer(ws)
}

func setupRoutes(port string) {
	http.HandleFunc("/", statusPage)
	http.HandleFunc("/restream/exchange", restream.ExchangeTokens)
	http.HandleFunc("/live/chat", liveChat)
	log.Printf("API Listening on port:%+v and serving...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	godotenv.Load()
	u1 := os.Getenv("RESTREAM_CHAT_ENDPOINT")
	u2 := os.Getenv("RESTREAM_UPDATES_ENDPOINT")

	tc := make(chan []byte)
	go chat.Connect(u1, &tc)
	go chat.Connect(u2, &tc)

	for {
		select {
		case message := <-tc:
			if message == nil {
				os.Exit(1)
			}
			log.Printf("Received Message: %+v\n", string(message))
		}
	}
	// p := os.Getenv("PORT")
	// setupRoutes(p)
}
