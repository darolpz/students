package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/darolpz/students/internal/model"
	"github.com/darolpz/students/internal/repository"
	"github.com/gin-gonic/gin"
)

// FindStudent godoc
// @Summary      Find Student
// @Description  get student by id
// @Tags         students
// @Param        student_id  path string  true  "student_id"  1
// @Success      200 {object} model.Student
// @Failure      404  {string} string
// @Failure      500 {string} string
// @Router       /students/{student_id} [get]
// @Security Authorization
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

// ListStudent godoc
// @Summary      List Student
// @Description  returns all students
// @Tags         students
// @Param        offset    query     string  false  "list offset"  0
// @Param        limit    query     string  false  "list limit"  0
// @Success      200 {array} model.Student
// @Failure      404  {string} string
// @Failure      500 {string} string
// @Router       /students/list [get]
// @Security Authorization
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

// CreateStudent godoc
// @Summary      Create student
// @Description  creates a new student
// @Tags         students
// @Param        student body model.Student true "user"
// @Accept       json
// @Produce      json
// @Success      200 {object} model.Student
// @Failure      500 {string} string
// @Router       /students [post]
// @Security Authorization
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

// UpdateStudent godoc
// @Summary      Update student
// @Description  modify student data
// @Tags         students
// @Param        student_id  path string  true  "student_id"  1
// @Param        student body model.Student true "user"
// @Accept       json
// @Produce      json
// @Success      200 {object} model.Student
// @Failure      400 {string} string "bad request"
// @Failure      500 {string} string
// @Router       /students/{student_id} [patch]
// @Security Authorization
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
