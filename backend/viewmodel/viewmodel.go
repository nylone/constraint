package viewmodel

import (
	"constraint/controller"
	"constraint/model"
	"sync"
)

type Viewmodel struct {
	mutex      sync.Mutex
	model      model.Model
	controller controller.Controller
	outputs    []chan<- (interface{})
	isOver     bool
}

func NewView() Viewmodel {
	model := model.NewModel(8)
	controller := controller.NewController(&model)

	return Viewmodel{
		model:      model,
		controller: controller,
	}
}

func (viewmodel *Viewmodel) AddClient(output chan<- (interface{})) chan<- (AddPos) {
	// determine if client is a player or a spectator
	viewmodel.mutex.Lock()
	defer viewmodel.mutex.Unlock()
	viewmodel.outputs = append(viewmodel.outputs, output)
	mark := model.NoMark
	if len(viewmodel.outputs) <= 2 {
		mark = model.Mark(len(viewmodel.outputs))
	}
	// send info about the game to the client
	defer func() {
		output <- StartingInfo{
			Field: viewmodel.model.Field,
			Mark:  mark,
		}
	}()

	// if we are a spectator, don't create an event listener
	if mark == model.NoMark {
		return nil
	}
	// create the input channel
	input := make(chan (AddPos))

	// goroutine to look at client events
	go func() {
		// as long as the client is connected the loop continues
		for v := range input {
			viewmodel.mutex.Lock()

			// call the controller and add the mark
			err := viewmodel.controller.AddMark(v.Pos, mark)
			if err != nil {
				output <- ControllerResponse{Succesful: false, Error: err.Error()}
				continue
			}
			output <- ControllerResponse{Succesful: true}
			// if all went well notify every player of the new model state
			modelUpdate := ModelUpdate{
				Field:  viewmodel.model.Field,
				Winner: viewmodel.model.CheckWinner(),
			}
			for _, c := range viewmodel.outputs {
				c <- modelUpdate
			}

			viewmodel.mutex.Unlock()
		}
		// if we're here it's because the client has shut down
		// acquire the lock and close the game if not already closed
		viewmodel.mutex.Lock()
		defer viewmodel.mutex.Unlock()
		if !viewmodel.isOver {
			for _, c := range viewmodel.outputs {
				close(c)
			}
			viewmodel.isOver = true
		}
	}()

	return input
}
