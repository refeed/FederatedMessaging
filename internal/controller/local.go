package controller

import (
	"FederatedMessaging/internal/message"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Local struct {
	R              *gin.Engine
	MessageService *message.MessageService
}

func (l *Local) RegisterRoutes() {
	l.R.GET("/api/msg", l.get)
	l.R.POST("/api/msg", l.post)
	l.R.PATCH("/api/msg", l.patch)
	l.R.DELETE("/api/msg", l.delete)
}

// This endpoint might be unused anymore if the websocket endpoint
// is used
func (l *Local) post(c *gin.Context) {
	type NewMsgReq struct {
		Passphrase string
		Body       string
	}
	var msg NewMsgReq
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	l.MessageService.CreateMessage(
		msg.Body,
		"local",
		message.GetUserIdentity(msg.Passphrase),
		strings.HasPrefix(msg.Body, "/cast"))

	c.JSON(http.StatusOK, gin.H{"status": "message received"})
}

func (l *Local) patch(c *gin.Context) {
	type UpdateMsgReq struct {
		Id   uuid.UUID
		Body string
	}
	var updateMsgReq UpdateMsgReq
	if err := c.ShouldBindJSON(&updateMsgReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	l.MessageService.Update(updateMsgReq.Id, updateMsgReq.Body)
	c.JSON(http.StatusOK, gin.H{"status": "message updated"})
}

func (l *Local) delete(c *gin.Context) {
	type DeleteReq struct {
		Id uuid.UUID
	}
	var deleteReq DeleteReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	l.MessageService.Delete(deleteReq.Id)
	c.JSON(http.StatusOK, gin.H{"status": "message deleted"})
}

func (l *Local) get(c *gin.Context) {
	messages, err := l.MessageService.GetBulk(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
