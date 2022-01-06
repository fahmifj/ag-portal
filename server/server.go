package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmifj/ag-portal/config"
	"github.com/fahmifj/ag-portal/logger"
	"github.com/fahmifj/ag-portal/server/handler"
	"github.com/fahmifj/ag-portal/service"
)

var (
	addr = "0.0.0.0:9000"
)

func Run() {
	config.Load()

	service := service.NewServices()

	h := handler.NewRouter(service)

	s := &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		logger.Log.Info("Server starting", "address", s.Addr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Failed to start the server", "reason", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	sig := <-c
	logger.Log.Info("Performing graceful shutdown", "signal", sig)

	// Graceful shutdown the server, waiting for max of 15 seconds until current operations is completed
	tc, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := s.Shutdown(tc)
	if err != nil {
		logger.Log.Error("Graceful shutdown failed", "reason", err)
		os.Exit(1)
	}
}
