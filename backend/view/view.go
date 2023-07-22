package view

import (
	"constraint/viewmodel"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleClient(
	w http.ResponseWriter,
	r *http.Request,
	vm *viewmodel.Viewmodel,
) error {
	// upgrade client connection to websocket
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	// start listening for viewmodel messages and sending them to the client
	output := make(chan interface{})
	go func() {
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
	}()

	// listen for client messages and send them to viewmodel
	input := vm.AddClient(output)
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

	return nil
}
