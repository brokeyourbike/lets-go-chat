package db

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UsersSuite struct {
	Suite
	repository *UsersRepo
}

func (s *UsersSuite) SetupTest() {
	s.SetupDatabase()
	s.repository = NewUsersRepo(s.db)
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}

func (s *UsersSuite) Test_repository_Create_ItCanCreateUser() {
	user := models.User{ID: uuid.New(), UserName: "test", PasswordHash: "super-hash"}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "users" ("id","user_name","password_hash") VALUES ($1,$2,$3)`)).
		WithArgs(user.ID, user.UserName, user.PasswordHash).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	u, err := s.repository.Create(user)
	require.NoError(s.T(), err)
	require.Equal(s.T(), user, u)
}

func (s *UsersSuite) Test_repository_GetByUserName_ItCanReturnUserByUserName() {
	user := models.User{ID: uuid.New(), UserName: "test", PasswordHash: "super-hash"}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE user_name = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(user.UserName).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "password_hash"}).
			AddRow(user.ID, user.UserName, user.PasswordHash))

	res, err := s.repository.GetByUserName(user.UserName)
	require.NoError(s.T(), err)
	require.Equal(s.T(), user, res)
}

func (s *UsersSuite) Test_repository_GetByUserName_ItCanReturnErrUserNotFound() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE user_name = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs("johndoe").
		WillReturnError(gorm.ErrRecordNotFound)

	res, err := s.repository.GetByUserName("johndoe")
	require.ErrorIs(s.T(), err, ErrUserNotFound)
	require.Equal(s.T(), models.User{}, res)
}

func (s *UsersSuite) Test_repository_GetByUserName_ItCanReturnGeneralError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE user_name = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs("johndoe").
		WillReturnError(gorm.ErrInvalidField)

	res, err := s.repository.GetByUserName("johndoe")
	require.ErrorIs(s.T(), err, gorm.ErrInvalidField)
	require.Equal(s.T(), models.User{}, res)
}
