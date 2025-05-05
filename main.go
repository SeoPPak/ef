package main

import (
	"log"
	"os"

	"ef/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := loadEnv(); err != nil {
		log.Printf("Warning: .env 파일을 로드하는데 실패했습니다: %v", err)
		log.Println("환경 변수가 직접 설정되어 있다면 이 경고는 무시해도 됩니다.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	router.POST("/categorize", handlers.CategorizeHandler)

	serverAddr := ":" + port
	log.Printf("서버가 %s 포트에서 시작됩니다...", port)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}

func loadEnv() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
			return nil
		}
		log.Println(".env 파일이 없습니다. .env.example 파일을 .env로 복사하고 값을 수정하세요.")
		return nil
	}

	return godotenv.Load()
}