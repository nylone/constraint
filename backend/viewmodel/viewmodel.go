package viewmodel

import (
	"errors"
	"sync"

	"constraint/controller"
	"constraint/model"
)

type Viewmodel struct {
	outputs    map[string]chan<- interface{}
	controller controller.Controller
	model      model.Model
	mutex      sync.Mutex
	isOver     bool
}

func NewViewmodel() Viewmodel {
	model := model.NewModel(8)
	controller := controller.NewController(&model)

	return Viewmodel{
		model:      model,
		controller: controller,
		outputs:    make(map[string]chan<- interface{}),
	}
}

func (viewmodel *Viewmodel) AddClient(nickname string, output chan<- interface{}) (chan<- Action, error) {
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
			Id:       OutputNewClient,
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
			Id:    OutputStarting,
			Field: viewmodel.model.Field,
			Mark:  mark,
		}
	}()
	// create the input channel
	input := make(chan Action)
	// client event listener
	go func() {
		// as long as the client is connected the loop continues
		for in := range input {
			viewmodel.mutex.Lock()
			switch in.Id {
			case InputAddPos:
				{
					// call the controller and add the mark
					err := viewmodel.controller.AddMark(in.Pos, mark)
					if err != nil {
						output <- ControllerResponse{
							Id:        OutputController,
							Succesful: false,
							Error:     err.Error(),
						}
						break
					}
					output <- ControllerResponse{
						Id:        OutputController,
						Succesful: true,
					}
					// if all went well notify every player of the new model state
					modelUpdate := ModelUpdate{
						Id:     OutputUpdate,
						Pos:    in.Pos,
						Winner: viewmodel.model.CheckWinner(),
					}
					for _, c := range viewmodel.outputs {
						c <- modelUpdate
					}
				}
			case InputMsg:
				{
					msg := ChatMessage{
						Id:  OutputChatMesage,
						By:  nickname,
						Msg: in.Msg,
					}
					for _, c := range viewmodel.outputs {
						c <- msg
					}
				}
			}
			viewmodel.mutex.Unlock()
		}
		// if we're here it's because the client has shut down
		// acquire the lock and close the game if not already closed
		viewmodel.mutex.Lock()
		defer viewmodel.mutex.Unlock()
		if !viewmodel.isOver {
			for _, c := range viewmodel.outputs {
				c <- GameClosed{Id: OutputClosed}
				close(c)
			}
			viewmodel.isOver = true
		}
	}()
	// return input channel to the client
	return input, nil
}
