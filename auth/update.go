package auth

import (
	"client/chat"
	"client/ws"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	KeyQuit   = tea.KeyMsg{Type: tea.KeyCtrlC}
	KeySwitch = tea.KeyMsg{Type: tea.KeyTab}       
	KeyToggle = tea.KeyMsg{Type: tea.KeyCtrlR}
	KeySubmit = tea.KeyMsg{Type: tea.KeyEnter}
)

const (
    RegistrationInfo = "Registration successful. Please login."
)


func (m *AuthScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case KeyQuit.String():      
			return m, tea.Quit
		case KeySwitch.String():   
			m.SwitchFocus()
			return m, nil
		case KeyToggle.String():   
			m.ToggleMode()
			return m, nil
		case KeySubmit.String():   
			return m, m.Submit()
		}
	case authMsg:
		return m.handleAuthMsg(msg)
	}

	return m.updateInputs(msg)
}


func (m *AuthScreenModel) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    m.Inputs.Username, cmd = m.Inputs.Username.Update(msg)
    m.Inputs.Password, _ = m.Inputs.Password.Update(msg)

    if m.State == RegisterMode {
        m.Inputs.Confirm, _ = m.Inputs.Confirm.Update(msg)
    }

    return m, cmd
}

func (a *AuthScreenModel) handleAuthMsg(msg authMsg) (tea.Model, tea.Cmd) {
    a.Busy = false

    if msg.err != nil {
        a.Error = msg.err.Error()
        return a, nil
    }

    if a.State == RegisterMode {
        a.Error = RegistrationInfo
        a.ToggleMode()
        return a, nil
    }

    userID := msg.resp.UserID
    token := msg.resp.AccessToken

    wsClient, err := ws.New(a.Service.WSURL, token)
    if err != nil {
        a.Error = err.Error()
        return a, nil
    }
    
    return chat.NewChatScreenModel(userID, token, wsClient), nil
}
