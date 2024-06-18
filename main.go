package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pavanmuttange/Bitespeed/controller"
	"github.com/pavanmuttange/Bitespeed/pkg/config"
)

func main() {

	config.Init()
	router := gin.Default()
	router.GET("/health", controller.Ping)
	router.GET("/create", controller.CreateTable)
	router.POST("/identify", controller.Identify)

	router.Run(":8002")
}
