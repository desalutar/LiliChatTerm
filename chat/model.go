package chat

import (
	"client/ws"
	"github.com/charmbracelet/bubbles/textinput"
)

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
	MsgChan  		chan incomingMsg
	WsClient 		ws.WsClienter
	Store 			*MessageStore
}

func NewChatScreenModel(userID int64, token string, wsClient *ws.Client) *ChatScreenModel{
    m :=  &ChatScreenModel{
        UserID: userID,
        Token:  token,
        WsClient: wsClient,
        Inputs: ChatInputs{
            ChatAreaInput:    NewChatAreaInput(),
            SearchUserInput:  NewSearchUserAreaInput(),
        },
        State: ChatState{},
        Messages: []Message{},
        MsgChan: make(chan incomingMsg, 100),
    }
	m.Inputs.ChatAreaInput.Focus()
	m.InitWS()
	return m
}