package websocket

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	logger := log.WithFields(log.Fields{"source": "websocket.Handler()", "host": r.Host})
	logger.Debugln("websocket handler started")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	logger.Debugln("upgrade connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return conn, err
	}
	logger.Debugln("websocket handler finished")
	return conn, nil

}

func Writer(conn *websocket.Conn) {
	logger := log.WithFields(log.Fields{"source": "websocket.Writer()", "local-network": conn.LocalAddr().String()})
	logger.Debugln("websocket writer started")
	for {
		logger.Debugln("create ticker")
		ticker := time.NewTicker(time.Second)

		for t := range ticker.C {
			logger.WithField("t", t).Traceln("tick")
			jsonString, err := json.Marshal(`{"message":"hey from the api"}`)
			if err != nil {
				logger.Errorln(err)
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonString)); err != nil {
				logger.Errorln(err)
				return
			}
			logger.WithField("message", jsonString).Debugln("wrote message")
		}
	}
}
