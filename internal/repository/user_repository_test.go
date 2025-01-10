package repository

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/nix-united/golang-gin-boilerplate/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testTime = time.Now()

func TestFindUserByEmail(t *testing.T) {
	db, mock, creatingMockErr := sqlmock.New()

	if creatingMockErr != nil {
		t.Fatalf(
			"an error '%s' occurred during creating a stub database connection",
			creatingMockErr,
		)
	}

	userID := uint(1)
	userEmail := "test@email.com"
	userPassword := "test pass"
	userFullName := "test full name"

	rows := mock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
		"email",
		"password",
		"full_name",
	}).
		AddRow(
			userID,
			testTime,
			testTime,
			nil,
			userEmail,
			userPassword,
			userFullName,
		)

	query := "SELECT * FROM `users`  WHERE email = ? AND `users`.`deleted_at` IS NULL"

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(mock.NewRows([]string{"VERSION()"}).AddRow("8.0.32"))
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs("test@email.com").
		WillReturnRows(rows)
	mockedDbConn, connErr := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if connErr != nil {
		t.Fatalf(
			"an error '%s' occurred during opening a stub database connection",
			connErr,
		)
	}

	want := model.User{
		Model: gorm.Model{
			ID:        userID,
			CreatedAt: testTime,
			UpdatedAt: testTime,
			DeletedAt: struct {
				Time  time.Time
				Valid bool
			}{Valid: false},
		},
		Email:    userEmail,
		Password: userPassword,
		FullName: userFullName,
	}

	got, err := NewUserRepository(mockedDbConn).FindUserByEmail(userEmail)

	assert.NoError(t, err)
	assert.Equal(t, want, got)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindUserByEmailReturnsError(t *testing.T) {
	db, mock, creatingMockErr := sqlmock.New()

	if creatingMockErr != nil {
		t.Fatalf(
			"an error '%s' occurred during creating a stub database connection",
			creatingMockErr,
		)
	}

	query := "SELECT * FROM `users`  WHERE email = ? AND `users`.`deleted_at` IS NULL"

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(mock.NewRows([]string{"VERSION()"}).AddRow("8.0.32"))
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnError(errors.New("test"))

	mockedDbConn, connErr := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if connErr != nil {
		t.Fatalf(
			"an error '%s' occurred during opening a stub database connection",
			connErr,
		)
	}

	want := model.User{}

	got, err := NewUserRepository(mockedDbConn).FindUserByEmail("test@email.com")

	assert.Error(t, err)
	assert.Equal(t, want, got)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindUserById(t *testing.T) {
	db, mock, creatingMockErr := sqlmock.New()

	if creatingMockErr != nil {
		t.Fatalf(
			"an error '%s' occurred during creating a stub database connection",
			creatingMockErr,
		)
	}

	userID := uint(1)
	userEmail := "test@email.com"
	userPassword := "test pass"
	userFullName := "test full name"

	rows := mock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
		"email",
		"password",
		"full_name",
	}).
		AddRow(
			userID,
			testTime,
			testTime,
			nil,
			userEmail,
			userPassword,
			userFullName,
		)

	query := "SELECT * FROM `users`  WHERE id = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1"

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(mock.NewRows([]string{"VERSION()"}).AddRow("8.0.32"))
	mock.
		ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID).
		WillReturnRows(rows)

	mockedDbConn, connErr := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if connErr != nil {
		t.Fatalf(
			"an error '%s' occurred during opening a stub database connection",
			connErr,
		)
	}

	want := model.User{
		Model: gorm.Model{
			ID:        userID,
			CreatedAt: testTime,
			UpdatedAt: testTime,
			DeletedAt: struct {
				Time  time.Time
				Valid bool
			}{Valid: false},
		},
		Email:    userEmail,
		Password: userPassword,
		FullName: userFullName,
	}

	got := NewUserRepository(mockedDbConn).FindUserByID(int(userID))

	assert.Equal(t, want, got)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUser(t *testing.T) {
	db, mock, creatingMockErr := sqlmock.New()

	if creatingMockErr != nil {
		t.Fatalf(
			"an error '%s' occurred during creating a stub database connection",
			creatingMockErr,
		)
	}

	userID := uint(1)
	userEmail := "test@email.com"
	userPassword := "test pass"
	userFullName := "test full name"

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(mock.NewRows([]string{"VERSION()"}).AddRow("8.0.32"))
	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(
			testTime,
			testTime,
			nil,
			userEmail,
			userPassword,
			userFullName,
			userID,
		).
		WillReturnResult(sqlmock.NewResult(int64(userID), 1))

	mock.ExpectCommit()

	mockedDbConn, connErr := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if connErr != nil {
		t.Fatalf(
			"an error '%s' occurred during opening a stub database connection",
			connErr,
		)
	}

	savingErr := NewUserRepository(mockedDbConn).StoreUser(model.User{
		Model: gorm.Model{
			ID:        userID,
			CreatedAt: testTime,
			UpdatedAt: testTime,
			DeletedAt: struct {
				Time  time.Time
				Valid bool
			}{Valid: false},
		},
		Email:    userEmail,
		Password: userPassword,
		FullName: userFullName,
	})

	assert.NoError(t, savingErr)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUserReturnsError(t *testing.T) {
	db, mock, creatingMockErr := sqlmock.New()

	if creatingMockErr != nil {
		t.Fatalf(
			"an error '%s' occurred during creating a stub database connection",
			creatingMockErr,
		)
	}

	userID := uint(1)
	userEmail := "test@email.com"
	userPassword := "test pass"
	userFullName := "test full name"

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(mock.NewRows([]string{"VERSION()"}).AddRow("8.0.32"))
	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
		WithArgs(
			testTime,
			testTime,
			nil,
			userEmail,
			userPassword,
			userFullName,
			userID,
		).
		WillReturnError(errors.New("test"))

	mock.ExpectRollback()

	mockedDbConn, connErr := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	if connErr != nil {
		t.Fatalf(
			"an error '%s' occurred during opening a stub database connection",
			connErr,
		)
	}

	savingErr := NewUserRepository(mockedDbConn).StoreUser(model.User{
		Model: gorm.Model{
			ID:        userID,
			CreatedAt: testTime,
			UpdatedAt: testTime,
			DeletedAt: struct {
				Time  time.Time
				Valid bool
			}{Valid: false},
		},
		Email:    userEmail,
		Password: userPassword,
		FullName: userFullName,
	})

	assert.Error(t, savingErr)

	assert.NoError(t, mock.ExpectationsWereMet())
}
