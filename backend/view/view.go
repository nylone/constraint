package view

import (
	"constraint/viewmodel"

	"github.com/gorilla/websocket"
)

func HandleClient(conn *websocket.Conn, nick string, vm *viewmodel.Viewmodel) {
	// channel to listen for viewmodel updates
	output := make(chan interface{})
	// listen for client messages and send them to viewmodel
	input, err := vm.AddClient(nick, output)
	if err != nil {
		conn.WriteJSON(viewmodel.JoinResponse{
			Id:        viewmodel.CONNECTED,
			Succesful: false,
			Error:     err.Error(),
		})
		return
	}
	conn.WriteJSON(viewmodel.JoinResponse{
		Id:        viewmodel.CONNECTED,
		Succesful: true,
	})
	// client event listener
	if input != nil {
		go func() {
			var v viewmodel.AddPos
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
		err := conn.WriteJSON(v)
		if err != nil {
			panic(err)
		}
	}
	// output has been closed, notify client
	conn.WriteJSON(struct {
		Id string `json:"id"`
	}{Id: "CLOSED"})
}
