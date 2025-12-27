package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
)

const (
    LoginMode    AuthState = "login"
    RegisterMode AuthState = "register"
)

type AuthScreenModel struct {
    Inputs  AuthInputs
    State   AuthState
    Error   string
    Busy    bool
    Service *AuthURLs
}

type AuthURLs struct {
    BackendURL 	string
	WSURL		string
}

type AuthInputs struct {
    Username textinput.Model
    Password textinput.Model
    Confirm  textinput.Model
    FocusIdx int
}

type AuthState string

func NewAuthScreenModel(service *AuthURLs) *AuthScreenModel {
    return &AuthScreenModel{
        Inputs: AuthInputs{
            Username: NewUsernameInput(),
            Password: NewPasswordInput(),
            Confirm:  NewConfirmPasswordInput(),
            FocusIdx: 0,
        },
        State:   LoginMode,
        Service: service,
    }
}
