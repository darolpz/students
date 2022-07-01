package database

import (
	"errors"
	"fmt"

	"github.com/darolpz/students/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IDatabaseService interface {
	FindStudent(id string) (model.Student, error)
	ListStudents(offset, limit string) ([]model.Student, error)
	CreateStudent(student model.Student) (model.Student, error)
	UpdateStudent(id string, student model.Student) (model.Student, error)
}

type databaseService struct {
	db *gorm.DB
}

const connectionFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"

var (
	ErrDatabaseConnection = errors.New("couldn't connect to database")
	ErrFindStudent        = errors.New("couldn't find student")
	ErrListStudents       = errors.New("couldn't list students")
	ErrStudentNotFound    = errors.New("student not found")
	ErrCreateStudent      = errors.New("couldn't create student")
	ErrUpdateStudent      = errors.New("couldn't update student")
)

func NewDatabaseService(dbUser, dbPass, dbHost, dbPort, dbName string) (*databaseService, error) {
	database, err := gorm.Open(mysql.Open(fmt.Sprintf(connectionFormat, dbUser, dbPass, dbHost, dbPort, dbName)))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseConnection, err)
	}
	return &databaseService{db: database}, nil
}

func (s databaseService) FindStudent(id string) (model.Student, error) {
	var student model.Student
	if err := s.db.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return student, fmt.Errorf("%w: %s", ErrStudentNotFound, err)
		}
		return student, fmt.Errorf("%w: %s", ErrFindStudent, err)
	}

	return student, nil
}

func (s databaseService) ListStudents(limit, offset string) ([]model.Student, error) {
	var students []model.Student
	if err := s.db.Find(&students).Error; err != nil {
		return students, fmt.Errorf("%w: %s", ErrListStudents, err)
	}
	return students, nil
}

func (s databaseService) CreateStudent(student model.Student) (model.Student, error) {
	if err := s.db.Create(&student).Error; err != nil {
		return student, fmt.Errorf("%w: %s", ErrCreateStudent, err)
	}
	return student, nil
}

func (s databaseService) UpdateStudent(id string, newStudent model.Student) (model.Student, error) {
	var student model.Student
	if err := s.db.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return student, fmt.Errorf("%w: %s", ErrStudentNotFound, err)
		}
		return student, fmt.Errorf("%w: %s", ErrFindStudent, err)
	}

	student.FirstName = newStudent.FirstName
	student.LastName = newStudent.LastName
	student.Email = newStudent.Email
	student.Age = newStudent.Age

	if err := s.db.Save(&student).Error; err != nil {
		return student, fmt.Errorf("%w: %s", ErrUpdateStudent, err)
	}
	return student, nil
}
