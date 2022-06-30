package main

import (
	"fmt"
	"net/http"
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
		panic(err)
	}

	studentRepository := repository.NewStudentsRepo(databaseService)

	student, err := studentRepository.FindStudent(2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", student)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	// listen and serve on
	r.Run(":8080")
}
