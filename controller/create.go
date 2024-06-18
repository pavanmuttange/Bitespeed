package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pavanmuttange/Bitespeed/model"
	"github.com/pavanmuttange/Bitespeed/pkg/config"
)

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ping") //for API testing
}

func CreateTable(ctx *gin.Context) {

	contact := new(model.Contact)

	config.DB.Migrator().CreateTable(&contact)

	ctx.JSON(http.StatusOK, "table created successfully")
}
