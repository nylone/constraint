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

// sent only after joining a lobby
type JoinResponse struct {
	Id        int    `json:"id"`
	Succesful bool   `json:"successful"`
	Error     string `json:"error,omitempty"`
}

// signals that your action was handled
type ControllerResponse struct {
	Id        int    `json:"id"`
	Succesful bool   `json:"successful"`
	Error     string `json:"error,omitempty"`
}

// signals the start of communications between client and view,
// with info about the game
type StartingInfo struct {
	Id    int         `json:"id"`
	Field model.Field `json:"field"`
	Mark  model.Mark  `json:"mark"`
}

// signals the start of communications between client and view,
// with info about the game
type NewClientInfo struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
}

type ModelUpdate struct {
	Id     int          `json:"id"`
	Pos    model.Pos    `json:"pos"`
	Winner model.Winner `json:"winner"`
}

type GameClosed struct {
	Id int `json:"id"`
}

type ChatMesage struct {
	Id  int    `json:"id"`
	Msg string `json:"msg"`
}
