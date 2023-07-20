package view

import (
	"constraint/controller"
	"constraint/model"
	"errors"
)

type Ok struct{}
type Close struct{}
type AddPos struct {
	Pos model.Pos
}

func Run(
	controller *controller.Controller,
	input <-chan (AddPos),
	output chan<- (interface{}),
) {
	modelChan := make(chan (model.UpdateMessage))
	mark := controller.AddView(modelChan)

	go func() {
		for {
			select {
			case action, ok := <-input:
				{
					// check if the channel is closed
					if !ok {
						controller.Close()
						return
					}

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
	}()
}
