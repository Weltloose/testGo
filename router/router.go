package router

import (
	"github.com/Weltloose/testGo/dal/mongodb"
	"github.com/Weltloose/testGo/views"
	"github.com/gin-gonic/gin"
)

func CreateServer() *gin.Engine {
	router := gin.Default()
	router.POST("/api/login", views.Login)
	router.POST("/api/register", views.Register)
	router.POST("/api/createGroup", views.CreateGroup)
	router.POST("/api/joinGroup", views.JoinGroup)
	router.POST("/api/leaveGroup", views.LeaveGroup)
	router.GET("/api/getGroups", views.GetGroupList)
	router.GET("/api/getEventList", views.GetEventList)
	router.POST("/api/addEvent", views.AddEvent)
	router.POST("/api/editEvent", views.EditEvent)
	router.POST("/api/deleteEvent", views.DeleteEvent)

	// use websocket
	go mongodb.MainHub.Run()

	router.GET("/ws", views.WsHandler)
	return router
}
