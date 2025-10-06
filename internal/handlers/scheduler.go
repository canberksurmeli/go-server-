package handlers

import (
	"fmt"
	"message-provider-go/internal/scheduler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchedulerHandler struct {
	scheduler *scheduler.Scheduler
}

func NewSchedulerHandler(sched *scheduler.Scheduler) *SchedulerHandler {
	return &SchedulerHandler{
		scheduler: sched,
	}
}

func (h *SchedulerHandler) StartJob(c *gin.Context) {
	err := h.scheduler.StartJob("fetch-messages")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to start job",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "FetchMessagesJob started successfully",
	})
}

func (h *SchedulerHandler) StopJob(c *gin.Context) {
	err := h.scheduler.StopJob("fetch-messages")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to stop job",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "FetchMessagesJob stop signal sent successfully",
	})
}

func (h *SchedulerHandler) JobStatus(c *gin.Context) {
	isRunning, err := h.scheduler.IsJobRunning("fetch-messages")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to get job status",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"job_name": "fetch-messages",
		"running":  isRunning,
		"message":  fmt.Sprintf("FetchMessagesJob is running: %v", isRunning),
	})
}
