package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IDatabaseService interface {
}

type databaseService struct {
	db *gorm.DB
}

const connectionFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"

var ErrDatabaseConnection = errors.New("couldn't connect to database")

func NewDatabaseService(dbUser, dbPass, dbHost, dbPort, dbName string) (*databaseService, error) {
	database, err := gorm.Open(mysql.Open(fmt.Sprintf(connectionFormat, dbUser, dbPass, dbHost, dbPort, dbName)))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseConnection, err)
	}
	return &databaseService{db: database}, nil
}

type Student struct {
	ID   int
	Name string
	Age  int
}

func (s databaseService) GetStudents() []Student {
	var students []Student
	s.db.Raw("SELECT * FROM student").Scan(&students)
	return students
}
