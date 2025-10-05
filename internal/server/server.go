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

func New(sched *scheduler.Scheduler) *Server {
	router := gin.New()
	cfg := config.Get()

	router.Use(middleware.GinCORS())
	router.Use(middleware.DBMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")

	msg := api.Group("/message")
	msg.GET("/", handlers.GetMessageHandler)

	schedulerHandler := handlers.NewSchedulerHandler(sched)
	schedulerGroup := api.Group("/scheduler")
	{
		schedulerGroup.GET("/job/start", schedulerHandler.StartJob)
		schedulerGroup.GET("/job/stop", schedulerHandler.StopJob)
	}

	api.GET("/health", healthCheck)

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

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
