package main

import (
	"fmt"
	"log"
	"os"

	"github.com/darolpz/students/cmd/handlers"
	_ "github.com/darolpz/students/docs"
	"github.com/darolpz/students/internal/auth"
	"github.com/darolpz/students/internal/database"
	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type services struct {
	studentRepository repository.IStudentsRepository
	userRepository    repository.IUsersRepository
	authService       auth.IAuthService
}

// @title           darolpz students
// @version         0.1
// @description     This is an example CRUD project.

// @contact.name   Dario Lopez
// @contact.url    http://www.github.com/darolpz
// @contact.email  daropl12@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization

func main() {
	services := initServices()

	app := gin.Default()
	handlers.CreateHealthEndpoints(app)
	handlers.CreateStudentsEndpoints(app, services.studentRepository, services.authService)
	handlers.CreateAuthEndpoints(app, services.userRepository, services.authService)
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
