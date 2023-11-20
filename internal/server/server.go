package server

import (
	"errors"
	"log"
	"net"
	"net/http"
	"sync/atomic"
)

type (
	Server interface {
		Start() error
		Stop() error
		IsStarted() bool
	}

	server struct {
		server    *http.Server
		isStarted bool
	}

	Config struct {
		Address string `json:"address"`
	}
)

func NewServer(cfg Config, handler http.Handler) Server {
	if cfg.Address == "" {
		cfg.Address = ":8000"
	}
	var (
		connCount  int32
		httpServer = &http.Server{
			Addr:    cfg.Address,
			Handler: handler,
			ConnState: func(conn net.Conn, state http.ConnState) {
				var delta int32
				switch state {
				case http.StateNew:
					delta = 1
				case http.StateClosed:
					delta = -1
				}

				atomic.AddInt32(&connCount, delta)
			},
		}
	)

	return &server{server: httpServer}
}

func (s *server) Start() error {
	s.isStarted = true
	defer func() { s.isStarted = false }()

	log.Println("Start http server on " + s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *server) Stop() error {
	return s.server.Close()
}

func (s *server) IsStarted() bool {
	return s.isStarted
}
