package message

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

func Parse(raw []byte) (Message, error) {
	logger := log.WithFields(log.Fields{"source": "message.Parse()", "raw-string": string(raw)})
	logger.Debugln("message parse started")
	parsedMsg := Message{}
	baseMsg := RawMessage{}
	if err := json.Unmarshal(raw, &baseMsg); err != nil {
		logger.Errorln(err)
		return parsedMsg, err
	}
	switch baseMsg.Action {
	case "event":
		{
			if err := json.Unmarshal(raw, &parsedMsg); err != nil {
				logger.Errorln(err)
				return parsedMsg, err
			}
			logger.Infof("%+v\n", parsedMsg)
			parsedMsg.Type = "chat"
			if parsedMsg.Payload.EventPayload.Author.Avatar == "" {
				if parsedMsg.Payload.EventPayload.Author.AvatarUrl != "" {
					parsedMsg.Payload.EventPayload.Author.Avatar = parsedMsg.Payload.EventPayload.Author.AvatarUrl
				}
				if parsedMsg.Payload.EventPayload.Author.Picture != "" {
					parsedMsg.Payload.EventPayload.Author.Avatar = parsedMsg.Payload.EventPayload.Author.Picture
				}
			}
		}
	case "upsert":
		{
			if err := json.Unmarshal(raw, &parsedMsg); err != nil {
				logger.Errorln(err)
				return parsedMsg, err
			}
			parsedMsg.Type = "upsert"
		}
	case "delete":
		{
			if err := json.Unmarshal(raw, &parsedMsg); err != nil {
				logger.Errorln(err)
				return parsedMsg, err
			}
			parsedMsg.Type = "delete"
		}
	case "updateStatuses":
		{
			stats := Stats{}
			if err := json.Unmarshal(raw, &stats); err != nil {
				logger.Errorln(err)
				return parsedMsg, err
			}
			parsedMsg.Action = baseMsg.Action
			parsedMsg.Stats = stats
			parsedMsg.Timestamp = int(time.Now().Unix())
			parsedMsg.Type = "stats"
		}
	}
	logger.Debugln("message parse finished")
	return parsedMsg, nil
}
