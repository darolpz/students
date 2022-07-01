package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

func craeteHealthEndpoints(app *gin.Engine) {
	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func createStudentsEndpoints(app *gin.Engine, repo repository.IStudentsRepository) {
	students := app.Group("students")

	students.GET("/:id", func(c *gin.Context) {
		studentID := c.Param("id")
		student, err := repo.FindStudent(studentID)
		if err != nil {
			if errors.Is(err, repository.ErrStudentNotFound) {
				log.Printf("couldnt find student with id %s: %s", studentID, err)
				c.String(http.StatusNotFound, err.Error())
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	})

	students.GET("/list", func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")
		students, err := repo.ListStudents(offset, limit)
		if err != nil {
			log.Printf("couldnt list student: %s", err)
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"students": students,
		})
	})
}
