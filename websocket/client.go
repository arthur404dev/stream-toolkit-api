package websocket

import (
	"time"

	"github.com/arthur404dev/404-api/message"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	ip   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump() {
	logger := log.WithFields(log.Fields{"source": "client.readPump()", "client-address": c.conn.LocalAddr().String()})
	logger.Debugln("read pump started")
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	logger.Debugln("setting reader options")
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorln(err)
			}
			break
		}
		logger.WithField("message", message).Debugln("received message from client")
	}
}

func (c *Client) writePump() {
	logger := log.WithFields(log.Fields{"source": "client.writePump()", "client-address": c.conn.LocalAddr().String()})
	logger.Debugln("write pump started")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			message, err := message.Parse(msg)
			if err != nil {
				logger.Errorln(err)
				return
			}
			if message.Action != "" {
				logger.WithFields(log.Fields{"message": message}).Warnln("sending message to client")
				c.conn.WriteJSON(message)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
