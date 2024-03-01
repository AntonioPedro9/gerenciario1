package main

import (
	"server/database"
	"server/handlers"
	"server/middlewares"
	"server/repositories"
	"server/services"

	"github.com/gin-contrib/cors"
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
}

func main() {
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"localhost"})

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	db := database.ConnectToDatabase()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// database.ClearTestDatabase(db)
	database.CreateDatabaseTables(db)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", userHandler.CreateUser)
		// userGroup.GET("/list", userHandler.ListUsers)
		userGroup.PATCH("/", middlewares.RequireAuth, userHandler.UpdateUser)
		userGroup.DELETE("/:id", middlewares.RequireAuth, userHandler.DeleteUser)
		userGroup.POST("/login", userHandler.LoginUser)
	}

	customerRepository := repositories.NewCustomerRepository(db)
	customerService := services.NewCustomerService(customerRepository)
	customerHandler := handlers.NewCustomerHandler(customerService)
	customerGroup := r.Group("/customers")
	{
		customerGroup.POST("/", middlewares.RequireAuth, customerHandler.CreateCustomer)
		customerGroup.GET("/list/", middlewares.RequireAuth, customerHandler.ListCustomers)
		customerGroup.GET("/:customerID", middlewares.RequireAuth, customerHandler.GetCustomer)
		customerGroup.PATCH("/", middlewares.RequireAuth, customerHandler.UpdateCustomer)
		customerGroup.DELETE("/:customerID", middlewares.RequireAuth, customerHandler.DeleteCustomer)
	}

	jobRepository := repositories.NewJobRepository(db)
	jobService := services.NewJobService(jobRepository)
	jobHandler := handlers.NewJobHandler(jobService)
	jobGroup := r.Group("/jobs")
	{
		jobGroup.POST("/", middlewares.RequireAuth, jobHandler.CreateJob)
		jobGroup.GET("/list/", middlewares.RequireAuth, jobHandler.ListJobs)
		jobGroup.GET("/:jobID", middlewares.RequireAuth, jobHandler.GetJob)
		jobGroup.PUT("/", middlewares.RequireAuth, jobHandler.UpdateJob)
		jobGroup.DELETE("/:jobID", middlewares.RequireAuth, jobHandler.DeleteJob)
	}

	budgetRepository := repositories.NewBudgetRepository(db)
	budgetService := services.NewBudgetService(budgetRepository)
	budgetHandler := handlers.NewBudgetHandler(budgetService)
	budgetGroup := r.Group("/budgets")
	{
		budgetGroup.POST("/", middlewares.RequireAuth, budgetHandler.CreateBudget)
		budgetGroup.GET("/list/", middlewares.RequireAuth, budgetHandler.ListBudgets)
		budgetGroup.GET("/:budgetID", middlewares.RequireAuth, budgetHandler.GetBudget)
		budgetGroup.GET("/list/jobs/:budgetID", middlewares.RequireAuth, budgetHandler.GetBudgetJobs)
		budgetGroup.DELETE("/:budgetID", middlewares.RequireAuth, budgetHandler.DeleteBudget)
	}

	r.Run()
}
