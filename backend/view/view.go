package view

import (
	"constraint/controller"
	"constraint/model"
	"sync"
)

type View struct {
	mutex      sync.Mutex
	model      model.Model
	controller controller.Controller
	outputs    []chan<- (interface{})
	isOver     bool
}

func NewView() View {
	model := model.NewModel(8)
	controller := controller.NewController(&model)

	return View{
		model:      model,
		controller: controller,
	}
}

func (view *View) AddClient(input <-chan (AddPos), output chan<- (interface{})) {
	// determine if client is a player or a spectator
	view.mutex.Lock()
	defer view.mutex.Unlock()
	view.outputs = append(view.outputs, output)
	mark := model.NoMark
	if len(view.outputs) <= 2 {
		mark = model.Mark(len(view.outputs))
	}
	// send info about the game to the client
	output <- StartingInfo{
		Field: view.model.Field,
		Mark:  mark,
	}
	// goroutine to look at client messages
	go func() {
		// as long as the client is connected the loop continues
		for v := range input {
			view.mutex.Lock()

			// call the controller and add the mark
			err := view.controller.AddMark(v.Pos, mark)
			if err != nil {
				output <- err
				continue
			}
			output <- Ok{}
			// if all went well notify every player of the new model state
			modelUpdate := ModelUpdate{
				Field:  view.model.Field,
				Winner: view.model.CheckWinner(),
			}
			for _, c := range view.outputs {
				c <- modelUpdate
			}

			view.mutex.Unlock()
		}
		// if we're here it's because the client has shut down
		// acquire the lock and close the game if not already closed
		view.mutex.Lock()
		defer view.mutex.Unlock()
		if !view.isOver {
			for _, c := range view.outputs {
				close(c)
			}
		}
	}()
}
