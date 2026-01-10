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
		msgMap, ok := data.(map[string]interface{})
		if !ok {
			log.Println("Received data is not a map:", data)
			return
		}

		switch msgType {
		case "message":
			c.processSingleMessage(msgMap)

		case "history":
			c.handleHistory(msgMap["messages"])
		case "connected":
			log.Println("WS connected:", msgMap)
		case "error":
			errMsg, _ := msgMap["error"].(string)
			log.Println("Server error:", errMsg)

		default:
			log.Println("Unknown WS message type:", msgType)
		}
	})
}

func (c *ChatScreenModel) processSingleMessage(m map[string]interface{}) {
	var msgID string
	if idVal, exists := m["id"]; exists && idVal != nil {
		msgID, _ = idVal.(string)
	} else {
		senderID := utils.SafeInt64(m["sender_id"])
		receiverID := utils.SafeInt64(m["receiver_id"])
		text := utils.SafeString(m["text"])
		msgID = fmt.Sprintf("%v-%v-%v", senderID, receiverID, text)
	}

	if _, exists := c.Store.seenIDs[msgID]; exists {
		return
	}

	senderID := utils.SafeInt64(m["sender_id"])
	receiverID := utils.SafeInt64(m["receiver_id"])
	text := utils.SafeString(m["text"])

	msg := Message{
		ID:         msgID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Text:       text,
	}

    c.Store.AddMessageIfIDNotExist(msg)
    c.Messages = append(c.Messages, msg)
}


func (c *ChatScreenModel) handleHistory(messages interface{}) {
	historySlice, ok := messages.([]interface{})
	if !ok {
		log.Println("History messages is not a slice:", messages)
		return
	}

	for _, m := range historySlice {
		mMap, ok := m.(map[string]interface{})
		if !ok {
			log.Println("History item is not a map:", m)
			continue
		}
		c.processSingleMessage(mMap)
	}
}
