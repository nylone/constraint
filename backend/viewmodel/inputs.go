package viewmodel

import "constraint/model"

type ActionId int

const (
	InputAddPos = iota
	InputMsg
)

// use to signal where your client wants a mark to be placed
type Action struct {
	Id  int       `json:"id"`
	Pos model.Pos `json:"pos"`
	Msg string    `json:"msg"`
}
