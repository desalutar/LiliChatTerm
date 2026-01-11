package chat

import (
	"client/utils"
	"fmt"
	"log"
)

func (c *ChatScreenModel) InitWS() {
	if c.Store == nil {
		c.Store = NewMessageStore()
	}

	go c.WsClient.ReadMessages(func(msgType string, data interface{}) {
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return
		}

		switch msgType {
		case "message":
			c.processSingleMessage(dataMap)
		case "history":
			if raw, ok := dataMap["messages"]; ok {
				c.handleHistory(raw)
			} else {
				log.Println("[WS] History key 'messages' missing:", dataMap)
			}
		case "connected":
			log.Println("[WS] Connected:", dataMap)
		case "error":
			errMsg, _ := dataMap["error"].(string)
			log.Println("[WS] Server error:", errMsg)
		default:
			log.Println("[WS] Unknown WS message type:", msgType)
		}
	})
}

func (c *ChatScreenModel) convertMapToMessage(m map[string]interface{}) Message {
	return Message{
		ID:         fmt.Sprintf("%v", m["id"]),
		SenderID:   utils.SafeInt64(m["sender_id"]),
		ReceiverID: utils.SafeInt64(m["receiver_id"]),
		Text:       utils.SafeString(m["text"]),
	}
}

func (c *ChatScreenModel) processSingleMessage(m map[string]interface{}) {
	msg := c.convertMapToMessage(m)

	if msg.ID == "" {
		msg.ID = fmt.Sprintf("%v-%v-%v", msg.SenderID, msg.ReceiverID, msg.Text)
	}

	if _, exists := c.Store.seenIDs[msg.ID]; exists {
		return
	}

	c.Store.AddMessageIfIDNotExist(msg)
	c.Messages = append(c.Messages, msg)
}

func (c *ChatScreenModel) handleHistory(messages interface{}) {
	historySlice, ok := messages.([]interface{})
	if !ok {
		log.Println("[WS] History messages is not a slice:", messages)
		return
	}

	for _, m := range historySlice {
		mMap, ok := m.(map[string]interface{})
		if !ok {
			log.Println("[WS] History item is not a map:", m)
			continue
		}
		c.processSingleMessage(mMap)
	}

	c.State.HistoryLoaded = true
}
