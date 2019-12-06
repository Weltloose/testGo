package views

import (
	"fmt"

	"github.com/Weltloose/testGo/dal/mongodb"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func WsHandler(c *gin.Context) {
	username := c.Query("username")
	groups := mongodb.GetGroupList(username)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("upgrade error ", err)
		return
	}
	client := &mongodb.Client{
		Conn:   conn,
		Hub:    mongodb.MainHub,
		Send:   make(chan []byte, 256),
		Groups: groups,
	}
	client.Hub.Register <- client
	go client.WritePump()

}
