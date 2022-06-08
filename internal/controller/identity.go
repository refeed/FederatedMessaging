package controller

import (
	"FederatedMessaging/internal/message"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Identity struct {
	R *gin.Engine
}

func (i *Identity) RegisterRoutes() {
	i.R.GET("/api/identity", getIdentity)
}

func getIdentity(c *gin.Context) {
	passphrase := c.Query("pass")

	if passphrase == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Please set the 'pass' param"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"identity": message.GetUserIdentity(passphrase)})
}
