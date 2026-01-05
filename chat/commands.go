package chat

import (
	"log"
	tea "github.com/charmbracelet/bubbletea"
)

func (c *ChatScreenModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	if c.State.IsSearchMode {
		c.Inputs.SearchUserInput, cmd = c.Inputs.SearchUserInput.Update(msg)
	} else {
		c.Inputs.ChatAreaInput, cmd = c.Inputs.ChatAreaInput.Update(msg)
	}

	return cmd
}

func (c *ChatScreenModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return c, tea.Quit
	case "esc":
		if c.State.IsSearchMode {
			c.State.IsSearchMode = false
			c.State.SearchMessage = ""
			return c, c.Inputs.ChatAreaInput.Focus()
		}
	case "ctrl+s":
		c.State.IsSearchMode = true
		return c, c.Inputs.SearchUserInput.Focus()
	case "enter":
		return c.handleEnter()
	}
	return c, nil
}

func (c *ChatScreenModel) limitMessages() {
	// TODO: изменить способ хранения сообщений
	const MaxMessages = 200
	if len(c.Messages) > MaxMessages {
		c.Messages = c.Messages[len(c.Messages)-MaxMessages:]
	}
}

func (c *ChatScreenModel) handleEnter() (tea.Model, tea.Cmd) {
	if c.State.IsSearchMode {
		return c.handleSearch()
	}
	return c.handleSendMessage()
}

func (c *ChatScreenModel) handleSearch() (tea.Model, tea.Cmd) {
	username := c.Inputs.SearchUserInput.Value()
	if username == "" {
		return c, nil
	}

	c.State.SearchMessage = "Searching..."
	return c, searchUserCmd(c.Token, username)
}

func (c *ChatScreenModel) handleSendMessage() (tea.Model, tea.Cmd) {
	text := c.Inputs.ChatAreaInput.Value()
	if text == "" {
		return c, nil
	}

	msg := Message{
		SenderID:   c.UserID,
		ReceiverID: c.State.ReceiverID,
		Text:       text,
	}

	if err := c.WsClient.SendMessage(msg.ReceiverID, msg.Text); err != nil {
		c.State.Error = err
	}

	c.Inputs.ChatAreaInput.SetValue("")
	return c, nil
}


func (c *ChatScreenModel) handleSearchResult(msg searchResultMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		c.State.SearchMessage = "Error: " + msg.Err.Error()
		return c, nil
	}

	c.State.ReceiverID = msg.UserID
	c.State.ReceiverName = msg.Username
	c.Messages = []Message{}
	c.State.SearchMessage = ""
	c.State.IsSearchMode = false

	err := c.WsClient.LoadHistory(c.State.ReceiverID)
	if err != nil {
		log.Println("Failed to send load_history:", err)
	}

	return c, c.Inputs.ChatAreaInput.Focus()
}
