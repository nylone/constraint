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
	Players map[string]model.Mark `json:"players"`
	Field   model.Field           `json:"field"`
	Id      int                   `json:"id"`
}

// signals the start of communications between client and view,
// with info about the game
type NewClientInfo struct {
	Nickname string     `json:"nickname"`
	Mark     model.Mark `json:"mark"`
	Id       int        `json:"id"`
}

type ModelUpdate struct {
	Id     int          `json:"id"`
	Pos    model.Pos    `json:"pos"`
	Mark   model.Mark   `json:"mark"`
	Winner model.Winner `json:"winner"`
}

type ClientLeft struct {
	Nickname string `json:"nickname"`
	Id       int    `json:"id"`
	Shutdown bool   `json:"shutdown"`
}

type ChatMessage struct {
	By  string `json:"by"`
	Msg string `json:"msg"`
	Id  int    `json:"id"`
}
