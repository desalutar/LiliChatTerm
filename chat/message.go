package chat

// import (
// 	"fmt"
// )


// type MessageStore struct {
// 	Messages []Message
// 	seenIDs  map[string]struct{}
// }

// func NewMessageStore() *MessageStore {
// 	return &MessageStore{
// 		Messages: []Message{},
// 		seenIDs:  make(map[string]struct{}),
// 	}
// }

// func (s *MessageStore) AddMessageIfIDNotExist(msg Message) bool {
// 	if _, exists := s.seenIDs[msg.ID]; exists {
// 		return false
// 	}
// 	s.Messages = append(s.Messages, msg)
// 	s.seenIDs[msg.ID] = struct{}{}
// 	return true
// }

// func CreageMessage(senderID, receiverID int64, text string, customID interface{}) Message {
// 	var id string
// 	if customID != nil {
// 		id = fmt.Sprintf("%v", customID)
// 	} else {
// 		id = fmt.Sprintf("%d-%d-%d", senderID, receiverID)
// 	}

// 	return Message {
// 		ID:         id,
// 		SenderID:   senderID,
// 		ReceiverID: receiverID,
// 		Text:       text,
// 	}
// }