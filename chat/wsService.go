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
            c.processSingleMessage(msgMap, true)

        case "history":
            historyIface, ok := msgMap["messages"]
            if !ok {
                log.Println("History payload missing 'messages' key")
				return
            }

            historySlice, ok := historyIface.([]interface{})
            if !ok {
				log.Println("History messages is not a slice:", historyIface)
				return
			}

            for _, m := range historySlice {
                mMap, ok := m.(map[string]interface{})
                if !ok {
                    log.Println("History item is not a map:", m)
					continue
                }
                c.processSingleMessage(mMap, false)
            }
        default:
            log.Println("Unknown WS message type:", msgType)
        }
    })
}

func (c *ChatScreenModel) processSingleMessage(m map[string]interface{}, sendToChan bool) {
    var msgID string
    if idVal, exists := m["id"]; exists && idVal != nil {
        msgID, _ = idVal.(string)
    } else {
		msgID = fmt.Sprintf("%v-%v-%v", m["sender_id"], m["receiver_id"], m["created_at"])
    }

    if _, exists := c.Store.seenIDs[msgID]; exists {
        return
    }

    senderID := utils.SafeInt64(m["sender_id"])
	receiverID := utils.SafeInt64(m["receiver_id"])
	text := utils.SafeString(m["text"])

    msg := Message {
        ID: msgID,
        SenderID: senderID,
        ReceiverID: receiverID,
        Text: text,
    }

    c.Store.AddMessageIfIDNotExist(msg)

    if mType, ok := m["type"].(string); ok && mType == "message" {
        c.MsgChan <- incomingMsg{
            ID:         msgID,
            SenderID:   senderID,
            ReceiverID: receiverID,
            Text:       text,
        }
    }
}
