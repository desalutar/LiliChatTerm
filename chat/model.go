package chat

import (
	"client/ws"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Message struct {
	SenderID   int64
	ReceiverID int64
	Text       string
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
}

type searchResultMsg struct {
	UserID   int64
	Username string
	Err      error
}

type ChatScreenModel struct {
	Inputs   ChatInputs
	State    ChatState
	UserID   int64
	Token    string
	Messages []Message
	MsgChan  chan incomingMsg
	WsClient ws.WsClienter
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

func searchUserCmd(token, username string) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("http://localhost:9900/api/1/users/%s", username)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return searchResultMsg{Err: err}
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return searchResultMsg{Err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(resp.Body)
			return searchResultMsg{Err: fmt.Errorf("user not found: %s", string(data))}
		}

		var user struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return searchResultMsg{Err: err}
		}

		return searchResultMsg{UserID: user.ID, Username: user.Username, Err: nil}
	}
}