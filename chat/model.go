package chat

import (
	"client/ws"
	"github.com/charmbracelet/bubbles/textinput"
)

type Message struct {
    ID         string `json:"id"`           // ← главный ключ уникальности!
    DialogID   string `json:"dialog_id"`
    SenderID   int64  `json:"sender_id"`
    ReceiverID int64  `json:"receiver_id"`
    Text       string `json:"text"`
    CreatedAt  string `json:"created_at"`   // для сортировки
}

type WSMessage struct {
	Type string
	Msg *Message
	List []Message
}

type incomingMsg Message

type ChatInputs struct {
	ChatAreaInput		textinput.Model
	SearchUserInput 	textinput.Model
}

type ChatState struct {
	IsSearchMode  bool
	SearchMessage string
	ReceiverID    int64
	ReceiverName  string
	HistoryLoaded bool
	Error 		  error
}

type searchResultMsg struct {
	UserID   int64
	Username string
	Err      error
}

type ChatScreenModel struct {
	Inputs   		ChatInputs
	State    		ChatState
	UserID   		int64
	Token    		string
	Messages 		[]Message
	WsClient 		ws.WsClienter
}

func NewChatScreenModel(userID int64, token string, wsClient ws.WsClienter) *ChatScreenModel {
	m := &ChatScreenModel{
		UserID:   userID,
		Token:    token,
		WsClient: wsClient,
		Inputs: ChatInputs{
			ChatAreaInput:   NewChatAreaInput(),
			SearchUserInput: NewSearchUserAreaInput(),
		},
		Messages: []Message{},
	}

	m.Inputs.ChatAreaInput.Focus()
	m.InitWS()

	return m
}