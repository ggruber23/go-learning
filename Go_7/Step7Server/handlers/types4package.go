package handlers

import "github.com/gorilla/websocket"

type UserMessage struct {
	UserID  string `json:"userid"`
	Message string `json:"message"`
}

type RegisterCommand struct {
	Name       string
	Connection *websocket.Conn
}
