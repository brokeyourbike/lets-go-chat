package db

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MessagesSuite struct {
	Suite
	repository *MessagesRepo
}

func (s *MessagesSuite) SetupTest() {
	s.SetupDatabase()
	s.repository = NewMessagesRepo(s.db)
}

func TestMessages(t *testing.T) {
	suite.Run(t, new(MessagesSuite))
}

func (s *MessagesSuite) Test_repository_Create_ItCanCreateMessage() {
	token := models.Message{ID: uuid.New(), UserID: uuid.New(), Text: "", CreatedAt: time.Now()}

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO "messages" ("id","user_id","text","created_at") VALUES ($1,$2,$3,$4)`)).
		WithArgs(token.ID, token.UserID, token.Text, token.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	t, err := s.repository.Create(token)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), token, t)
}
