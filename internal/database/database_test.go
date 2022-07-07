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

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
			AddRow("1", "John", "Doe", "john.doe@gmail.com", "33"))

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
		WithArgs("5").
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `students` WHERE `students`.`id` = ? ORDER BY `students`.`id` LIMIT 1")).
		WithArgs("x").
		WillReturnError(errors.New("error"))

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
		want          model.Student
		expectedError error
	}{
		{
			name:       "should_return_student",
			student_id: "1",
			want: model.Student{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@gmail.com",
				Age:       33,
			},
		},
		{
			name:          "should_return_error_student_not_found",
			student_id:    "5",
			expectedError: ErrStudentNotFound,
		},
		{
			name:          "should_return_error_student_not_found",
			student_id:    "5",
			expectedError: ErrFindStudent,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `students` LIMIT 10")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "email", "age"}).
			AddRow("1", "John", "Doe", "john.doe@gmail.com", "33").
			AddRow("2", "Dario", "Lopez", "daropl12@gmail.com", "26"))

	mock.ExpectQuery(
		regexp.QuoteMeta("SELECT * FROM `students` LIMIT 25 OFFSET 9")).
		WillReturnError(errors.New("error"))

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
		want          []model.Student
		expectedError error
	}{
		{
			name:   "should_return_students",
			limit:  10,
			offset: 0,
			want: []model.Student{
				{ID: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@gmail.com", Age: 33},
				{ID: 2, FirstName: "Dario", LastName: "Lopez", Email: "daropl12@gmail.com", Age: 26},
			},
		},
		{
			name:          "should_return_error_list_students",
			limit:         25,
			offset:        9,
			expectedError: ErrListStudents,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
