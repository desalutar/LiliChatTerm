package chat

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (c *ChatScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, specialCmd := c.handleKeyMsg(msg)
		if specialCmd != nil {
			if casted, ok := m.(*ChatScreenModel); ok {
				c = casted
			}
			return m, specialCmd
		}

		if casted, ok := m.(*ChatScreenModel); ok {
			c = casted
		}
	case searchResultMsg:
		return c.handleSearchResult(msg)
	}
	cmd = c.updateInputs(msg)

	for len(c.MsgChan) > 0 {
		incoming := <-c.MsgChan
		if incoming.SenderID != c.UserID {
			c.Messages = append(c.Messages, Message(incoming))
			c.limitMessages()
		}
	}
	return c, cmd
}

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

	c.Messages = append(c.Messages, msg)
	c.limitMessages()

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
	c.Messages = []Message{}
	c.State.SearchMessage = fmt.Sprintf("Found user: %s (ID: %d)", msg.Username, msg.UserID)
	c.State.IsSearchMode = false
	return c, c.Inputs.ChatAreaInput.Focus()
}
