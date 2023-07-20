package view

import (
	"constraint/controller"
	"constraint/model"
	"errors"
)

// signals that your action was handled successfully
type Ok struct{}

// signals the MVC is shutting down
type Close struct{}

// use to signal where your client wants a mark to be placed
type AddPos struct {
	Pos model.Pos
}

// signals the start of communications between client and view,
// with info about the game
type StartingInfo struct {
	model.Field
	model.Mark
}

func Run(
	controller *controller.Controller,
	input <-chan (AddPos),
	output chan<- (interface{}),
) {
	modelChan := make(chan (model.UpdateMessage))
	mark, field := controller.AddView(modelChan)
	output <- StartingInfo{
		Mark:  mark,
		Field: field,
	}
	for {
		select {
		case action, ok := <-input:
			{
				// check if the channel is closed
				if !ok {
					controller.Close()
					return
				}
				// spectators are handled by the view only
				if mark == model.NoMark {
					output <- errors.New("invalid action")
					continue
				}

				err := controller.AddMark(action.Pos, mark)
				if err != nil {
					output <- err
					continue
				}
				output <- (Ok{})
			}
		case modelUpdate, ok := <-modelChan:
			{
				// check if the channel is closed
				if !ok {
					output <- Close{}
					return
				}

				output <- modelUpdate
			}
		}
	}
}
