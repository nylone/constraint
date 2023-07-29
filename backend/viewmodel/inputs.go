package viewmodel

import "constraint/model"

type ActionId int

const (
	InputAddPos = iota
	InputMsg
)

// use to signal where your client wants a mark to be placed
type Action struct {
	Msg string    `json:"msg"`
	Id  int       `json:"id"`
	Pos model.Pos `json:"pos"`
}
