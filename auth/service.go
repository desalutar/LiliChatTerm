package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewAuthService(url AuthURLs) *AuthURLs {
    return &AuthURLs{
		BackendURL: url.BackendURL,
		WSURL: url.WSURL,
	}
}

type AuthResponse struct {
	Message     string
	UserID      int64
	AccessToken string
}

var LoginResp struct {
	UserID      int64    `json:"user_id"`
	AccessToken string `json:"access_token"`
}

func (s *AuthURLs) Login(username, password string) (*AuthResponse, error) {
	reqBody := map[string]string{
		"username": username,
		"password": password,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	resp, err := http.Post(s.BackendURL+"/login", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", data)
	}

	if err := json.NewDecoder(resp.Body).Decode(&LoginResp); err != nil {
		return nil, err
	}

	return &AuthResponse{
		UserID:      LoginResp.UserID,
		AccessToken: LoginResp.AccessToken,
	}, nil
}

func (au *AuthURLs) Register(username, password string) (*AuthResponse, error) {
	reqBody := map[string]string{
		"username": username,
		"password": password,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	resp, err := http.Post(au.BackendURL+"/register", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("register failed: %s", data)
	}

	var registerResp struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		return nil, err
	}

	return &AuthResponse{Message: registerResp.Message}, nil
}
