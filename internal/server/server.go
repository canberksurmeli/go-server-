package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"message-provider-go/internal/config"
	"message-provider-go/internal/handlers"
	"message-provider-go/internal/middleware"
	"message-provider-go/internal/scheduler"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	server *http.Server
}

// New creates a new HTTP server with configured routes and middleware
func New(sched *scheduler.Scheduler) *Server {
	router := gin.New()
	cfg := config.Get()

	// Add middleware
	router.Use(middleware.GinCORS())
	router.Use(middleware.DBMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api/v1")

	// User/Message routes - handler'ı parametre olarak geçiriyoruz
	msg := api.Group("/message")
	msg.GET("/", handlers.GetMessageHandler)

	// Scheduler routes
	schedulerHandler := handlers.NewSchedulerHandler(sched)
	schedulerGroup := api.Group("/scheduler")
	{
		schedulerGroup.POST("/job/start", schedulerHandler.StartJob)
		schedulerGroup.POST("/job/stop", schedulerHandler.StopJob)
		schedulerGroup.POST("/job/restart", schedulerHandler.RestartJob)
		schedulerGroup.GET("/job/status", schedulerHandler.GetJobStatus)
		schedulerGroup.POST("/stop", schedulerHandler.StopScheduler)
	}

	// Health check
	api.GET("/health", healthCheck)

	// Static file serving (optional)
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Http.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		engine: router,
		server: server,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// healthCheck is a simple health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
