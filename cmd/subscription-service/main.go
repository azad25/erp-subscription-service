package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"subscription-service/pkg/config"
	"subscription-service/pkg/database"
	"subscription-service/pkg/redis"
	"subscription-service/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient, err := redis.NewClient(cfg.Redis)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize gRPC server
	// TODO: Implement gRPC server
	// grpcServer := grpc.NewServer(cfg, db, redisClient, logger)

	// Start gRPC server
	// TODO: Uncomment when gRPC is implemented
	/*
		go func() {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
			if err != nil {
				logger.Fatalf("Failed to listen for gRPC: %v", err)
			}

			// Enable reflection for development
			reflection.Register(grpcServer)

			logger.Infof("gRPC server listening on port %d", cfg.GRPC.Port)
			if err := grpcServer.Serve(lis); err != nil {
				logger.Fatalf("Failed to serve gRPC: %v", err)
			}
		}()
	*/

	// Initialize HTTP server
	router := gin.Default()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(utils.CORSMiddleware())
	router.Use(utils.RequestIDMiddleware())
	router.Use(utils.LoggingMiddleware(logger))

	// Initialize handlers
	// TODO: Implement handlers
	// handlers.InitializeRoutes(router, cfg, db, redisClient, logger)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "subscription-service",
			"timestamp": time.Now().UTC(),
		})
	})

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}

	go func() {
		logger.Infof("HTTP server listening on port %d", cfg.HTTP.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down subscription service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP server forced to shutdown: %v", err)
	}

	// Shutdown gRPC server
	// TODO: Uncomment when gRPC is implemented
	// grpcServer.GracefulStop()

	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	// Close Redis connection
	redisClient.Close()

	logger.Info("Subscription service stopped")
}
