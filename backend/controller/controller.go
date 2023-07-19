package controller

import (
	"constraint/model"
	"errors"
)

type Winner int8

const (
	NoWinner Winner = iota
	Player1
	Player2
)

type Controller struct {
	model model.Model
}

func NewController(size uint) Controller {
	return Controller{
		model: model.NewModel(size),
	}
}

func (controller Controller) isSafePos(pos model.Pos) bool {
	return pos.X > 0 && pos.Y > 0 && pos.X < int(controller.model.Size) && pos.Y < int(controller.model.Size)
}

func (controller Controller) isMarkAdjacent(pos model.Pos) bool {
	varX := controller.model.History[len(controller.model.History)-1].X - pos.X
	varY := controller.model.History[len(controller.model.History)-1].Y - pos.Y
	return varX < -1 || varX > 1 || varY < -1 || varY > 1
}

func (controller Controller) AddMark(pos model.Pos, mark model.Mark) error {
	// special rules apply if it's the first round
	if len(controller.model.History) == 0 {
		// check to see if it's the correct turn
		if mark == model.P2Mark {
			return errors.New("wrong turn")
		}
		// check to see if it's the position is safe
		if !controller.isSafePos(pos) {
			return errors.New("unsafe position")
		}
		// finalize move
		controller.model.PutMark(pos, mark)
		return nil
	}
	// check to see if it's the correct turn
	lastPos := controller.model.History[len(controller.model.History)-1]
	if mark == controller.model.Field[lastPos.X][lastPos.Y] {
		return errors.New("wrong turn")
	}
	// check to see if it's the position is safe
	if !controller.isSafePos(pos) {
		return errors.New("unsafe position")
	}
	// check to see il move is adjacent the previous one
	if controller.isMarkAdjacent(pos) {
		return errors.New("invalid position")
	}
	// finalize move
	controller.model.PutMark(pos, mark)
	return nil
}

func (controller Controller) CheckWinner() Winner {
	lastPos := controller.model.History[len(controller.model.History)-1]
	posToCheck := [8]model.Pos{
		{X: lastPos.X - 1, Y: lastPos.Y - 1},
		{X: lastPos.X + 0, Y: lastPos.Y - 1},
		{X: lastPos.X + 1, Y: lastPos.Y - 1},
		{X: lastPos.X - 1, Y: lastPos.Y + 0},
		{X: lastPos.X + 1, Y: lastPos.Y + 0},
		{X: lastPos.X - 1, Y: lastPos.Y + 1},
		{X: lastPos.X + 0, Y: lastPos.Y + 1},
		{X: lastPos.X + 1, Y: lastPos.Y + 1},
	}
	for _, p := range posToCheck {
		// check if p is out of bounds, if so skip it
		if p.X < 0 || p.X > int(controller.model.Size) || p.Y < 0 || p.Y > int(controller.model.Size) {
			continue
		}
		// p is guaranteed in bounds, check to see if it is occupied
		if controller.model.Field[p.X][p.Y] == model.NoMark {
			return NoWinner
		}
	}
	// all positions are occupied, there is a winner
	if controller.model.Field[lastPos.X][lastPos.Y] == model.P1Mark {
		return Player2
	} else {
		return Player1
	}
}
