package websocket

import (
	"os"

	log "github.com/sirupsen/logrus"

	"net/http"
	"time"

	"github.com/arthur404dev/404-api/message"
	"github.com/gorilla/websocket"
)

const (

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

func Writer(conn *websocket.Conn, source *chan []byte) {
	logger := log.WithFields(log.Fields{"source": "websocket.Writer()", "local-network": conn.LocalAddr().String()})
	logger.Debugln("websocket writer started")

	logger.Debugln("create ticker")
	ticker := time.NewTicker(5 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case rawMsg := <-*source:
			if rawMsg == nil {
				os.Exit(1)
			}
			msg, err := message.Parse(rawMsg)
			if err != nil {
				logger.Errorln(err)
				return
			}
			conn.WriteJSON(msg)

		case t := <-ticker.C:
			logger.Infoln("Heartbeat on %+v", t)
			hb := message.Message{
				Action:    "heartbeat",
				Timestamp: int(time.Now().Unix()),
			}
			if err := conn.WriteJSON(hb); err != nil {
				return
			}
		}
	}
}
