package main

import (
    "github.com/gin-gonic/gin"
    "C:\Users\marwq\Avito_task\internal\api" // Импортируйте пакет api из вашего проекта
)

func main() {
    // Инициализация Gin router
    router := api.SetupRouter()

    // Запуск HTTP сервера
    router.Run(":8080")
}
