package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"ef/handlers"
)

func main() {
	router := gin.Default()

	router.POST("/categorize", handlers.CategorizeHandler)

	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}