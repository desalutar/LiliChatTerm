package chat

func (c *ChatScreenModel) InitWS() {
    go c.WsClient.ReadMessages(func(msg map[string]interface{}) {
        if msg["type"] != "message" {
            return
        }
        senderID := int64(msg["sender_id"].(float64))
        receiverID := int64(msg["receiver_id"].(float64))
        text := msg["text"].(string)
        c.MsgChan <- incomingMsg{SenderID: senderID, ReceiverID: receiverID, Text: text}
    })
}
