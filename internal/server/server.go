package server

import (
	"fmt"
	"github.com/POMBNK/avito_test_task/pkg/config"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	logs   *logger.Logger
	router *httprouter.Router
	cfg    *config.Config
}

func New(logs *logger.Logger, router *httprouter.Router, cfg *config.Config) *Server {
	return &Server{
		logs:   logs,
		router: router,
		cfg:    cfg,
	}
}

func (s *Server) Start() {
	var listener net.Listener
	var listenErr error
	listener, listenErr = net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Server.Port))
	if listenErr != nil {
		s.logs.Fatal(listenErr)
	}

	server := http.Server{
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.logs.Println("Server error:", err)
		}
	}()

	s.logs.Println("Server started")

	<-interrupt
	s.logs.Println("Shutting down server...")

	shutdownErr := server.Shutdown(nil)
	if shutdownErr != nil {
		s.logs.Println("Failed to gracefully shutdown server:", shutdownErr)
	} else {
		s.logs.Println("Server shutdown")
	}
}
