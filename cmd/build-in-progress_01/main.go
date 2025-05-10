package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/config"
	"github.com/kartikey1188/build-in-progress_01/internal/http/routes"
	"github.com/kartikey1188/build-in-progress_01/internal/kafka"
	"github.com/kartikey1188/build-in-progress_01/internal/storage/postgres"
)

func main() {
	// Disabling [GIN-debug] logs for route registration
	gin.SetMode(gin.ReleaseMode)

	// loading config

	cfg := config.MustLoad()

	// setting up database

	storage, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// creating a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// creating a WaitGroup
	var wg sync.WaitGroup

	// starting kafka listeners
	kafka.StartKafkaListeners(ctx, &wg, cfg, storage)

	// setting up router and routes

	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(router, storage)

	//setting up server (with graceful shutdown)

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", cfg.Port),
		Handler: router,
	}

	slog.Info("server started", slog.String("address", fmt.Sprintf("localhost:%s", cfg.Port)))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	// cancelling the context to signal the listener goroutines to stop
	cancel()

	// waiting for all goroutines to finish
	wg.Wait()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("sever shutdown successfully")
}
