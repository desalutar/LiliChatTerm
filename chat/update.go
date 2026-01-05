package chat

import (
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

		exists := false
		for _, msg := range c.Messages {
			if msg.ID == incoming.ID {
				exists = true
				break
			}
		}

		if !exists {
			c.Messages = append(c.Messages, Message(incoming))
			c.limitMessages()
		}
	}

	return c, cmd
}
