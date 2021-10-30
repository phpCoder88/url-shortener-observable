package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/phpCoder88/url-shortener/internal/config"
	"github.com/phpCoder88/url-shortener/internal/http/routes"
	"github.com/phpCoder88/url-shortener/internal/ioc"
	"go.uber.org/zap"
)

// Server struct
type Server struct {
	server    http.Server
	logger    *zap.SugaredLogger
	conf      *config.Config
	container *ioc.Container
	errors    chan error
}

func NewServer(logger *zap.SugaredLogger, conf *config.Config, container *ioc.Container, errChan chan error) *Server {
	return &Server{
		server: http.Server{
			Addr:         net.JoinHostPort("", fmt.Sprint(conf.Server.Port)),
			Handler:      routes.Routes(logger, container),
			IdleTimeout:  conf.Server.IdleTimeout,
			ReadTimeout:  conf.Server.ReadTimeout,
			WriteTimeout: conf.Server.WriteTimeout,
		},
		logger:    logger,
		conf:      conf,
		container: container,
		errors:    errChan,
	}
}

func (s *Server) Run() {
	go func() {
		s.logger.Infof("API Server is listening on PORT: %d...", s.conf.Server.Port)
		err := s.server.ListenAndServe()
		if err != nil {
			s.errors <- err
		}
	}()
}

func (s *Server) Stop() error {
	s.logger.Info("Starting to shutdown the API server...")
	ctx, cancel := context.WithTimeout(context.Background(), s.conf.Server.ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
