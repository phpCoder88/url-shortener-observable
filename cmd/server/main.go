// Package main URL shortener API.
//
// Open API for URL shortener service
//
// Terms Of Service:
//
//     Schemes: http
//     Host: localhost:8000
//     BasePath: /api
//     Version: 1.0.0
//     License: MIT https://opensource.org/licenses/MIT
//     Contact: Pavel Bobylev<p_bobylev@bk.ru> https://github.com/phpCoder88
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/phpCoder88/url-shortener/internal/config"
	"github.com/phpCoder88/url-shortener/internal/ioc"
	"github.com/phpCoder88/url-shortener/internal/server"
	"github.com/phpCoder88/url-shortener/internal/storages/postgres"
	"github.com/phpCoder88/url-shortener/internal/version"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	logger = logger.With(
		zap.String("Version", version.Version),
		zap.String("BuildDate", version.BuildDate),
		zap.String("BuildCommit", version.BuildCommit),
	)

	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()
	slogger := logger.Sugar()

	slogger.Info("Starting the application...")
	slogger.Info("Reading configuration and initializing resources...")
	conf, err := config.GetConfig()
	if err != nil {
		slogger.Error(err)
		return
	}

	db, err := postgres.NewPgConnection(conf.DB)
	if err != nil {
		slogger.Fatal("Can't connect to the database.", "err", err)
	}

	slogger.Info("Configuring the application units...")
	container := ioc.NewContainer(db, conf.DB.QueryTimeout)
	errChan := make(chan error, 1)

	slogger.Info("Starting the servers...")
	apiServer := server.NewServer(slogger, conf, container, errChan)
	apiServer.Run()

	metricServer := server.NewMetricServer(slogger, 9000, conf.Server.ShutdownTimeout, errChan)
	metricServer.Run()

	// Graceful shutdown
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case x := <-osSignalChan:
		slogger.Infow("Received a signal.", "signal", x.String())
	case err = <-errChan:
		slogger.Errorw("Received an error from the business logic server.", "err", err)
	}

	slogger.Info("Stopping the servers...")
	err = apiServer.Stop()
	if err != nil {
		slogger.Error("Got an error while stopping the API server.", "err", err)
	}

	err = metricServer.Stop()
	if err != nil {
		slogger.Error("Got an error while stopping the metrics server.", "err", err)
	}

	slogger.Info("The app is calling the last defers and will be stopped.")
}
