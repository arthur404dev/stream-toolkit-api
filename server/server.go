package server

import (
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/arthur404dev/api/restream"
	"github.com/arthur404dev/api/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(port string, hub *websocket.Hub) {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodConnect, http.MethodPost},
	}))
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("ACCESS_API_KEY"), nil
		},
	}))
	e.Use(loggingMiddleware)

	e.GET("/", statusPage)
	e.GET("/ws", func(c echo.Context) error { return websocket.ServeWs(hub, c) })
	e.POST("/restream/exchange", restream.ExchangeTokens)

	log.WithField("port", port).Info("api is listening and serving...")
	log.Fatal(e.Start(":" + port))
}

func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		res := next(c)
		log.WithFields(log.Fields{
			"method":     c.Request().Method,
			"path":       c.Path(),
			"status":     c.Response().Status,
			"latency_ns": time.Since(start).Nanoseconds(),
		}).Info("request details")
		return res
	}
}
