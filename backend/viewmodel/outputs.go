package viewmodel

import "constraint/model"

const (
	CONTROLLER = "CONTROLLER"
	STARTING   = "STARTING"
	UPDATE     = "UPDATE"
)

// signals that your action was handled
type ControllerResponse struct {
	Id        string `json:"id"`
	Succesful bool   `json:"successful"`
	Error     string `json:"error,omitempty"`
}

// signals the start of communications between client and view,
// with info about the game
type StartingInfo struct {
	Id    string      `json:"id"`
	Field model.Field `json:"field"`
	Mark  model.Mark  `json:"mark"`
}

type ModelUpdate struct {
	Id     string       `json:"id"`
	Pos    model.Pos    `json:"pos"`
	Winner model.Winner `json:"winner"`
}
