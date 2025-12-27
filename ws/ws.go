package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsClienter interface {
	SendMessage(receiverID int64, text string) error
	ReadMessages(handler func(map[string]interface{}))
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
	req := map[string]interface{}{
		"receiver_id": receiverID,
		"text":        text,
	}
	return c.Conn.WriteJSON(req)
}

func (c *Client) ReadMessages(handler func(map[string]interface{})) {
	go func() {
		for {
			var msg map[string]interface{}
			if err := c.Conn.ReadJSON(&msg); err != nil {
				log.Println("WS read error:", err)
				break
			}
			handler(msg)
		}
	}()
}
