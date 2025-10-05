package handlers

import (
	"message-provider-go/internal/database"
	"message-provider-go/internal/repository"
	"message-provider-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMessages returns exactly 2 messages from database using transaction managed from handler
func GetMessageHandler(c *gin.Context) {
	ctx := c.Request.Context()

	db := database.Get()

	repo := repository.NewMessageRepository(db)
	service := service.NewMessageService(repo)

	// Transaction'ı service'e geç
	messages, err := service.GetSentMessages(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get messages: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"count":    len(messages),
		"info":     "Retrieved using transaction managed from handler layer",
	})
}
