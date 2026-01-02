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
		if c.State.IsSearchMode {
			c.Inputs.SearchUserInput, cmd = c.Inputs.SearchUserInput.Update(msg)
		} else {
			c.Inputs.ChatAreaInput, cmd = c.Inputs.ChatAreaInput.Update(msg)
		}
	case searchResultMsg:
		return c.handleSearchResult(msg)
	default:
		if c.State.IsSearchMode {
			c.Inputs.SearchUserInput, cmd = c.Inputs.SearchUserInput.Update(msg)
		} else {
			c.Inputs.ChatAreaInput, cmd = c.Inputs.ChatAreaInput.Update(msg)
		}
	}

	for {
		select {
		case incoming := <-c.MsgChan:
			if incoming.SenderID != c.UserID {
				c.Messages = append(c.Messages, Message(incoming))
			}
		default:
			return c, cmd
	}
}

}

func (c *ChatScreenModel) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	c.Inputs.ChatAreaInput, _ = c.Inputs.ChatAreaInput.Update(msg)
	c.Inputs.SearchUserInput, _ = c.Inputs.SearchUserInput.Update(msg)


	for {
		select {
		case incoming := <-c.MsgChan:
			c.Messages = append(c.Messages, Message(incoming))
		default:
			return c, cmd
		}
	}
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
		if c.State.IsSearchMode {
			username := c.Inputs.SearchUserInput.Value()
			if username != "" {
				c.State.SearchMessage = "Searching..."
				return c, searchUserCmd(c.Token, username)
			}
		} else {
			text := c.Inputs.ChatAreaInput.Value()
			if text != "" {
				// Добавляем локально, чтобы сразу отображалось
				c.Messages = append(c.Messages, Message{
					SenderID:   c.UserID,
					ReceiverID: c.State.ReceiverID,
					Text:       text,
				})

				// Отправляем на сервер
				_ = c.WsClient.SendMessage(c.State.ReceiverID, text)

				// Очищаем input
				c.Inputs.ChatAreaInput.SetValue("")
			}
		}
	}
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
