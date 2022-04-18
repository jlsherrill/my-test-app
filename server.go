package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	kafka "github.com/segmentio/kafka-go"
	"log"
	"strconv"
)
import "net/http"

type Widget struct {
	Name   string `form:"name"`
	Flavor string `form:"flavor"`
	Id     int64  `form:"id"`
}

func getUrl() (string, error) {
	cfg := clowder.LoadedConfig
	if cfg.Kafka.Brokers[0].Hostname != "" {
		return fmt.Sprintf("%s:%v", cfg.Kafka.Brokers[0].Hostname, *cfg.Kafka.Brokers[0].Port), nil
	} else {
		return "", errors.New("empty name")
	}
}

func getTopic() (string, error) {
	topic := "foo"
	for _, topicConfig := range clowder.KafkaTopics {
		topic = topicConfig.Name
	}
	if topic == "" {
		return "", errors.New("empty name")
	} else {
		return topic, nil
	}

}

func getKafkaWriter() *kafka.Writer {
	kafkaUrl, _ := getUrl()
	topic, _ := getTopic()
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaUrl),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func sendMessage(context context.Context, id int64, name string) {
	if clowder.IsClowderEnabled() {
		kafkaWriter := getKafkaWriter()

		msg := kafka.Message{
			Key:   []byte(fmt.Sprintf("id-%s", id)),
			Value: []byte(name),
		}
		err := kafkaWriter.WriteMessages(context, msg)

		if err == nil {
			log.Fatalln("Sent message")
		} else {
			log.Fatalln(err)
		}
	}
}

func main() {
	myWidgets := make(map[int64]Widget)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/widgets/", func(c *gin.Context) {
		list := make([]Widget, 0, len(myWidgets))
		for _, w := range myWidgets {
			list = append(list, w)
		}
		c.JSON(http.StatusOK, list)
	})
	r.POST("/widgets/", func(c *gin.Context) {
		var widget Widget
		if c.BindJSON(&widget) == nil {
			if widget.Name == "send" {
				sendMessage(c, widget.Id, widget.Name)
			}
			log.Println(widget.Id)
			log.Println(widget.Name)
			myWidgets[widget.Id] = widget

			c.JSON(http.StatusOK, widget)
		}
	})
	r.GET("/widgets/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		widget, found := myWidgets[id]
		if found {
			c.JSON(http.StatusOK, widget)
		} else {
			c.String(404, "Not Found")
		}
	})
	r.Run(":8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
