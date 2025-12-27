package main

import (
	"client/auth"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	urls := auth.AuthURLs {
		BackendURL: "http://localhost:9900/api/1/auth",
		WSURL: 		"ws://localhost:9900/api/1/ws",
	}
	authService := auth.NewAuthService(urls)

	model := auth.NewAuthScreenModel(authService)

	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		panic(err)
	}

}
