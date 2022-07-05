package handlers

import (
	"net/http"

	"github.com/darolpz/students/cmd/handlers/middleware"
	"github.com/darolpz/students/internal/auth"
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

func CreateStudentsEndpoints(
	app *gin.Engine,
	studentsRepo repository.IStudentsRepository,
	authService auth.IAuthService) {
	students := app.Group("students")
	students.Use(middleware.AuthMiddleware(authService))
	students.GET("/:id", FindStudent(studentsRepo))

	students.GET("/list", ListStudents(studentsRepo))

	students.POST("/", CreateStudent(studentsRepo))

	students.PATCH("/:id", UpdateStudent(studentsRepo))
}

func CreateAuthEndpoints(app *gin.Engine, usersRepo repository.IUsersRepository, authService auth.IAuthService) {
	auth := app.Group("auth")

	auth.POST("/login", Login(usersRepo, authService))
	auth.POST("/register", Register(usersRepo, authService))
}
