package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strconv"
)
import "net/http"

type Widget struct {
	gorm.Model
	Name   string `form:"name"`
	Flavor string `form:"flavor"`
	Id     int64  `form:"id"`
}

func apiServer() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Widget{})

	//myWidgets := make(map[int64]Widget)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/widgets/", func(c *gin.Context) {
		//list := make([]Widget, 0, len(myWidgets))
		//for _, w := range myWidgets {
		//	list = append(list, w)
		//}
		var list []Widget
		db.Find(&list)
		c.JSON(http.StatusOK, list)
	})
	r.POST("/widgets/", func(c *gin.Context) {
		var widget Widget
		if c.BindJSON(&widget) == nil {
			log.Println(widget.Id)
			log.Println(widget.Name)
			if result := db.Create(&widget); result.Error != nil {
				log.Println(result.Error)
				c.JSON(http.StatusInternalServerError, "")
			} else {
				c.JSON(http.StatusOK, widget)
			}

		}
	})
	r.GET("/widgets/:id", func(c *gin.Context) {
		var widget Widget
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

		result := db.Where("id = ?", id).First(&widget)
		//widget, found := myWidgets[id]
		if result.RowsAffected == 1 {
			c.JSON(http.StatusOK, widget)
		} else {
			c.String(404, "Not Found")
		}
	})
	r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func main() {
	apiServer()
}
