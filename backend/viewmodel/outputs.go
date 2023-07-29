package viewmodel

import (
	"constraint/model"
)

const (
	OutputController = iota
	OutputStarting
	OutputUpdate
	OutputNewClient
	OutputConnected
	OutputClosed
	OutputChatMesage
)

// signals that your action was handled
type ControllerResponse struct {
	Error     string `json:"error,omitempty"`
	Id        int    `json:"id"`
	Succesful bool   `json:"successful"`
}

// signals the start of communications between client and view,
// with info about the game
type StartingInfo struct {
	Field model.Field `json:"field"`
	Id    int         `json:"id"`
	Mark  model.Mark  `json:"mark"`
}

// signals the start of communications between client and view,
// with info about the game
type NewClientInfo struct {
	Nickname string `json:"nickname"`
	Id       int    `json:"id"`
}

type ModelUpdate struct {
	Id     int          `json:"id"`
	Pos    model.Pos    `json:"pos"`
	Winner model.Winner `json:"winner"`
}

type GameClosed struct {
	Id int `json:"id"`
}

type ChatMessage struct {
	By  string `json:"by"`
	Msg string `json:"msg"`
	Id  int    `json:"id"`
}
