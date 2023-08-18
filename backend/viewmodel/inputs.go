package viewmodel

import "constraint/model"

// aliased integer type, used to represent the wishes of the user to the game engine
type ActionId int

const (
	InputAddPos ActionId = iota // this value, as an ID, signals to the game engine that the user wishes to add his mark to a specified position
	InputMsg                    // this values, as an ID, signals to the game engine that the user has sent a chat message
)

// use to signal where your client wants a mark to be placed
type Action struct {
	Msg string    `json:"msg"` // if ID is InputMsg then this field is the chat message to be sent
	Id  ActionId  `json:"id"`  // signals the intent of the action message
	Pos model.Pos `json:"pos"` // if ID is InputAddPos then this field is the position to be marked
}
