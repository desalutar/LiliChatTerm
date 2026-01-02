package chat

import (
	"client/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type User struct {
	ID 			int64 	`json:"id"`
	Username 	string 	`json:"username"`
}

func SearchUser(token, username string) (User, error) {
	var user User
	url := utils.BaseURL + utils.UserEndpoint + username
	err := utils.GetJson(url, token, &user)

	return user,  err
}

func searchUserCmd(token, username string) tea.Cmd {
	return func() tea.Msg {
		user, err := SearchUser(token, username)
		if err != nil {
			return searchResultMsg{Err: err}
		}
		return searchResultMsg{UserID: user.ID, Username: username}
	}
}