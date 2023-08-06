package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kkyr/fig"

	"constraint/view"
	"constraint/viewmodel"
)

type Config struct {
	Proxy      string `fig:"proxy" default:"127.0.0.1"`
	Bind       string `fig:"bind" default:"127.0.0.1:8080"`
	FEOrigin   string `fig:"fe_origin" default:"127.0.0.1:80"`
	LobbyLimit int    `fig:"lobby_limit" default:"1000"`
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == cfg.FEOrigin
	},
}

type lobby struct {
	vm *viewmodel.Viewmodel
}

var (
	cfg          Config
	lobbies      map[string]lobby = make(map[string]lobby)
	lobbiesMutex sync.Mutex
)

func init() {
	err := fig.Load(&cfg)
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{cfg.Proxy})

	r.GET("/", func(c *gin.Context) {
		lobbyID := c.Query("lobby")
		if lobbyID == "" {
			c.String(http.StatusBadRequest, "missing lobby id")
			return
		}
		lobbiesMutex.Lock()
		defer lobbiesMutex.Unlock()
		if len(lobbies) >= cfg.LobbyLimit {
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
			c.String(http.StatusNotAcceptable, err.Error())
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

	r.Run(cfg.Bind)
}
