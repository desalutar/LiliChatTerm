package ws

import (
	"client/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsClienter interface {
	SendMessage(receiverID int64, text string) error
	ReadMessages(handler func(msgType string, data interface{}))
	LoadHistory(receiverID int64) error
}

type Client struct {
	Conn *websocket.Conn
}

func New(url string, token string) (*Client, error) {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)

	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, err
	}

	return &Client{Conn: conn}, nil
}

func (c *Client) SendMessage(receiverID int64, text string) error {
	msgID := fmt.Sprintf("%v-%v-%v", receiverID, text, fmt.Sprint(utils.NowUnixMilli()))
	req := map[string]interface{}{
		"type":        "message",
		"id": 			msgID,
		"receiver_id": receiverID,
		"text":        text,
	}
	return c.Conn.WriteJSON(req)
}


func (c *Client) ReadMessages(handler func(msgType string, data interface{})) {
	go func ()  {
		for {
			var msg map[string]interface{}
			if err := c.Conn.ReadJSON(&msg); err != nil {
				log.Println("WS read error:", err)
				break
			}

			msgType, ok := msg["type"].(string)
			if !ok {
				log.Println("received WS message without type, skipping:", msg)
				continue
			}

			if msg["id"] == nil {
				senderID := utils.SafeInt64(msg["sender_id"])
				receiverID := utils.SafeInt64(msg["receiver_id"])
				createdAt := utils.SafeInt64(msg["created_at"])
				msg["id"] = fmt.Sprintf("%v-%v-%v", senderID, receiverID, createdAt)
			}

			handler(msgType, msg)
		}
	}()
}

func (c *Client) LoadHistory(receiverID int64) error {
	req := map[string]interface{}{
		"type":        "load_history",
		"receiver_id": receiverID,
	}
	return c.Conn.WriteJSON(req)
}
