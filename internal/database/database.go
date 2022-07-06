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
	FindUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	DeleteStudent(id string) error
}

type databaseService struct {
	db *gorm.DB
}

const connectionFormat = "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local"

var (
	ErrDatabaseConnection = errors.New("couldn't connect to database")
	ErrFindStudent        = errors.New("couldn't find student")
	ErrStudentNotFound    = errors.New("student not found")
	ErrListStudents       = errors.New("couldn't list students")
	ErrCreateStudent      = errors.New("couldn't create student")
	ErrUpdateStudent      = errors.New("couldn't update student")
	ErrUserNotFound       = errors.New("user not found")
	ErrFindUser           = errors.New("couldn't find user")
	ErrCreateUser         = errors.New("couldn't create user")
	ErrDeleteStudent      = errors.New("couldn't delete student")
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

func (s databaseService) FindUserByEmail(email string) (model.User, error) {
	var user model.User
	if err := s.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, fmt.Errorf("%w: %s", ErrUserNotFound, err)
		}
		return user, fmt.Errorf("%w: %s", ErrFindUser, err)
	}
	return user, nil
}

func (s databaseService) CreateUser(user model.User) (model.User, error) {
	if err := s.db.Create(&user).Error; err != nil {
		return user, fmt.Errorf("%w: %s", ErrCreateUser, err)
	}
	return user, nil
}

func (s databaseService) DeleteStudent(id string) error {
	var student model.Student
	if err := s.db.First(&student, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %s", ErrStudentNotFound, err)
		}
		return fmt.Errorf("%w: %s", ErrFindStudent, err)
	}

	if err := s.db.Delete(&student).Error; err != nil {
		return fmt.Errorf("%w: %s", ErrDeleteStudent, err)
	}
	return nil
}
