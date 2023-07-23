package view

import (
	"constraint/viewmodel"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleClient(
	w http.ResponseWriter,
	r *http.Request,
	nick string,
	vm *viewmodel.Viewmodel,
) error {
	// upgrade client connection to websocket
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	// used to syncronize the connection successful message with the rest of the application
	writerMutex := sync.Mutex{}
	writerMutex.Lock()
	defer writerMutex.Unlock()

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
		return nil
	}
	conn.WriteJSON(viewmodel.JoinResponse{
		Id:        viewmodel.CONNECTED,
		Succesful: true,
	})
	// start listening for viewmodel messages and sending them to the client
	go func() {
		writerMutex.Lock()
		defer writerMutex.Unlock()
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

	return nil
}
