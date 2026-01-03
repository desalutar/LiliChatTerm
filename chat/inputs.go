package chat

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func NewChatAreaInput() textinput.Model {
	chatTi := textinput.New()
	chatTi.Placeholder = "Type a message..."
	chatTi.Focus()
	chatTi.CharLimit = 512
	chatTi.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))   // синий
	chatTi.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))    // белый
	chatTi.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	return chatTi
}

func NewSearchUserAreaInput() textinput.Model {
	searchTi := textinput.New()
	searchTi.Placeholder = "Enter username..."
	searchTi.CharLimit = 32
	searchTi.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	searchTi.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	searchTi.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	return searchTi
}