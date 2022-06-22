package message

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type PartnerServer struct {
	Address string
}

func (p *PartnerServer) TellMessage(msg Message) {
	jsonReq, _ := json.Marshal(msg)
	req, _ := http.NewRequest(http.MethodPost, p.getApiAdress(), bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	_, err := client.Do(req)

	if err != nil {
		log.Fatalln("PartnerServer.TellMessage: ", err)
	}
}

func (p *PartnerServer) TellDelete(id uuid.UUID, passphrase string) {
	jsonReq, _ := json.Marshal(map[string]interface{}{"id": id, "passphrase": passphrase})
	req, _ := http.NewRequest(http.MethodDelete, p.getApiAdress(), bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	_, err := client.Do(req)

	if err != nil {
		log.Fatalln("PartnerServer.TellDelete: ", err)
	}
}

func (p *PartnerServer) TellUpdate(msg Message) {
	jsonReq, _ := json.Marshal(msg)
	req, _ := http.NewRequest(http.MethodPatch, p.getApiAdress(), bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	_, err := client.Do(req)

	if err != nil {
		log.Fatalln("PartnerServer.TellUpdate: ", err)
	}
}

func (p *PartnerServer) getApiAdress() string {
	return p.Address + "/api/outside"
}
