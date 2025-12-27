package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *AuthScreenModel) Init() tea.Cmd {
    return textinput.Blink
}