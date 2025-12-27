package auth

import "strings"

const (
    Login       = "=== LOGIN ===\n"
    Register    = "=== REGISTER ===\n"
    HelpFooter = "\nTAB — switch field\nENTER — submit\nCTRL+R — toggle mode\n"
)

func (m *AuthScreenModel) View() string {
    var b strings.Builder

    if m.State == LoginMode {
        b.WriteString(Login)
    } else {
        b.WriteString(Register)
    }

    b.WriteString(m.Inputs.Username.View() + "\n")
    b.WriteString(m.Inputs.Password.View() + "\n")

    if m.State == RegisterMode {
        b.WriteString(m.Inputs.Confirm.View() + "\n")
    }

    if m.Error != "" {
        b.WriteString("\n" + m.Error + "\n")
    }

    b.WriteString(HelpFooter)

    return b.String()
}
