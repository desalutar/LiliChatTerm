package chat

import "fmt"

func (c *ChatScreenModel) InitWS() {
	go c.WsClient.ReadMessages(func(msgType string, data interface{}) {
		switch msgType {
		case "message":
			m := data.(map[string]interface{})

			// безопасно получаем id
			var msgID string
			if m["id"] != nil {
				msgID, _ = m["id"].(string)
			} else {
				// генерируем временный ID
				msgID = fmt.Sprintf("%v-%v-%v", m["sender_id"], m["receiver_id"], m["created_at"])
			}

			// проверка на дубликат
			for _, existing := range c.Messages {
				if existing.ID == msgID {
					return // уже есть, пропускаем
				}
			}

			c.MsgChan <- incomingMsg{
				ID:         msgID,
				SenderID:   int64(m["sender_id"].(float64)),
				ReceiverID: int64(m["receiver_id"].(float64)),
				Text:       m["text"].(string),
			}

		case "history":
			history := data.(map[string]interface{})["messages"].([]interface{})
			for _, m := range history {
				mMap := m.(map[string]interface{})

				var msgID string
				if mMap["id"] != nil {
					msgID, _ = mMap["id"].(string)
				} else {
					// генерируем временный ID
					msgID = fmt.Sprintf("%v-%v-%v", mMap["sender_id"], mMap["receiver_id"], mMap["created_at"])
				}

				// проверка на дубликат
				duplicate := false
				for _, existing := range c.Messages {
					if existing.ID == msgID {
						duplicate = true
						break
					}
				}
				if duplicate {
					continue
				}

				msg := Message{
					ID:         msgID,
					SenderID:   int64(mMap["sender_id"].(float64)),
					ReceiverID: int64(mMap["receiver_id"].(float64)),
					Text:       mMap["text"].(string),
				}
				c.Messages = append(c.Messages, msg)
			}
		}
	})
}
