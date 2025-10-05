package handlers

import (
	"message-provider-go/internal/database"
	"message-provider-go/internal/repository"
	"message-provider-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMessageHandler(c *gin.Context) {
	ctx := c.Request.Context()

	db := database.Get()

	repo := repository.NewMessageRepository(db)
	service := service.NewMessageService(repo)

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
	})
}
