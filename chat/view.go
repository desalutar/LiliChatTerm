package chat

import (
	"fmt"
	"strings"
)

func (m ChatScreenModel) View() string {
	var b strings.Builder

	if m.State.IsSearchMode {
		b.WriteString("=== Search User ===\n\n")
		if m.State.SearchMessage != "" {
			b.WriteString(m.State.SearchMessage + "\n\n")
		}
		b.WriteString("Enter username: " + m.Inputs.SearchUserInput.View())
		b.WriteString("\n\nPress Enter to search, Esc to cancel.\n")
		return b.String()
	}

	// Chat view
	b.WriteString(fmt.Sprintf("=== Chat with User %d ===\n", m.State.ReceiverID))
	if m.State.SearchMessage != "" {
		b.WriteString(m.State.SearchMessage + "\n")
	}
	b.WriteString("\n")

	for _, msg := range m.Messages {
		if msg.SenderID == m.UserID {
			b.WriteString(msg.Text + "\n")
		} else {
			b.WriteString(fmt.Sprintf("User %d: %s\n", msg.SenderID, msg.Text))
		}
	}

	b.WriteString("\n" + m.Inputs.ChatAreaInput.View())
	b.WriteString("\n\nPress Enter to send, Ctrl+S to search user, q to quit.\n")
	return b.String()
}
