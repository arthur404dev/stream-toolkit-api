package websocket

import (
	"time"

	"github.com/arthur404dev/404-api/restream"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func consume(url string, tc *chan []byte, quit *chan bool) {
	logger := log.WithFields(log.Fields{"source": "websocket.Connect()", "url": url, "target-channel": tc})
	logger.Debugln("websocket connection started")
	accessToken, err := restream.GetAccessToken()
	if err != nil {
		logger.Errorln(err)
	}
	logger.Debugln("building options")

	u := url + "/ws?accessToken=" + accessToken
	logger.Infoln("connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		logger.Error(err)
	}
	defer c.Close()
	done := make(chan struct{})
	logger.Debugln("shipped OnMessage goroutine")
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.WithField("error", err).Error("exiting reader go routine...")
				return
			}
			logger.WithField("message", string(message)).Debugln("send message to target channel")
			*tc <- message
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	logger.Debugln("started watcher")
	for {
		select {
		case <-*quit:
			logger.Warnln("graceful shutdown")
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		case <-done:
			return
		case t := <-ticker.C:
			logger.WithField("t", t).Traceln("tick")
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}
