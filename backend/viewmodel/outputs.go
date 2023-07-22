package viewmodel

import "constraint/model"

// signals that your action was handled
type ControllerResponse struct {
	Succesful bool   `json:"successful"`
	Error     string `json:"error"`
}

// signals the start of communications between client and view,
// with info about the game
type StartingInfo struct {
	Field model.Field `json:"field"`
	Mark  model.Mark  `json:"mark"`
}

type ModelUpdate struct {
	Field  model.Field  `json:"field"`
	Winner model.Winner `json:"winner"`
}
