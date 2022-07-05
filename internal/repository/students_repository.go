package repository

import (
	"errors"

	"github.com/darolpz/students/internal/database"
	"github.com/darolpz/students/internal/model"
)

var ErrStudentNotFound = database.ErrStudentNotFound

type IStudentsRepository interface {
	FindStudent(id string) (model.Student, error)
	ListStudents(offset, limit string) ([]model.Student, error)
	CreateStudent(student model.Student) (model.Student, error)
	UpdateStudent(id string, student model.Student) (model.Student, error)
}

type studentsRepo struct {
	db database.IDatabaseService
}

func NewStudentsRepo(db database.IDatabaseService) IStudentsRepository {
	return studentsRepo{db: db}
}

func (s studentsRepo) FindStudent(id string) (model.Student, error) {
	student, err := s.db.FindStudent(id)
	if err != nil {
		if errors.Is(err, database.ErrStudentNotFound) {
			return model.Student{}, ErrStudentNotFound
		}
		return model.Student{}, err
	}
	return student, nil
}

func (s studentsRepo) ListStudents(offset, limit string) ([]model.Student, error) {
	return s.db.ListStudents(offset, limit)
}

func (s studentsRepo) CreateStudent(student model.Student) (model.Student, error) {
	return s.db.CreateStudent(student)
}

func (s studentsRepo) UpdateStudent(id string, student model.Student) (model.Student, error) {
	return s.db.UpdateStudent(id, student)
}
