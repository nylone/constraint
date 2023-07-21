package view

import "constraint/model"

// signals that your action was handled successfully
type Ok struct{}

// signals the start of communications between client and view,
// with info about the game
type StartingInfo struct {
	Field model.Field
	Mark  model.Mark
}

type ModelUpdate struct {
	Field  model.Field
	Winner model.Winner
}
