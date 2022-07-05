package main

import (
	"fmt"
	"log"
	"os"

	"github.com/darolpz/students/cmd/handlers"
	"github.com/darolpz/students/internal/auth"
	"github.com/darolpz/students/internal/database"
	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

type services struct {
	studentRepository repository.IStudentsRepository
	userRepository    repository.IUsersRepository
	authService       auth.IAuthService
}

func main() {
	services := initServices()

	app := gin.Default()
	handlers.CreateHealthEndpoints(app)
	handlers.CreateStudentsEndpoints(app, services.studentRepository, services.authService)
	handlers.CreateAuthEndpoints(app, services.userRepository, services.authService)

	// listen and serve on
	app.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func initServices() *services {
	services := &services{}
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	databaseService, err := database.NewDatabaseService(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}

	services.studentRepository = repository.NewStudentsRepo(databaseService)
	services.userRepository = repository.NewUsersRepository(databaseService)

	services.authService = auth.NewAuthService(os.Getenv("JWT_SECRET"))
	return services
}
