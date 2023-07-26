package main

import (
	"constraint/view"
	"constraint/viewmodel"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type lobby struct {
	vm          *viewmodel.Viewmodel
	clientCount int
}

var (
	lobbies      map[string]lobby = make(map[string]lobby)
	lobbiesMutex sync.Mutex
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.LoadHTMLFiles("index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		nickname := c.DefaultQuery("nickname", "Guest "+time.Now().String())
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
		// upgrade client connection to websocket
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		lobby := lobbies[lobbyID]
		vm := lobby.vm
		lobby.clientCount++
		lobbies[lobbyID] = lobby

		go func() {
			view.HandleClient(conn, nickname, vm)

			lobbiesMutex.Lock()
			defer lobbiesMutex.Unlock()
			lobby := lobbies[lobbyID]
			lobby.clientCount--
			if lobby.clientCount == 0 {
				delete(lobbies, lobbyID)
				return
			}
			lobbies[lobbyID] = lobby
			// when client is done, see if the lobby needs to be freed
		}()
	})

	r.Run("localhost:8080")
}
