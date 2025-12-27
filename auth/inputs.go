package auth

import "github.com/charmbracelet/bubbles/textinput"

func NewUsernameInput() textinput.Model {
    u := textinput.New()
    u.Placeholder = "Login"
    u.CharLimit = 32
    u.Width = 32
    u.Focus()
    return u
}

func NewPasswordInput() textinput.Model {
    p := textinput.New()
    p.Placeholder = "Password"
    p.EchoMode = textinput.EchoPassword
    p.CharLimit = 32
    p.Width = 32
    return p
}

func NewConfirmPasswordInput() textinput.Model {
    c := textinput.New()
    c.Placeholder = "Confirm password"
    c.EchoMode = textinput.EchoPassword
    c.CharLimit = 32
    c.Width = 32
    return c
}
