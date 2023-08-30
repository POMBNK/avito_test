package main

import (
	"context"
	segmentHttp "github.com/POMBNK/avito_test_task/internal/segment/delivery/http"
	segmentRepository "github.com/POMBNK/avito_test_task/internal/segment/repository"
	segmentUseCase "github.com/POMBNK/avito_test_task/internal/segment/useCase"
	"github.com/POMBNK/avito_test_task/internal/server"
	"github.com/POMBNK/avito_test_task/pkg/client/postgresql"
	"github.com/POMBNK/avito_test_task/pkg/config"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/julienschmidt/httprouter"
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

	client, err := postgresql.NewClient(context.Background(), 3, cfg)
	if err != nil {
		logs.Fatalln(err)
	}

	logs.Println("Router initialization...")
	router := httprouter.New()
	logs.Println("Router initialized.")
	segmentStorage := segmentRepository.NewPostgresDB(logs, client)
	segmentService := segmentUseCase.NewService(logs, segmentStorage)
	segmentHandler := segmentHttp.NewHandler(logs, segmentService)
	segmentHandler.Register(router)

	srv := server.New(logs, router, cfg)
	srv.Start()
}
