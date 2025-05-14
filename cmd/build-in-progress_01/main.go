package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/config"
	"github.com/kartikey1188/build-in-progress_01/internal/http/routes"
	"github.com/kartikey1188/build-in-progress_01/internal/pub_sub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage/postgres"
)

func main() {
	// loading config

	cfg := config.MustLoad()

	// setting up database

	storage, err := postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// initializing Pub/Sub client
	pubsubClient, err := pubsub.NewClient(context.Background(), cfg.GCPProjectID)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub client: %v", err)
	}
	defer pubsubClient.Close()

	// initializing Pub/Sub topics

	err = pub_sub.InitTopics(pubsubClient)
	if err != nil {
		log.Fatalf("Failed to initialize Pub/Sub topics: %v", err)
	}
	slog.Info("Pub/Sub topics initialized")

	// initializing Pub/Sub subscriptions

	err = pub_sub.InitSubscriptions(pubsubClient)
	if err != nil {
		log.Fatalf("Failed to initialize subscriptions: %v", err)
	}
	slog.Info("Pub/Sub subscriptions initialized")

	// starting listeners in a separate goroutine

	go func() {
		err = pub_sub.StartListeners(context.Background(), pubsubClient, cfg, storage)
		if err != nil {
			log.Fatalf("Failed to start listeners: %v", err)
		}
		slog.Info("Pub/Sub listeners started")
	}()

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

	routes.SetupRoutes(router, storage, pubsubClient)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("sever shutdown successfully")
}
