package main

import (
	"github.com/RustReh/go-project-278/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/ping", handler.Ping)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
