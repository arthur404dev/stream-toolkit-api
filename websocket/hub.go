package websocket

import (
	"time"

	"github.com/arthur404dev/404-api/restream"
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
	logger.Debugln("hub service started")
	consumers := NewConsumers(h)
	ttr := 30 * time.Minute
	if err := restream.RefreshTokens(ttr); err != nil {
		logger.Error(err)
		return
	}
	logger.Debugln("hub watcher started")
	ticker := time.NewTicker(ttr)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			logger.WithField("tick", time.Now().Unix()).Debugln("refreshing tokens")
			if err := restream.RefreshTokens(ttr); err != nil {
				logger.Error(err)
				return
			}
			if len(h.clients) > 0 {
				logger.Debugln("refreshing consumers")
				consumers.Down()
				consumers.Run()
			}
		case client := <-h.register:
			logger.WithField("client-ip", client.ip).Debugln("registering client")
			if len(h.clients) == 0 {
				logger.Warnln("no clients found, launching consumers")
				consumers.Run()
			}
			h.clients[client] = true
			h.ips[client.ip] = true
			logger.WithField("client-ip", client.ip).Debugln("client registered")
		case client := <-h.unregister:
			logger.Debugln("unregistering client")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.ips, client.ip)
				close(client.send)
			}
			if len(h.clients) == 0 {
				logger.Warnln("no more clients online, shutting down consumers")
				consumers.Down()
			}
		case message := <-h.broadcast:
			logger.WithField("message", string(message)).Debugln("broadcasting message")
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
