package viewmodel

import "constraint/model"

// use to signal where your client wants a mark to be placed
type AddPos struct {
	Pos model.Pos `json:"pos"`
}
