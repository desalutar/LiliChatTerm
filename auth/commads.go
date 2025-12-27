package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
    PasswordErr = "Passwords do not match"
)

type authMsg struct {
	resp *AuthResponse
	err  error
}

func (m *AuthScreenModel) Submit() tea.Cmd {
    if m.Busy {
        return nil
    }
    m.Busy = true

    u := m.Inputs.Username.Value()
    p := m.Inputs.Password.Value()

    if m.State == RegisterMode {
        if p != m.Inputs.Confirm.Value() {
            m.Error = PasswordErr
            m.Busy = false
            return nil
        }
        return func() tea.Msg {
            resp, err := m.Service.Register(u, p)
            return authMsg{resp, err}
        }
    }

    return func() tea.Msg {
        resp, err := m.Service.Login(u, p)
        return authMsg{resp, err}
    }
}

func (m *AuthScreenModel) ToggleMode() {
    if m.State == LoginMode {
        m.State = RegisterMode
    } else {
        m.State = LoginMode
    }
    m.Inputs.FocusIdx = 0
    m.Error = ""
    m.clearInputs()
}

func (m *AuthScreenModel) clearInputs() {
    m.Inputs.Username.SetValue("")
    m.Inputs.Password.SetValue("")
    m.Inputs.Confirm.SetValue("")
}

func (m *AuthScreenModel) SwitchFocus() {
    max := 2
    if m.State == RegisterMode {
        max = 3
    }

    m.Inputs.FocusIdx = (m.Inputs.FocusIdx + 1) % max

    inputs := []*textinput.Model{
        &m.Inputs.Username,
        &m.Inputs.Password,
        &m.Inputs.Confirm,
    }

    for i := 0; i < max; i++ {
        if i == m.Inputs.FocusIdx {
            inputs[i].Focus()
        } else {
            inputs[i].Blur()
        }
    }
}
