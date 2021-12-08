package db

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TokensSuite struct {
	Suite
	repository *TokensRepo
}

func (s *TokensSuite) SetupTest() {
	s.SetupDatabase()
	s.repository = NewTokensRepo(s.db)
}

func TestTokens(t *testing.T) {
	suite.Run(t, new(TokensSuite))
}

func (s *TokensSuite) Test_repository_Create_ItCanCreateToken() {
	token := models.Token{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: time.Now()}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "tokens" ("id","user_id","expires_at") VALUES ($1,$2,$3)`)).
		WithArgs(token.ID, token.UserID, token.ExpiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.Create(token)
	require.NoError(s.T(), err)
}

func (s *TokensSuite) Test_repository_Get_ItCanReturnTokenById() {
	token := models.Token{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: time.Now()}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tokens" WHERE id = $1 ORDER BY "tokens"."id" LIMIT 1`)).
		WithArgs(token.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "expires_at"}).
			AddRow(token.ID, token.UserID, token.ExpiresAt))

	res, err := s.repository.Get(token.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), token, res)
}

func (s *TokensSuite) Test_repository_Get_ItCanReturnErrTokenNotFound() {
	tokenID := uuid.New()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tokens" WHERE id = $1 ORDER BY "tokens"."id" LIMIT 1`)).
		WithArgs(tokenID).
		WillReturnError(gorm.ErrRecordNotFound)

	res, err := s.repository.Get(tokenID)
	require.ErrorIs(s.T(), err, ErrTokenNotFound)
	require.Equal(s.T(), models.Token{}, res)
}

func (s *TokensSuite) Test_repository_Get_ItCanReturnGeneralError() {
	tokenID := uuid.New()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "tokens" WHERE id = $1 ORDER BY "tokens"."id" LIMIT 1`)).
		WithArgs(tokenID).
		WillReturnError(gorm.ErrInvalidField)

	res, err := s.repository.Get(tokenID)
	require.ErrorIs(s.T(), err, gorm.ErrInvalidField)
	require.Equal(s.T(), models.Token{}, res)
}

func (s *TokensSuite) Test_repository_InvalidateByUserId() {
	userID := uuid.New()

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "tokens" WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	err := s.repository.InvalidateByUserId(userID)
	require.NoError(s.T(), err)
}
