package message

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:id`
	UpdatedAt time.Time `json:updated_at`
	Source    string    `json:source`

	// TODO: We can use []byte instead of string for Sender
	// SHA-1 in bytes = 20 bytes
	// SHA-1 in string = 40 bytes
	// 50% reduced memory
	Sender string `json:sender`

	Body           string `json:body`
	IsAnnouncement bool   `json:isAnnouncement`
}

func NewMessageService(localDB map[uuid.UUID]Message, mutex *sync.RWMutex,
	connectedServers []PartnerServer) MessageService {
	return MessageService{
		localDB,
		mutex,
		connectedServers,
	}
}

type MessageService struct {
	localDB          map[uuid.UUID]Message
	mutex            *sync.RWMutex
	connectedServers []PartnerServer
}

func (as *MessageService) CreateMessage(body, source, sender string, isAnnouncement bool) (Message, error) {
	// TODO: Limit the body length
	message := Message{
		// TODO: Add sendAt
		ID:             uuid.New(),
		UpdatedAt:      time.Now().UTC(),
		Body:           body,
		Source:         source,
		Sender:         sender,
		IsAnnouncement: isAnnouncement,
	}

	as.mutex.Lock()
	as.localDB[message.ID] = message
	as.mutex.Unlock()
	if message.IsAnnouncement {
		go as.Broadcast(message)
	}

	return message, nil
}

func (as *MessageService) GetBulk(after *uuid.UUID) ([]Message, error) {
	// FIXME: Use the `after`
	if after == nil {
		after = &uuid.UUID{}
	}

	as.mutex.RLock()
	messages := make([]Message, len(as.localDB))
	index := 0
	for _, message := range as.localDB {
		messages[index] = message
		index++
	}
	as.mutex.RUnlock()

	return messages, nil
}

func (as *MessageService) Receive(msg Message) {
	// TODO: Verify the source
	as.mutex.Lock()
	as.localDB[msg.ID] = msg
	as.mutex.Unlock()
}

func (as *MessageService) Delete(id uuid.UUID, passphrase string) error {
	if GetUserIdentity(passphrase) != as.localDB[id].Sender {
		return fmt.Errorf("invalid passphrase")
	}

	if as.localDB[id].IsAnnouncement {
		go as.BroadcastDelete(id, passphrase)
	}
	as.mutex.Lock()
	delete(as.localDB, id)
	as.mutex.Unlock()

	return nil
}

func (as *MessageService) Update(id uuid.UUID, body string, passphrase string) error {
	if GetUserIdentity(passphrase) != as.localDB[id].Sender {
		return fmt.Errorf("invalid passphrase")
	}
	as.mutex.Lock()
	newMsgObj := as.localDB[id]
	newMsgObj.Body = body
	newMsgObj.UpdatedAt = time.Now().UTC()
	as.localDB[id] = newMsgObj
	as.mutex.Unlock()

	if newMsgObj.IsAnnouncement {
		go as.BroadcastUpdate(newMsgObj)
	}
	return nil
}

// Send message to every peer
func (as *MessageService) Broadcast(msg Message) {
	for _, partnerServer := range as.connectedServers {
		partnerServer.TellMessage(msg)
	}
}

func (as *MessageService) BroadcastDelete(id uuid.UUID, passphrase string) {
	for _, partnerServer := range as.connectedServers {
		partnerServer.TellDelete(id, passphrase)
	}
}

func (as *MessageService) BroadcastUpdate(msg Message) {
	for _, partnerServer := range as.connectedServers {
		partnerServer.TellUpdate(msg)
	}
}
