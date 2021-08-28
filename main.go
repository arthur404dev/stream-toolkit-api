package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/arthur404dev/404-api/server"
	"github.com/arthur404dev/404-api/websocket"
	"github.com/joho/godotenv"
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
	hub := websocket.NewHub()
	go hub.Run()
	p := os.Getenv("PORT")
	server.Start(p, hub)
}
