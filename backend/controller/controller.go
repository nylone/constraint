package controller

import (
	"constraint/model"
	"errors"
	"sync"
)

type Controller struct {
	model model.Model
	mutex sync.Mutex
}

func NewController(size uint) Controller {
	return Controller{
		model: model.NewModel(size),
	}
}

func (controller *Controller) AddView(c chan<- (model.UpdateMessage)) (model.Mark, model.Field) {
	// acquire lock on mutex
	controller.mutex.Lock()
	defer controller.mutex.Unlock()

	return controller.model.AddListener(c)
}

func isSafePos(pos model.Pos, size uint) bool {
	return pos.X > 0 && pos.Y > 0 && pos.X < int(size) && pos.Y < int(size)
}

func isMarkAdjacent(pos model.Pos, lastPos model.Pos) bool {
	varX := lastPos.X - pos.X
	varY := lastPos.Y - pos.Y
	return varX < -1 || varX > 1 || varY < -1 || varY > 1
}

func (controller *Controller) AddMark(pos model.Pos, mark model.Mark) error {
	// acquire lock on mutex
	controller.mutex.Lock()
	defer controller.mutex.Unlock()

	//check if game is over
	if controller.model.CheckWinner() != model.NoWinner {
		return errors.New("game over")
	}
	// special rules apply if it's the first round
	if len(controller.model.History) == 0 {
		// check to see if it's the correct turn
		if mark == model.P2Mark {
			return errors.New("wrong turn")
		}
		// check to see if it's the position is safe
		if !isSafePos(pos, controller.model.Size) {
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
	if !isSafePos(pos, controller.model.Size) {
		return errors.New("unsafe position")
	}
	// check to see il move is adjacent the previous one
	if isMarkAdjacent(pos, controller.model.History[len(controller.model.History)-1]) {
		return errors.New("invalid position")
	}
	// finalize move
	controller.model.PutMark(pos, mark)
	return nil
}

func (controller *Controller) Close() {
	// acquire lock on mutex
	controller.mutex.Lock()
	defer controller.mutex.Unlock()

	controller.model.Close()
}
