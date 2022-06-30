package database

import (
	"errors"
	"fmt"

	"github.com/darolpz/students/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IDatabaseService interface {
	FindStudent(id int) (model.Student, error)
}

type databaseService struct {
	db *gorm.DB
}

const connectionFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"

var (
	ErrDatabaseConnection = errors.New("couldn't connect to database")
	ErrFindStudent        = errors.New("couldn't find student")
	ErrListStudents       = errors.New("couldn't list students")
)

func NewDatabaseService(dbUser, dbPass, dbHost, dbPort, dbName string) (*databaseService, error) {
	database, err := gorm.Open(mysql.Open(fmt.Sprintf(connectionFormat, dbUser, dbPass, dbHost, dbPort, dbName)))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseConnection, err)
	}
	return &databaseService{db: database}, nil
}

func (s databaseService) FindStudent(id int) (model.Student, error) {
	var student model.Student
	if err := s.db.Raw("SELECT * FROM student WHERE ID = ?", id).Scan(&student).Error; err != nil {
		return student, fmt.Errorf("%w: %s", ErrFindStudent, err)
	}
	return student, nil
}

func (s databaseService) ListStudent() ([]model.Student, error) {
	var students []model.Student
	if err := s.db.Raw("SELECT * FROM student").Scan(&students).Error; err != nil {
		return students, fmt.Errorf("%w: %s", ErrListStudents, err)
	}
	return students, nil
}
