package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Username    string `json:"username"`
	MsgContents string `json:"content"`
	ChatID      int    `json:"chatID,string"`
}

type Chat struct {
	gorm.Model
}

func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*")

	m := melody.New()
	dsn := "host=localhost user=ted password=topsecretpassword dbname=chats"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to datbase :(")
	}

	db.AutoMigrate(&Message{})

	router.GET("/chats/:id", func(c *gin.Context) {
		chatID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error_message": fmt.Sprintf("No chat found with ID %v", c.Param("id"))})
			return
		}

		var messages []Message
		db.Where(&Message{ChatID: chatID}).Order("id desc").Find(&messages)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"messages": messages,
		})
	})

	router.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			panic(err)
		}

		db.Create(&message)

		m.Broadcast(msg)
	})

	router.Run(":5000")
}
