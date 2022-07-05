package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/darolpz/students/internal/model"
	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

func FindStudent(studentsRepo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get query params
		studentID := c.Param("id")
		// Retrieve student from repository
		student, err := studentsRepo.FindStudent(studentID)
		if err != nil {
			// Check if error was student not found
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

func ListStudents(studentsRepo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get query params
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		// Retrieve students from repository
		students, err := studentsRepo.ListStudents(offset, limit)
		if err != nil {
			log.Printf("couldnt list student: %s", err)
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"students": students,
		})
	}
}

func CreateStudent(studentsRepo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		newStudent := model.Student{}
		// Bind the JSON body to the newStudent struct
		if err := c.BindJSON(&newStudent); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		// Persist the new student to repository
		student, err := studentsRepo.CreateStudent(newStudent)
		if err != nil {
			log.Printf("couldnt create student: %s", err)
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}

func UpdateStudent(studentsRepo repository.IStudentsRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		// Get query params
		studentID := c.Param("id")
		newStudent := model.Student{}
		// Bind the JSON body to the newStudent struct
		if err := c.BindJSON(&newStudent); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		// Update the student in the repository
		student, err := studentsRepo.UpdateStudent(studentID, newStudent)
		if err != nil {
			log.Printf("couldnt create student: %s", err)
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"student": student,
		})
	}
}
