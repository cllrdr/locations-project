// @title           Locations Project API
// @version         1.0
// @description     REST API для управления локациями и игровыми заявками
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"os"

	"locations-project/internal/app/config"
	"locations-project/internal/app/dsn"
	"locations-project/internal/app/handler"
	"locations-project/internal/app/repository"
	"locations-project/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "locations-project/docs"
)

func main() {
	router := gin.Default()

	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	// MinIO config (can be overridden via env)
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minio")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minio124")
	bucketName := getEnv("MINIO_BUCKET", "locations")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	rep, err := repository.New(postgresString, minioEndpoint, minioAccessKey, minioSecretKey, bucketName, useSSL)
	if err != nil {
		logrus.Fatalf("error initializing repository: %v", err)
	}

	hand := handler.NewHandler(rep, conf)

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запускаем приложение (внутри будет вызван RegisterAPI)
	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}