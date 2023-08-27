package main

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/segment/db"
	"github.com/POMBNK/avito_test_task/pkg/client/postgresql"
	"github.com/POMBNK/avito_test_task/pkg/config"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
)

// @title Avito Segment Service API
// @version 1.0
// @description API Server for Avito Segment Service
// @contact.name Uchaev Roman
// @contact.url https://github.com/POMBNK
// @contact.email uchaevroman11@gmail.com
// @host 127.0.0.1:8080
// @BasePath /

func main() {
	logs := logger.GetLogger()
	logs.Println("Logger initialized.")

	logs.Println("Config initialization...")
	cfg := config.GetCfg()
	logs.Println("Config initialized.")

	client, err := postgresql.NewClient(context.Background(), cfg)
	if err != nil {
		logs.Fatalln(err)
	}
	//TODO: change to fiber or echo
	logs.Println("Router initialization...")
	router := httprouter.New()
	logs.Println("Router initialized.")
	segmentStorage := db.NewPostgresDB(logs, client)
	segmentService := segment.NewService(logs, segmentStorage)
	segmentHandler := segment.NewHandler(logs, segmentService)
	segmentHandler.Register(router)

	start(logs, router, cfg)
}

func start(logs *logger.Logger, router *httprouter.Router, cfg *config.Config) {
	var listener net.Listener
	var listenErr error
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", "127.0.0.1", "8080"))
	if listenErr != nil {
		logs.Fatal(listenErr)
	}

	server := http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		logs.Fatalf("Server error:%s", err)
	}
}
