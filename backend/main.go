package main

import (
	"constraint/view"
	"constraint/viewmodel"

	"github.com/gin-gonic/gin"
)

func main() {

	vm := viewmodel.NewView()

	r := gin.Default()

	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		err := view.HandleClient(c.Writer, c.Request, "user", &vm)
		if err != nil {
			println(err.Error())
		}
	})

	r.Run("localhost:8080")
}
