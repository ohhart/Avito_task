package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршруты и возвращает готовый маршрутизатор Gin
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Обработчики маршрутов
	router.GET("/ping", pingHandler)

	return router
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
