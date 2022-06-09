// Controller to handle local API in websocket (for real-time messaging)
package controller

import (
	"encoding/json"
	"log"
	"strings"

	"FederatedMessaging/internal/message"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type LocalWs struct {
	M              *melody.Melody
	MessageService *message.MessageService
	R              *gin.Engine
}

func (l *LocalWs) RegisterHandlers() {
	l.R.GET("/ws", func(c *gin.Context) {
		l.M.HandleRequest(c.Writer, c.Request)
	})

	l.M.HandleMessage(func(s *melody.Session, b []byte) {
		// Parse the json, add message to db, if ok then BROADCAST
		type NewMsgReq struct {
			Passphrase string
			Body       string
		}
		var msgReq NewMsgReq

		err := json.Unmarshal(b, &msgReq)
		if err != nil {
			log.Println("LocalWs.HandleMessage: ", err)
			return
		}

		message, err := l.MessageService.CreateMessage(
			msgReq.Body,
			"local",
			message.GetUserIdentity(msgReq.Passphrase),
			strings.HasPrefix(msgReq.Body, "/cast"))

		if err != nil {
			log.Println("LocalWs.HandleMessage: ", err)
			return
		}

		msgJson, err := json.Marshal(message)
		if err != nil {
			log.Println("LocalWs.HandleMessage: ", err)
			return
		}

		l.M.Broadcast(msgJson)
	})
}
