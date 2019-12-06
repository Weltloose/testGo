package main

import (
	"github.com/Weltloose/testGo/router"
)

func main() {
	server := router.CreateServer()
	server.Run()
}
