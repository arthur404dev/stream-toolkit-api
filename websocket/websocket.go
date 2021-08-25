package websocket

import (
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func ServeWs(hub *Hub, c echo.Context) error {
	logger := log.WithFields(log.Fields{"source": "websocket.ServeWs()", "hub": hub, "req-ip": c.RealIP()})
	logger.Debugln("websocket handler started")
	w := c.Response()
	r := c.Request()
	logger.Debugln("enabling access origin")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	logger.Debugln("upgrading connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	logger.Debugln("connection upgraded")
	logger.WithField("client-ip", c.RealIP()).Infoln("creating client")
	client := &Client{ip: c.RealIP(), hub: hub, conn: conn, send: make(chan []byte, 256)}
	logger.Debugln("registering client")
	client.hub.register <- client
	logger.Infoln("write/read pumps started")
	go client.writePump()
	go client.readPump()
	logger.Debugln("websocket handler finished")
	return nil
}
