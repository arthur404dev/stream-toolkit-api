package websocket

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	ips        map[string]bool
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		ips:        make(map[string]bool),
	}
}

func (h *Hub) Run() {
	logger := log.WithFields(log.Fields{"source": "Hub.Run()", "hub": h})
	logger.Infoln()
	urls := strings.Split(os.Getenv("SOCKET_ENDPOINTS"), ",")
	quit := make(chan bool)

	for {
		select {
		case client := <-h.register:
			logger.Debugln("registering client")
			if len(h.clients) == 0 {
				logger.Warnln("no clients found, launching consumers")
				for _, url := range urls {
					go consume(url, &h.broadcast, &quit)
				}
			}
			h.clients[client] = true
			h.ips[client.ip] = true
		case client := <-h.unregister:
			logger.Debugln("unregistering client")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.ips, client.ip)
				close(client.send)
			}
			if len(h.clients) == 0 {
				logger.Warnln("no more clients online, shutting down consumers")
				for range urls {
					quit <- true
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
