package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/darolpz/students/internal/model"
	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

func FindStudent(repo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func ListStudents(repo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}

func CreateStudent(repo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		newStudent := model.Student{}
		if err := c.BindJSON(&newStudent); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		student, err := repo.CreateStudent(newStudent)
		if err != nil {
			log.Printf("couldnt create student: %s", err)
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}

func UpdateStudent(repo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		studentID := c.Param("id")
		newStudent := model.Student{}
		if err := c.BindJSON(&newStudent); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		student, err := repo.UpdateStudent(studentID, newStudent)
		if err != nil {
			log.Printf("couldnt create student: %s", err)
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
