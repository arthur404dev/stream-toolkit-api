package websocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func ServeWs(hub *Hub, c echo.Context) error {
	w := c.Response()
	r := c.Request()
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	client := &Client{ip: c.RealIP(), hub: hub, conn: conn, send: make(chan []byte, 256)}
	log.Warnln("checking ip", client.ip)
	if _, ok := hub.ips[client.ip]; ok {
		log.Warnln("bouncing ip...", client.ip)
		return nil
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	return nil
}
