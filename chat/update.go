package chat

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (c *ChatScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
				return c.handleSearch()
			}
			return c.handleSendMessage()
		}
	case searchResultMsg:
		return c.handleSearchResult(msg)
	}

	if c.State.IsSearchMode {
		c.Inputs.SearchUserInput, cmd = c.Inputs.SearchUserInput.Update(msg)
	} else {
		c.Inputs.ChatAreaInput, cmd = c.Inputs.ChatAreaInput.Update(msg)
	}

	return c, cmd
}
