package repository

import (
	"errors"

	"github.com/darolpz/students/internal/database"
	"github.com/darolpz/students/internal/model"
)

var ErrUserNotFound = database.ErrUserNotFound

type IUsersRepository interface {
	FindUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
}

type usersRepository struct {
	db database.IDatabaseService
}

func NewUsersRepository(db database.IDatabaseService) IUsersRepository {
	return usersRepository{db: db}
}

func (u usersRepository) FindUserByEmail(email string) (model.User, error) {
	user, err := u.db.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return user, nil
}

func (u usersRepository) CreateUser(user model.User) (model.User, error) {
	return u.db.CreateUser(user)
}
