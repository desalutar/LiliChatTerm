package chat

import (
	"fmt"
	"strings"
)

func (m *ChatScreenModel) View() string {
	var b strings.Builder

	if m.State.IsSearchMode {
		m.renderSearchMode(&b)
	} else {
		m.renderChatMode(&b)
	}

	return b.String()
}

func (m *ChatScreenModel) renderSearchMode(b *strings.Builder) {
	b.WriteString("=== Search User ===\n\n")

	if m.State.SearchMessage != "" {
		b.WriteString(m.State.SearchMessage + "\n\n")
	}

	b.WriteString("Enter username: " + m.Inputs.SearchUserInput.View())
	b.WriteString("\n\nEnter — search, Esc — cancel\n")
}

func (m *ChatScreenModel) renderChatMode(b *strings.Builder) {
	b.WriteString(fmt.Sprintf("=== Chat with %s ===\n\n", m.State.ReceiverName))

	if len(m.Messages) == 0 {
		m.renderEmptyState(b)
	} else {
		m.renderMessages(b)
	}

	b.WriteString("\n" + m.Inputs.ChatAreaInput.View())
	b.WriteString("\n\nEnter — send, Ctrl+S — search user, Ctrl+C — quit\n")
}

func (m *ChatScreenModel) renderEmptyState(b *strings.Builder) {
	if !m.State.HistoryLoaded {
		b.WriteString("Loading history...\n\n")
	} else {
		b.WriteString("No messages yet.\n\n")
	}
}

func (m *ChatScreenModel) renderMessages(b *strings.Builder) {
	for _, msg := range m.Messages {
		if msg.SenderID == m.UserID {
			b.WriteString(strings.Repeat(" ", 80) + msg.Text + "\n")
		} else {
			b.WriteString(fmt.Sprintf("User %d: %s\n", msg.SenderID, msg.Text))
		}
	}
}