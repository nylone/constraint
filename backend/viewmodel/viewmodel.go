package viewmodel

import (
	"errors"
	"sync"

	"constraint/controller"
	"constraint/model"
)

type player struct {
	output chan<- interface{}
	mark   model.Mark
}

type Viewmodel struct {
	players    map[string]player
	controller controller.Controller
	model      model.Model
	Mutex      sync.Mutex
	IsOver     bool
}

func NewViewmodel() Viewmodel {
	model := model.NewModel(8)
	controller := controller.NewController(&model)

	return Viewmodel{
		model:      model,
		controller: controller,
		players:    make(map[string]player),
	}
}

func (viewmodel *Viewmodel) AddClient(nickname string, output chan<- interface{}) (chan<- Action, error) {
	// determine if client is a player or a spectator
	viewmodel.Mutex.Lock()
	defer viewmodel.Mutex.Unlock()
	// if the game is over don't add a new client
	if viewmodel.IsOver {
		return nil, errors.New("game over")
	}
	// check if nickname is already in use
	if _, ok := viewmodel.players[nickname]; ok {
		return nil, errors.New("nickname in use")
	}
	// decide if view is player or spectator
	mark := model.NoMark
	if len(viewmodel.players) < 2 {
		mark = model.Mark(len(viewmodel.players) + 1)
	}
	// tell other clients that a new player has joined
	for _, p := range viewmodel.players {
		p.output <- NewClientInfo{
			Id:       OutputNewClient,
			Mark:     mark,
			Nickname: nickname,
		}
	}
	// subscribe new view with nickname
	viewmodel.players[nickname] = player{mark: mark, output: output}
	// send info about the game to the client asyncronously
	go func() {
		viewmodel.Mutex.Lock()
		defer viewmodel.Mutex.Unlock()
		players := make(map[string]model.Mark)

		for k, p := range viewmodel.players {
			players[k] = p.mark
		}

		output <- StartingInfo{
			Players: players,
			Id:      OutputStarting,
			Field:   viewmodel.model.Field,
		}
	}()
	// create the input channel
	input := make(chan Action)
	// client event listener
	go func() {
		// as long as the client is connected the loop continues
		for in := range input {
			viewmodel.Mutex.Lock()
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
						Mark:   mark,
						Winner: viewmodel.model.CheckWinner(),
					}
					for _, p := range viewmodel.players {
						p.output <- modelUpdate
					}
				}
			case InputMsg:
				{
					msg := ChatMessage{
						Id:  OutputChatMesage,
						By:  nickname,
						Msg: in.Msg,
					}
					for _, p := range viewmodel.players {
						p.output <- msg
					}
				}
			}
			viewmodel.Mutex.Unlock()
		}
		// if we're here it's because the client has shut down
		// acquire the lock and close the game if not already closed
		viewmodel.Mutex.Lock()
		defer viewmodel.Mutex.Unlock()

		if !viewmodel.IsOver {
			if mark == model.NoMark {
				delete(viewmodel.players, nickname)
				close(output)
				for _, p := range viewmodel.players {
					p.output <- ClientLeft{
						Id:       OutputClosed,
						Shutdown: false,
						Nickname: nickname,
					}
				}
				return
			}

			viewmodel.IsOver = true
			for _, p := range viewmodel.players {
				p.output <- ClientLeft{
					Id:       OutputClosed,
					Shutdown: true,
					Nickname: nickname,
				}
				close(p.output)
			}
		}
	}()
	// return input channel to the client
	return input, nil
}
