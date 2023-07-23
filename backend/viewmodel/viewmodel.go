package viewmodel

import (
	"constraint/controller"
	"constraint/model"
	"errors"
	"sync"
)

type Viewmodel struct {
	mutex      sync.Mutex
	model      model.Model
	controller controller.Controller
	outputs    map[string]chan<- (interface{})
	isOver     bool
}

func NewView() Viewmodel {
	model := model.NewModel(8)
	controller := controller.NewController(&model)

	return Viewmodel{
		model:      model,
		controller: controller,
		outputs:    make(map[string]chan<- interface{}),
	}
}

func (viewmodel *Viewmodel) AddClient(nickname string, output chan<- (interface{})) (chan<- (AddPos), error) {
	// determine if client is a player or a spectator
	viewmodel.mutex.Lock()
	defer viewmodel.mutex.Unlock()
	// if the game is over don't add a new client
	if viewmodel.isOver {
		return nil, errors.New("game over")
	}
	// check if nickname is already in use
	if _, ok := viewmodel.outputs[nickname]; ok {
		return nil, errors.New("nickname in use")
	}
	// tell other clients that a new player has joined
	for _, c := range viewmodel.outputs {
		c <- NewClientInfo{
			Id:       NEWCLIENT,
			Nickname: nickname,
		}
	}
	// subscribe new view with nickname
	viewmodel.outputs[nickname] = output
	// decide if view is player or spectator
	mark := model.NoMark
	if len(viewmodel.outputs) <= 2 {
		mark = model.Mark(len(viewmodel.outputs))
	}
	// send info about the game to the client asyncronously
	go func() {
		viewmodel.mutex.Lock()
		defer viewmodel.mutex.Unlock()
		output <- StartingInfo{
			Id:    STARTING,
			Field: viewmodel.model.Field,
			Mark:  mark,
		}
	}()
	// if we are a spectator, don't create an event listener
	if mark == model.NoMark {
		return nil, nil
	}
	// create the input channel
	input := make(chan (AddPos))
	// client event listener
	go func() {
		// as long as the client is connected the loop continues
		for in := range input {
			viewmodel.mutex.Lock()
			// call the controller and add the mark
			err := viewmodel.controller.AddMark(in.Pos, mark)
			if err != nil {
				output <- ControllerResponse{
					Id:        CONTROLLER,
					Succesful: false,
					Error:     err.Error(),
				}
				viewmodel.mutex.Unlock()
				continue
			}
			output <- ControllerResponse{
				Id:        CONTROLLER,
				Succesful: true,
			}
			// if all went well notify every player of the new model state
			modelUpdate := ModelUpdate{
				Id:     UPDATE,
				Pos:    in.Pos,
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
	// return input channel to the client
	return input, nil
}
