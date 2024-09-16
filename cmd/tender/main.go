package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Te8va/Tender/internal/tender/config"
	"github.com/Te8va/Tender/internal/tender/handler"
	"github.com/Te8va/Tender/internal/tender/middleware"
	"github.com/Te8va/Tender/internal/tender/repository"
	"github.com/Te8va/Tender/internal/tender/service"
	"github.com/Te8va/Tender/pkg/logger"
	"github.com/caarlos0/env/v6"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Logger().Fatalln("Failed to parse env: %v", err)
	}

	m, err := migrate.New("file://migrations", cfg.PostgresConn)
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}

	err = repository.ApplyMigrations(m)
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}

	logger.Logger().Infoln("Migrations applied successfully")

	pool, err := repository.GetPgxPool(cfg.PostgresConn)
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}

	logger.Logger().Infoln("Postgres connection pool created")

	pg := repository.NewPostgres(pool)

	var wg sync.WaitGroup

	pingRep := repository.NewPingProvider(pg)
	pingService := service.NewPingProvider(pingRep)
	pingHandler := handler.NewPingProvider(pingService)

	tenderRep := repository.NewTenderService(pool)
	tenderService := service.NewTender(tenderRep)
	tenderHandler := handler.NewTenderHandler(tenderService)

	_, cancelDeleteCtx := context.WithCancel(context.Background())

	mux := http.NewServeMux()

	mux.Handle("/api/ping", middleware.Log(http.HandlerFunc(pingHandler.PingHandler)))

	mux.Handle("GET /api/tenders", middleware.Log(http.HandlerFunc(tenderHandler.ListTenderHandler)))
	mux.Handle("POST /api/tender/new", middleware.Log(http.HandlerFunc(tenderHandler.CreateTenderHandler)))
	mux.Handle("GET /api/tenders/my", middleware.Log(http.HandlerFunc(tenderHandler.GetUserTendersHandler)))
	mux.Handle("PATCH /api/tenders/{tenderId}/edit", middleware.Log(http.HandlerFunc(tenderHandler.UpdatePartTenderHandler)))
	mux.Handle("GET /api/tenders/{tenderId}/status", middleware.Log(http.HandlerFunc(tenderHandler.GetTenderStatusHandler)))
	mux.Handle("PUT /api/tenders/{tenderId}/status", middleware.Log(http.HandlerFunc(tenderHandler.UpdateTenderStatusHandler)))
	mux.Handle("PUT /api/tenders/{tenderId}/rollback/{version}", middleware.Log(http.HandlerFunc(tenderHandler.RollbackTenderHandler)))

	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", cfg.ServiceHost, cfg.ServicePort),
		ErrorLog: log.New(logger.Logger(), "", 0),
		Handler:  mux,
	}

	go func() {
		logger.Logger().Infoln("Server started, listening on port", cfg.ServicePort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger().Fatalln("ListenAndServe failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logger.Logger().Infoln("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Logger().Fatalln("Server was forced to shutdown:", zap.Error(err))
	}

	waitGroupChan := make(chan struct{})
	go func() {
		wg.Wait()
		waitGroupChan <- struct{}{}
	}()

	select {
	case <-waitGroupChan:
		logger.Logger().Infoln("All delete goroutines successfully finished")
	case <-time.After(time.Second * 3):
		cancelDeleteCtx()
		logger.Logger().Infoln("Some of delete goroutines have not completed their job due to shutdown timeout")
	}

	logger.Logger().Infoln("Server was shut down")
}
