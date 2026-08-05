package main

import (
	"log"
	"task-management-api/internal/bootstrap"
	"task-management-api/internal/config"
	"task-management-api/internal/database"
	_ "task-management-api/internal/docs"
	"task-management-api/internal/middleware"
	"task-management-api/internal/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Task Management API
// @version 1.0
// @description REST API for task management built with Go, Gin, and GORM.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatal(err)
	}

	if err := database.Seed(db); err != nil {
		log.Fatal(err)
	}

	app := bootstrap.New(db, cfg.JWTSecret)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggerMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Task Management API",
		})
	})

	routes.Register(router, app, cfg.JWTSecret)

	router.Run(":" + cfg.ServerPort)
}
