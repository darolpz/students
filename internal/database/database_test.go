package database

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/darolpz/students/internal/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_databaseService_FindStudent(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		DriverName:                "mysql",
		DSN:                       "sqlmock_db",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	s := databaseService{
		db: db,
	}

	tests := []struct {
		name          string
		student_id    string
		setMock       func(mock sqlmock.Sqlmock)
		want          model.Student
		expectedError error
	}{
		{
			name:       "should_find_student",
			student_id: "1",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
						AddRow("1", "John", "Doe", "john.doe@gmail.com", "33"))
			},
			want: model.Student{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       33,
			},
		},
		{
			name:       "should_return_error_student_not_found",
			student_id: "5",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("5").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: ErrStudentNotFound,
		},
		{
			name:       "should_return_error_find_student",
			student_id: "5",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("x").
					WillReturnError(errors.New("error"))
			},
			expectedError: ErrFindStudent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setMock(mock)
			got, err := s.FindStudent(tt.student_id)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_databaseService_ListStudents(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		DriverName:                "mysql",
		DSN:                       "sqlmock_db",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	s := databaseService{
		db: db,
	}

	tests := []struct {
		name          string
		limit         int
		offset        int
		setMock       func(mock sqlmock.Sqlmock)
		want          []model.Student
		expectedError error
	}{
		{
			name:   "should_return_students_list",
			limit:  10,
			offset: 0,
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` LIMIT 10")).
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
						AddRow("1", "John", "Doe", "john.doe@gmail.com", "33").
						AddRow("2", "Dario", "Lopez", "daropl12@gmail.com", "26"))
			},
			want: []model.Student{
				{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@gmail.com", Age: 33},
				{ID: 2, FirstName: "Dario", LastName: "Lopez", Email: "daropl12@gmail.com", Age: 26},
			},
		},
		{
			name:   "should_return_error_list_students",
			limit:  25,
			offset: 9,
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` LIMIT 25 OFFSET 9")).
					WillReturnError(errors.New("error"))
			},
			expectedError: ErrListStudents,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setMock(mock)
			got, err := s.ListStudents(tt.limit, tt.offset)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_databaseService_CreateStudent(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		DriverName:                "mysql",
		DSN:                       "sqlmock_db",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	s := databaseService{
		db: db,
	}

	tests := []struct {
		name          string
		student       model.Student
		want          model.Student
		setMock       func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "should_create_student",
			student: model.Student{
				FirstName: "dario",
				LastName:  "lopez",
				Age:       26,
				Email:     "daropl12@gmail.com",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `students`").
					WithArgs("dario", "lopez", 26, "daropl12@gmail.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			want: model.Student{
				ID:        1,
				FirstName: "dario",
				LastName:  "lopez",
				Age:       26,
				Email:     "daropl12@gmail.com",
			},
		},
		{
			name: "should_return_error_create_student",
			student: model.Student{
				FirstName: "john",
				LastName:  "doe",
				Age:       33,
				Email:     "john.doe@gmail.com",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `students`").
					WithArgs("john", "doe", 33, "john.doe@gmail.com").
					WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
			expectedError: ErrCreateStudent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setMock(mock)
			got, err := s.CreateStudent(tt.student)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_databaseService_UpdateStudent(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		DriverName:                "mysql",
		DSN:                       "sqlmock_db",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	s := databaseService{
		db: db,
	}

	tests := []struct {
		name          string
		student_id    string
		newStudent    model.Student
		setMock       func(mock sqlmock.Sqlmock)
		want          model.Student
		expectedError error
	}{
		{
			name:       "should_update_student",
			student_id: "1",
			newStudent: model.Student{
				FirstName: "Dario",
				LastName:  "Lopez",
				Age:       26,
				Email:     "daropl12@gmail.com",
			},
			want: model.Student{
				ID:        1,
				FirstName: "Dario",
				LastName:  "Lopez",
				Age:       26,
				Email:     "daropl12@gmail.com",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
						AddRow("1", "John", "Doe", "john.doe@gmail.com", "33"))

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `students` SET `first_name`=?,`last_name`=?,`age`=?,`email`=? WHERE `id` = ?")).
					WithArgs("Dario", "Lopez", 26, "daropl12@gmail.com", 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:       "should_return_error_student_not_found",
			student_id: "99",
			newStudent: model.Student{
				FirstName: "Dario",
				LastName:  "Lopez",
				Age:       26,
				Email:     "daropl12@gmail.com",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("99").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: ErrStudentNotFound,
		},
		{
			name:       "should_return_error_find_student",
			student_id: "00",
			newStudent: model.Student{
				FirstName: "Dario",
				LastName:  "Lopez",
				Age:       26,
				Email:     "daropl12@gmail.com",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("00").
					WillReturnError(errors.New("error"))
			},
			expectedError: ErrFindStudent,
		},
		{
			name:       "should_return_error_update_student",
			student_id: "2",
			newStudent: model.Student{
				FirstName: "Dario",
				LastName:  "Lopez",
				Age:       26,
				Email:     "daropl12@gmail.com",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("2").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
						AddRow("1", "John", "Doe", "john.doe@gmail.com", "33"))

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE `students` SET `first_name`=?,`last_name`=?,`age`=?,`email`=? WHERE `id` = ?")).
					WithArgs("Dario", "Lopez", 26, "daropl12@gmail.com", 1).
					WillReturnError(errors.New("error"))
				mock.ExpectRollback()
			},
			expectedError: ErrUpdateStudent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setMock(mock)
			got, err := s.UpdateStudent(tt.student_id, tt.newStudent)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_databaseService_FindUserByEmail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		DriverName:                "mysql",
		DSN:                       "sqlmock_db",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	s := databaseService{
		db: db,
	}
	tests := []struct {
		name          string
		email         string
		setMock       func(mock sqlmock.Sqlmock)
		want          model.User
		expectedError error
	}{
		{
			name:  "should_return_user_by_email",
			email: "john.doe@gmail.com",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs("john.doe@gmail.com").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "role"}).
						AddRow("1", "John", "john.doe@gmail.com", "123456", "admin"))
			},
			want: model.User{
				ID:       1,
				Name:     "John",
				Email:    "john.doe@gmail.com",
				Password: "123456",
				Role:     "admin",
			},
		},
		{
			name:  "should_return_error_user_not_found",
			email: "john.doe@gmail.com",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs("john.doe@gmail.com").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: ErrUserNotFound,
		},
		{
			name:  "should_return_error_user_not_found",
			email: "john.doe@gmail.com",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs("john.doe@gmail.com").
					WillReturnError(errors.New("somer error"))
			},
			expectedError: ErrFindUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setMock(mock)
			got, err := s.FindUserByEmail(tt.email)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_databaseService_CreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		DriverName:                "mysql",
		DSN:                       "sqlmock_db",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	s := databaseService{
		db: db,
	}

	tests := []struct {
		name          string
		user          model.User
		setMock       func(mock sqlmock.Sqlmock)
		want          model.User
		expectedError error
	}{
		{
			name: "should_create_user",
			user: model.User{
				Name:     "John",
				Email:    "john.doe@gmail.com",
				Password: "123456",
				Role:     "admin",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs("John", "john.doe@gmail.com", "123456", "admin").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			want: model.User{
				ID:       1,
				Name:     "John",
				Email:    "john.doe@gmail.com",
				Password: "123456",
				Role:     "admin",
			},
		},
		{
			name: "should_return_error_create_user",
			user: model.User{
				Name:     "John",
				Email:    "john.doe@gmail.com",
				Password: "123456",
				Role:     "admin",
			},
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").
					WithArgs("John", "john.doe@gmail.com", "123456", "admin").
					WillReturnError(errors.New("somer error"))
				mock.ExpectRollback()
			},
			expectedError: ErrCreateUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setMock(mock)
			got, err := s.CreateUser(tt.user)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_databaseService_DeleteStudent(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		DriverName:                "mysql",
		DSN:                       "sqlmock_db",
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector)
	require.NoError(t, err)
	s := databaseService{
		db: db,
	}

	tests := []struct {
		name          string
		student_id    string
		setMock       func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:       "should_delete_student",
			student_id: "1",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
						AddRow("1", "John", "Doe", "john.doe@gmail.com", "33"))
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `students`").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:       "should_return_error_student_not_found",
			student_id: "1",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("1").
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError: ErrStudentNotFound,
		},
		{
			name:       "should_return_error_find_student",
			student_id: "1",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("1").
					WillReturnError(errors.New("somer error"))
			},
			expectedError: ErrFindStudent,
		},
		{
			name:       "should_return_error_delete_student",
			student_id: "1",
			setMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
						AddRow("1", "John", "Doe", "john.doe@gmail.com", "33"))
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `students`").
					WithArgs(1).
					WillReturnError(errors.New("somer error"))
				mock.ExpectRollback()
			},
			expectedError: ErrDeleteStudent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setMock(mock)
			err := s.DeleteStudent(tt.student_id)
			if tt.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tt.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
