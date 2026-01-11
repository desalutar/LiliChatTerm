package chat

import (
	"client/utils"
	"log"
	"sort"
	"time"
)

func (c *ChatScreenModel) InitWS() {
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
			}
		case "connected":
			// silent
		case "error":
			errMsg, _ := dataMap["error"].(string)
			log.Println("[WS error]", errMsg)
		}
	})
}

func (c *ChatScreenModel) convertMapToMessage(m map[string]interface{}) Message {
	return Message{
		ID:         utils.SafeString(m["id"]),
		DialogID:   utils.SafeString(m["dialog_id"]),
		SenderID:   utils.SafeInt64(m["sender_id"]),
		ReceiverID: utils.SafeInt64(m["receiver_id"]),
		Text:       utils.SafeString(m["text"]),
		CreatedAt:  utils.SafeString(m["created_at"]),
	}
}

func (c *ChatScreenModel) processSingleMessage(data map[string]interface{}) {
	msg := c.convertMapToMessage(data)

	if msg.ID == "" {
		return // без id от сервера не принимаем
	}

	// Простая проверка на дубликат по id
	for _, existing := range c.Messages {
		if existing.ID == msg.ID {
			return
		}
	}

	c.Messages = append(c.Messages, msg)
	c.sortMessages()
}

func (c *ChatScreenModel) handleHistory(raw interface{}) {
	slice, ok := raw.([]interface{})
	if !ok {
		return
	}

	for _, item := range slice {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		c.processSingleMessage(m)
	}

	c.State.HistoryLoaded = true
	c.sortMessages()
}

func (c *ChatScreenModel) sortMessages() {
	sort.Slice(c.Messages, func(i, j int) bool {
		ti := parseTime(c.Messages[i].CreatedAt)
		tj := parseTime(c.Messages[j].CreatedAt)
		return ti.Before(tj)
	})
}

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, _ = time.Parse("2006-01-02 15:04:05", s)
	}
	if err != nil {
		return time.Time{} // fallback
	}
	return t
}