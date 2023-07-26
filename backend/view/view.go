package view

import (
	"constraint/viewmodel"

	"github.com/gorilla/websocket"
)

// sent only after joining a lobby
type JoinResponse struct {
	Id        int    `json:"id"`
	Succesful bool   `json:"successful"`
	Error     string `json:"error,omitempty"`
}

func HandleClient(conn *websocket.Conn, nick string, vm *viewmodel.Viewmodel) {
	// channel to listen for viewmodel updates
	output := make(chan interface{})
	// listen for client messages and send them to viewmodel
	input, err := vm.AddClient(nick, output)
	if err != nil {
		conn.WriteJSON(JoinResponse{
			Id:        viewmodel.OutputConnected,
			Succesful: false,
			Error:     err.Error(),
		})
		return
	}
	conn.WriteJSON(JoinResponse{
		Id:        viewmodel.OutputConnected,
		Succesful: true,
	})
	// client event listener
	if input != nil {
		go func() {
			var v viewmodel.Action
			for {
				err := conn.ReadJSON(&v)
				if err != nil {
					close(input)
					return
				}
				input <- v
			}
		}()
	}
	// start listening for viewmodel messages and sending them to the client
	for v := range output {
		conn.WriteJSON(v)
	}
}
