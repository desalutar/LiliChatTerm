package chat

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)


func (c *ChatScreenModel) Init() tea.Cmd {
	return textinput.Blink
}
