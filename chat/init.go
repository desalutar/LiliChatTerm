package chat

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


func (c *ChatScreenModel) Init() tea.Cmd {
	c.Inputs.ChatAreaInput.Focus()
	c.InitWS()
	return textinput.Blink
}
