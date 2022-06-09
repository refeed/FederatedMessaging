package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"

	"FederatedMessaging/internal/configparser"
	"FederatedMessaging/internal/controller"
	"FederatedMessaging/internal/message"
)

var (
	partnersFilepath = flag.String("partners", "", "Path to the partners file containing newline-separated URL of FederatedMessaging servers")
)

func main() {
	messageService := createMessageService()
	r := gin.Default()
	m := melody.New()

	(&controller.LocalWs{m, &messageService, r}).RegisterHandlers()
	(&controller.Outside{r, &messageService}).RegisterRoutes()
	(&controller.Local{r, &messageService}).RegisterRoutes()
	(&controller.Identity{r}).RegisterRoutes()
	registerIndexHandler(r)

	r.Run()
}

func registerIndexHandler(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"title": "Server",
		})
	})
}

// Parse the config file to get the partner servers then instantiate the
// MessageService
func createMessageService() message.MessageService {
	flag.Parse()
	var partnerServers []message.PartnerServer

	if *partnersFilepath != "" {
		partners, err := configparser.ReadLines(*partnersFilepath)
		if err != nil {
			log.Fatalln(err.Error())
		}

		for _, partner := range partners {
			partnerServers = append(
				partnerServers, message.PartnerServer{partner})
			log.Println("Added partner server: ", partner)
		}
	}

	messageService := message.NewMessageService(
		make(map[uuid.UUID]message.Message),
		&sync.RWMutex{},
		partnerServers)

	return messageService
}
