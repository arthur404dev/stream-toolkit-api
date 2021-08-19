package main

import (
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/arthur404dev/404-api/restream"
	"github.com/arthur404dev/404-api/websocket"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	godotenv.Load()
	env := os.Getenv("ENVIRONMENT")
	if env == "prod" {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&log.JSONFormatter{})
		log.SetReportCaller(true)
	}
	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
	if logLevel == log.TraceLevel {
		log.SetReportCaller(true)
	}
}

func main() {
	runSockets()
	// p := os.Getenv("PORT")
	// setupServer(p)
}

func setupServer(port string) {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodConnect, http.MethodPost},
	}))
	e.Use(loggingMiddleware)

	e.GET("/", statusPage)
	e.POST("/restream/exchange", restream.ExchangeTokens)
	e.GET("/ws", liveChat)
	log.WithField("port", port).Info("api is listening and serving...")
	log.Fatal(e.Start(":" + port))
}

func statusPage(c echo.Context) error {
	return c.String(http.StatusOK, "api is online!")
}

func liveChat(c echo.Context) error {
	ws, err := websocket.Handler(c.Response(), c.Request())
	if err != nil {
		return err
	}
	go websocket.Writer(ws)
	return nil
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

func runSockets() {
	if err := restream.RefreshTokens(); err != nil {
		log.Fatal("tokens couldn't be refreshed")
	}
	u1 := os.Getenv("RESTREAM_CHAT_ENDPOINT")
	u2 := os.Getenv("RESTREAM_UPDATES_ENDPOINT")

	tc := make(chan []byte)
	go websocket.Connect(u1, &tc)
	go websocket.Connect(u2, &tc)

	for {
		select {
		case message := <-tc:
			if message == nil {
				os.Exit(1)
			}
			log.WithField("message", string(message)).Infoln("received message")
		}
	}
}
