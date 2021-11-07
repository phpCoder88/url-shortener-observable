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
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/phpCoder88/url-shortener-observable/internal/config"
	"github.com/phpCoder88/url-shortener-observable/internal/ioc"
	"github.com/phpCoder88/url-shortener-observable/internal/server"
	"github.com/phpCoder88/url-shortener-observable/internal/storages/postgres"
	"github.com/phpCoder88/url-shortener-observable/internal/tracing"
	"github.com/phpCoder88/url-shortener-observable/internal/version"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	tracer, closer, err := tracing.InitJaeger("jaeger:6831", "shortener", logger)
	if err != nil {
		logger.Error(fmt.Sprintf("can't initialize Jaeger: %v", err))
		return
	}

	defer func() {
		err = closer.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	logger = logger.With(
		zap.String("Version", version.Version),
		zap.String("BuildDate", version.BuildDate),
		zap.String("BuildCommit", version.BuildCommit),
	)

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
	container := ioc.NewContainer(db, conf.DB.QueryTimeout, tracer)
	apiServer := server.NewServer(slogger, conf, container, tracer)
	err = apiServer.Run()
	if err != nil {
		slogger.Error("Occurred error during stopping the API server.", "err", err)
	}

	slogger.Info("The app is calling the last defers and will be stopped.")
}
