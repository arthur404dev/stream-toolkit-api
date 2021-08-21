package websocket

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/arthur404dev/404-api/restream"
)

func Start(dumpChannel *chan []byte) error {
	logger := log.WithFields(log.Fields{"source": "websocket.Start()", "dumpChannel": dumpChannel})

	if err := restream.RefreshTokens(); err != nil {
		logger.Error(err)
		return err
	}

	u1 := os.Getenv("RESTREAM_CHAT_ENDPOINT")
	u2 := os.Getenv("RESTREAM_UPDATES_ENDPOINT")

	go Connect(u1, dumpChannel)
	go Connect(u2, dumpChannel)

	return nil
}
