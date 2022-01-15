package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brokeyourbike/lets-go-chat/api/server"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/brokeyourbike/lets-go-chat/mocks"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_users_HandleChat(t *testing.T) {
	cases := map[string]struct {
		token      string
		statusCode int
		message    string
		setupMock  func(activeUsersRepo *mocks.ActiveUsersRepo, tokensRepo *mocks.TokensRepo)
	}{
		"token can not be empty": {
			token:      "",
			statusCode: http.StatusBadRequest,
			message:    "Token format invalid\n",
			setupMock:  func(activeUsersRepo *mocks.ActiveUsersRepo, tokensRepo *mocks.TokensRepo) {},
		},
		"token must be valid uuid": {
			token:      "not-uuid",
			statusCode: http.StatusBadRequest,
			message:    "Token format invalid\n",
			setupMock:  func(activeUsersRepo *mocks.ActiveUsersRepo, tokensRepo *mocks.TokensRepo) {},
		},
		"token must exist": {
			token:      "c0834646-95ce-4d71-9cc3-e54ae187d1b9",
			statusCode: http.StatusBadRequest,
			message:    "Token invalid\n",
			setupMock: func(activeUsersRepo *mocks.ActiveUsersRepo, tokensRepo *mocks.TokensRepo) {
				id := uuid.MustParse("c0834646-95ce-4d71-9cc3-e54ae187d1b9")
				tokensRepo.On("Get", id).Return(models.Token{}, db.ErrTokenNotFound)
			},
		},
		"token should be returned from the query": {
			token:      "c0834646-95ce-4d71-9cc3-e54ae187d1b9",
			statusCode: http.StatusInternalServerError,
			message:    "Token cannot be validated\n",
			setupMock: func(activeUsersRepo *mocks.ActiveUsersRepo, tokensRepo *mocks.TokensRepo) {
				id := uuid.MustParse("c0834646-95ce-4d71-9cc3-e54ae187d1b9")
				tokensRepo.On("Get", id).Return(models.Token{}, errors.New("cannot quary token"))
			},
		},
		"token must not be expired": {
			token:      "c0834646-95ce-4d71-9cc3-e54ae187d1b9",
			statusCode: http.StatusBadRequest,
			message:    "Token expired\n",
			setupMock: func(activeUsersRepo *mocks.ActiveUsersRepo, tokensRepo *mocks.TokensRepo) {
				t := models.Token{
					ID:        uuid.MustParse("c0834646-95ce-4d71-9cc3-e54ae187d1b9"),
					UserID:    uuid.New(),
					ExpiresAt: time.Now().Add(-time.Minute),
				}
				tokensRepo.On("Get", t.ID).Return(t, nil)
			},
		},
		"can't upgrade request if it's not upgradable": {
			token:      "c0834646-95ce-4d71-9cc3-e54ae187d1b9",
			statusCode: http.StatusBadRequest,
			message:    "Bad Request\nCannot upgrade request to websocket protocol\n",
			setupMock: func(activeUsersRepo *mocks.ActiveUsersRepo, tokensRepo *mocks.TokensRepo) {
				t := models.Token{
					ID:        uuid.MustParse("c0834646-95ce-4d71-9cc3-e54ae187d1b9"),
					UserID:    uuid.New(),
					ExpiresAt: time.Now().Add(time.Minute),
				}
				tokensRepo.On("Get", t.ID).Return(t, nil)
				tokensRepo.On("InvalidateByUserId", t.UserID).Return(nil)
				activeUsersRepo.On("Add", t.UserID).Return(nil)
				activeUsersRepo.On("Delete", t.UserID).Return(nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			activeUsersRepo := new(mocks.ActiveUsersRepo)
			tokensRepo := new(mocks.TokensRepo)
			chat := NewChat(NewHub(), activeUsersRepo, tokensRepo, nil)

			c.setupMock(activeUsersRepo, tokensRepo)

			req := httptest.NewRequest(http.MethodGet, "/chat/ws.rtm.start", nil)
			w := httptest.NewRecorder()

			srv := server.NewServer(NewUsers(nil, nil, nil), chat)
			srv.WsRTMStart(w, req, server.WsRTMStartParams{Token: c.token})

			assert.Equal(t, c.statusCode, w.Result().StatusCode)
			assert.Equal(t, c.message, w.Body.String())

			activeUsersRepo.AssertExpectations(t)
			tokensRepo.AssertExpectations(t)
		})
	}
}
