package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/layzy-wolf/timeTrackerTest/docs"
	"github.com/layzy-wolf/timeTrackerTest/internal/database"
	"github.com/layzy-wolf/timeTrackerTest/internal/env"
	"github.com/layzy-wolf/timeTrackerTest/internal/transport/http"
	log "github.com/sirupsen/logrus"
	net "net/http"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	formatter := new(log.TextFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(formatter)
	formatter.FullTimestamp = true

	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log.Infoln("config setup")
	conf := env.Setup()

	if !conf.Debug {
		log.Infoln("debug mode is off")
		log.SetLevel(log.InfoLevel)
	} else {
		log.Infoln("debug mode is on")
		log.SetLevel(log.DebugLevel)
	}

	log.Infoln("database setup connection")
	db := database.Setup(conf)

	r := http.Handler(conf, db)

	srv := &net.Server{
		Addr:    fmt.Sprintf(":%v", conf.Port),
		Handler: r,
	}

	go func(srv *net.Server) {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, net.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}(srv)

	<-ctx.Done()
	stop()

	log.Infoln("Server shutting down gracefully, press Ctrl+C to force")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Infoln("Server exiting")
}
