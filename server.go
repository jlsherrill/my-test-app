package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)
import "net/http"

type Widget struct {
	Name   string `form:"name"`
	Flavor string `form:"flavor"`
	Id     int64  `form:"id"`
}

//my_widgets := make([]widget, 5, 10)

func main() {
	my_widgets := make(map[int64]Widget)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/widgets/", func(c *gin.Context) {
		list := make([]Widget, 0, len(my_widgets))
		for _, w := range my_widgets {
			list = append(list, w)
		}
		c.JSON(http.StatusOK, list)
	})
	r.POST("/widgets/", func(c *gin.Context) {
		var widget Widget
		if c.BindJSON(&widget) == nil {
			log.Println(widget.Id)
			log.Println(widget.Name)
			my_widgets[widget.Id] = widget

			c.JSON(http.StatusOK, widget)
		}
	})
	r.GET("/widgets/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		widget, found := my_widgets[id]
		if found {
			c.JSON(http.StatusOK, widget)
		} else {
			c.String(404, "Not Found")
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
