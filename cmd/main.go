package main

import (
	"fmt"
	"log"
	"os"

	"github.com/darolpz/students/internal/database"
	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	databaseService, err := database.NewDatabaseService(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatal(err)
	}

	studentRepository := repository.NewStudentsRepo(databaseService)

	app := gin.Default()
	craeteHealthEndpoints(app)
	createStudentsEndpoints(app, studentRepository)

	// listen and serve on
	app.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
