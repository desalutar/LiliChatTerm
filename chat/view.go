package chat

import (
	"fmt"
	"strings"
)

func (m *ChatScreenModel) View() string {
	var b strings.Builder

	if m.State.IsSearchMode {
		m.searchUserView(&b)
	} else {
		m.chatView(&b)

		for _, msg := range m.Messages {
			if msg.SenderID == m.UserID {
				b.WriteString(strings.Repeat(" ", 80)+ msg.Text + "\n")
			} else {
				// TODO: при получении сообщения показать от кого пришло а не UserID
				b.WriteString(fmt.Sprintf("User %d: %s\n", msg.SenderID, msg.Text))
			}
		}

		b.WriteString("\n" + m.Inputs.ChatAreaInput.View())
		b.WriteString("\n\nPress Enter to send, Ctrl+S to search user, q to quit.\n")
	}

	return b.String()
}

func (m *ChatScreenModel) searchUserView(sb *strings.Builder) string {
	sb.WriteString("=== Search User ===\n\n")
	if m.State.SearchMessage != "" {
		sb.WriteString(m.State.SearchMessage + "\n\n")
	}
	sb.WriteString("Enter username: " + m.Inputs.SearchUserInput.View())
	sb.WriteString("\n\nPress Enter to search, Esc to cancel.\n")
	return sb.String()
}

func (m *ChatScreenModel) chatView(sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("=== Chat with User %s ===\n\n", m.State.ReceiverName))

	msgs := m.Store.Messages
	if len(msgs) == 0 {
		if !m.State.HistoryLoaded {
			sb.WriteString("Loading history...\n\n")
		} else {
			sb.WriteString("No messages yet.\n\n")
		}
	}

	for _, msg := range msgs {
		if msg.SenderID == m.UserID {
			sb.WriteString(strings.Repeat(" ", 80) + msg.Text + "\n")
		} else {
			sb.WriteString(fmt.Sprintf("User %d: %s\n", msg.SenderID, msg.Text))
		}
	}
}
