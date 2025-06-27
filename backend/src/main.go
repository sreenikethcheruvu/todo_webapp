package main

import (
    "log"
    "todoapp-backend/src"
)

func main() {
    src.InitDB()

    router := src.SetupRoutes()

    if err := router.Run(":8080"); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}