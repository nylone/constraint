package viewmodel

import (
	"constraint/model"
)

// aliased integer value, used to represent the fields inside a message sent to the user
type OutputId int

const (
	OutputController OutputId = iota // signals ControllerResponse
	OutputStarting                   // signals StartingInfo
	OutputUpdate                     // signals ModelUpdate
	OutputNewClient                  // signals NewClientInfo
	OutputConnected                  // signals JoinResponse
	OutputClosed                     // signals ClientLeft
	OutputChatMesage                 // signals ChatMessage
)

// signals that your action was handled
type ControllerResponse struct {
	Error     string   `json:"error,omitempty"` // only present if Successful is false, signals the error from the controller
	Id        OutputId `json:"id"`              // can only be OutputController
	Succesful bool     `json:"successful"`      // if true, the action was carried out with no issues
}

// signals the start of communications between client and view,
// with info about the game
type StartingInfo struct {
	Players map[string]model.Mark `json:"players"` // list of all nickanames mapped to their respective marks, player included
	Field   model.Field           `json:"field"`   // current state of the game field at the time of joining the game
	Id      OutputId              `json:"id"`      // can only be OutputStarting
}

// signals the start of communications between client and view,
// with info about the game
type NewClientInfo struct {
	Nickname string     `json:"nickname"` // nickname of the new client that connected
	Mark     model.Mark `json:"mark"`     // signals the mark of the new client
	Id       OutputId   `json:"id"`       // can only be OutputNewClient
}

// signals the underlying model has changed.
// contains info about the new update to the position, the mark and
// the winner, if present
type ModelUpdate struct {
	Id     OutputId     `json:"id"`     // can only be OutputUpdate
	Pos    model.Pos    `json:"pos"`    // is position that just got updated
	Mark   model.Mark   `json:"mark"`   // the mark that has been put at Pos
	Winner model.Winner `json:"winner"` // if this was a winning move, reports the winner
}

// signals the end of communications between client and view
type ClientLeft struct {
	Nickname string   `json:"nickname"` // nickname of the user that left
	Id       OutputId `json:"id"`       // can only be OutputClosed
	Shutdown bool     `json:"shutdown"` // if true, the game is over
}

// signals the presence of a new messsage in the chat
type ChatMessage struct {
	By  string   `json:"by"`  // nickname of the user that sent the message
	Msg string   `json:"msg"` // message in the chat
	Id  OutputId `json:"id"`  // can only be OutputChatMesage
}

// sent only after joining a lobby
type JoinResponse struct {
	Error     string   `json:"error,omitempty"` // upon a failed connection, reports the error in connecting to the lobby
	Succesful bool     `json:"successful"`      // if false the connection is to be closed
	Id        OutputId `json:"id"`              // can only be OutputConnected
}
