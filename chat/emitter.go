package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	// Handle cors connection
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Handler.upgrader.Upgrade error =%+v\n", err)
		return conn, err
	}

	return conn, nil

}

func Writer(conn *websocket.Conn) {
	for {
		ticker := time.NewTicker(5 * time.Second)

		for t := range ticker.C {
			log.Printf("this is the t=%+v\n", t)

			jsonString, err := json.Marshal(`{"message":"hey from the api"}`)
			if err != nil {
				log.Fatalf("error from json=%+v\n", err)
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				log.Fatalf("error here=%+v\n", err)
				return
			}
		}
		log.Printf("Shipped message")
	}
}
