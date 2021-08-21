package server

import (
	"net/http"

	"github.com/arthur404dev/404-api/websocket"
	"github.com/labstack/echo/v4"
)

func statusPage(c echo.Context) error {
	return c.String(http.StatusOK, "api is online!")
}

func liveChat(c echo.Context) error {
	ws, err := websocket.Handler(c.Response(), c.Request())
	if err != nil {
		return err
	}
	defer ws.Close()
	tc := make(chan []byte)
	go websocket.Start(&tc)
	websocket.Writer(ws, &tc)
	return nil
}
