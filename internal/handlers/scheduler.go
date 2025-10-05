package handlers

import (
	"message-provider-go/internal/scheduler"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SchedulerHandler holds the scheduler instance
type SchedulerHandler struct {
	scheduler *scheduler.Scheduler
}

// NewSchedulerHandler creates a new scheduler handler
func NewSchedulerHandler(sched *scheduler.Scheduler) *SchedulerHandler {
	return &SchedulerHandler{
		scheduler: sched,
	}
}

// StartJob starts the FetchMessagesJob
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

// StopJob stops the FetchMessagesJob
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
		"message": "FetchMessagesJob stopped successfully",
	})
}

// RestartJob restarts the FetchMessagesJob
func (h *SchedulerHandler) RestartJob(c *gin.Context) {
	// First stop the job
	err := h.scheduler.StopJob("fetch-messages")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to stop job",
			"message": err.Error(),
		})
		return
	}

	// Then start it again
	err = h.scheduler.StartJob("fetch-messages")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to start job",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "FetchMessagesJob restarted successfully",
	})
}

// GetJobStatus returns the status of all jobs
func (h *SchedulerHandler) GetJobStatus(c *gin.Context) {
	jobs := h.scheduler.ListJobs()

	jobDetails := make([]gin.H, 0, len(jobs))
	for _, jobName := range jobs {
		name, interval, err := h.scheduler.GetJobInfo(jobName)
		if err != nil {
			continue
		}
		jobDetails = append(jobDetails, gin.H{
			"name":     name,
			"interval": interval.String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"jobs":   jobDetails,
		"count":  len(jobDetails),
	})
}

// StopScheduler stops the entire scheduler
func (h *SchedulerHandler) StopScheduler(c *gin.Context) {
	h.scheduler.Stop()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Scheduler stopped successfully",
	})
}
