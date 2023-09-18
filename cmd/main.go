package main

import (
	"server/database"
	"server/handlers"
	"server/middlewares"
	"server/repositories"
	"server/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	// configure logger
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "04/04/2001 15:00",
	})

	// load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectToDatabase()
	database.CreateDatabaseTables()
}

//	@title			Gerenciario Backend
//	@version		1.0
//	@description	REST API to help small business owners manage their business.

//	@contact.name	Antônio Pedro Rodrigues Santos
//	@contact.email	antoniopedro.rs9@gmail.com

//	@host		localhost:8080
//	@BasePath	/

func main() {
	userRepository := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	r := gin.Default()

	userGroup := r.Group("/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		userGroup.GET("/", userHandler.ListUsers)
		userGroup.PUT("/", middlewares.RequireAuth, userHandler.UpdateUser)
		userGroup.DELETE("/:id", middlewares.RequireAuth, userHandler.DeleteUser)
		userGroup.POST("/login", userHandler.LoginUser)
	}

	r.Run()
}