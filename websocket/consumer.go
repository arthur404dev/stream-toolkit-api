package websocket

import (
	"os"
	"strings"
	"time"

	"github.com/arthur404dev/404-api/restream"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Consumers struct {
	urls []string
	quit chan bool
	hub  *Hub
}

func NewConsumers(h *Hub) *Consumers {
	return &Consumers{
		urls: strings.Split(os.Getenv("SOCKET_ENDPOINTS"), ","),
		quit: make(chan bool),
		hub:  h,
	}
}

func (c *Consumers) Run() {
	for _, url := range c.urls {
		go consume(url, &c.hub.broadcast, &c.quit)
	}
}

func (c *Consumers) Down() {
	for range c.urls {
		c.quit <- true
	}
}

func consume(url string, tc *chan []byte, quit *chan bool) {
	logger := log.WithFields(log.Fields{"source": "websocket.consume()", "url": url})
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
	logger.Debugln("shipped onmessage routine")
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				logger.Errorln(err)
				return
			}
			logger.WithField("message", string(message)).Debugln("message received")
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
