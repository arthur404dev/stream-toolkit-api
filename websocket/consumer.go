package websocket

import (
	"os"
	"os/signal"
	"time"

	"github.com/arthur404dev/404-api/restream"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func Connect(url string, tc *chan []byte) {
	logger := log.WithFields(log.Fields{"source": "websocket.Connect()", "url": url, "target-channel": tc})
	logger.Debugln("websocket connection started")
	accessToken, err := restream.GetAccessToken()
	if err != nil {
		logger.Errorln(err)
	}
	logger.Debugln("building options")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url + "/ws?accessToken=" + accessToken
	logger.Infoln("connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		logger.Error(err)
	}
	defer c.Close()
	done := make(chan struct{})
	logger.Debugln("shipped OnMessage goroutine")
	go OnMessage(c, done, *tc)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	logger.Debugln("started watcher")
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			logger.WithField("t", t).Traceln("tick")
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Error(err)
				return
			}
		case <-interrupt:
			logger.Warn("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Error(err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func OnMessage(c *websocket.Conn, done chan struct{}, target chan []byte) {
	logger := log.WithFields(log.Fields{"source": "websocket.OnMessage()", "local-network": c.LocalAddr().String(), "target-channel": target})
	logger.Debugln("websocket onmessage started")
	defer close(done)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			logger.Error(err)
			target <- nil
			return
		}
		logger.WithField("message", string(message)).Debugln("send message to target channel")
		target <- message
	}
}
