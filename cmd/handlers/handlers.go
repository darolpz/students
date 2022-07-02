package handlers

import (
	"net/http"

	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

func CreateHealthEndpoints(app *gin.Engine) {
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func CreateStudentsEndpoints(app *gin.Engine, repo repository.IStudentsRepository) {
	students := app.Group("students")

	students.GET("/:id", FindStudent(repo))

	students.GET("/list", ListStudents(repo))

	students.POST("/", CreateStudent(repo))

	students.PATCH("/:id", UpdateStudent(repo))
}
