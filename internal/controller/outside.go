package controller

import (
	"FederatedMessaging/internal/message"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Outside struct {
	R              *gin.Engine
	MessageService *message.MessageService
}

func (o *Outside) RegisterRoutes() {
	o.R.GET("/api/outside", o.get)
	o.R.POST("/api/outside", o.post)
	o.R.PATCH("/api/outside", o.post)
	o.R.DELETE("/api/outside", o.delete)
}

func (o *Outside) post(c *gin.Context) {
	var message message.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o.MessageService.Receive(message)
	c.JSON(http.StatusOK, gin.H{"status": "message received"})
}

func (o *Outside) delete(c *gin.Context) {
	type DeleteReq struct {
		Id         uuid.UUID
		Passphrase string
	}
	var deleteReq DeleteReq
	if err := c.ShouldBindJSON(&deleteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o.MessageService.Delete(deleteReq.Id, deleteReq.Passphrase)
	c.JSON(http.StatusOK, gin.H{"status": "message deleted"})
}

func (o *Outside) get(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not implemented yet"})
}
