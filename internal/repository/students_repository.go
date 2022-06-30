package repository

import (
	"github.com/darolpz/students/internal/database"
	"github.com/darolpz/students/internal/model"
)

type IStudentsRepository interface {
	FindStudent(id int) (model.Student, error)
}

type StudentsRepo struct {
	db database.IDatabaseService
}

func NewStudentsRepo(db database.IDatabaseService) IStudentsRepository {
	return StudentsRepo{db: db}
}

func (s StudentsRepo) FindStudent(id int) (model.Student, error) {
	return s.db.FindStudent(id)
}
