package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"constraint/view"
	"constraint/viewmodel"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type lobby struct {
	vm *viewmodel.Viewmodel
}

var (
	lobbies      map[string]lobby = make(map[string]lobby)
	lobbiesMutex sync.Mutex
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", struct{}{})
	})

	r.GET("/ws", func(c *gin.Context) {
		lobbyID := c.Query("lobby")
		if lobbyID == "" {
			c.String(http.StatusBadRequest, "missing lobby id")
			return
		}
		lobbiesMutex.Lock()
		defer lobbiesMutex.Unlock()
		if len(lobbies) == 1000 {
			c.String(http.StatusForbidden, "too many lobbies")
			return
		}
		if _, ok := lobbies[lobbyID]; !ok {
			vm := viewmodel.NewViewmodel()
			lobbies[lobbyID] = lobby{
				vm: &vm,
			}
		}

		// get the lobby
		lobby := lobbies[lobbyID]
		vm := lobby.vm

		nickname := c.Query("nickname")
		if nickname == "" {
			c.String(http.StatusBadRequest, "missing nickname")
			return
		}
		// upgrade client connection to websocket
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		go func() {
			view.HandleClient(conn, nickname, vm)
			// when client is done, see if the lobby needs to be freed
			lobbiesMutex.Lock()
			defer lobbiesMutex.Unlock()

			lobby, ok := lobbies[lobbyID]
			if ok {
				lobby.vm.Mutex.Lock()
				defer lobby.vm.Mutex.Unlock()
				if lobby.vm.IsOver {
					delete(lobbies, lobbyID)
					return
				}
			}
		}()
	})

	r.Run("localhost:8080")
}
