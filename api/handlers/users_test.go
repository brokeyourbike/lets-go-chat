package handlers

import (
	"bytes"
	"encoding/json"
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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type userPayload struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func preparePayload(t *testing.T, p interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	require.NoError(t, err)
	return &buf
}

func Test_users_HandleUserCreate(t *testing.T) {
	cases := map[string]struct {
		payload    userPayload
		statusCode int
		message    string
		setupMock  func(usersRepo *mocks.UsersRepo)
	}{
		"username is required": {
			payload: userPayload{
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "UserName is too short\n",
			setupMock:  func(usersRepo *mocks.UsersRepo) {},
		},
		"password is required": {
			payload: userPayload{
				UserName: "john",
			},
			statusCode: http.StatusBadRequest,
			message:    "Password is too short\n",
			setupMock:  func(usersRepo *mocks.UsersRepo) {},
		},
		"username should be min 4 symbols": {
			payload: userPayload{
				UserName: "joh",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "UserName is too short\n",
			setupMock:  func(usersRepo *mocks.UsersRepo) {},
		},
		"password should be min 8 symbols": {
			payload: userPayload{
				UserName: "john",
				Password: "1234567",
			},
			statusCode: http.StatusBadRequest,
			message:    "Password is too short\n",
			setupMock:  func(usersRepo *mocks.UsersRepo) {},
		},
		"user will not be created if it's exists": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "User with userName john already exists\n",
			setupMock: func(usersRepo *mocks.UsersRepo) {
				usersRepo.On("GetByUserName", "john").Return(models.User{}, nil)
			},
		},
		"user will not be created if db return error": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusInternalServerError,
			message:    "User cannot be created\n",
			setupMock: func(usersRepo *mocks.UsersRepo) {
				usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("user not found"))
				usersRepo.On("Create", mock.AnythingOfType("User")).Return(models.User{}, errors.New("cannot create user"))
			},
		},
		"user can be created": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusOK,
			message:    "{\"id\":\"50b470c4-223f-443b-b525-26be5e005063\",\"userName\":\"john\"}\n",
			setupMock: func(usersRepo *mocks.UsersRepo) {
				u := models.User{ID: uuid.MustParse("50b470c4-223f-443b-b525-26be5e005063"), UserName: "john"}

				usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("user not found"))
				usersRepo.On("Create", mock.AnythingOfType("User")).Return(u, nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			usersRepo := new(mocks.UsersRepo)
			users := NewUsers(usersRepo, nil, nil)
			c.setupMock(usersRepo)

			req := httptest.NewRequest(http.MethodPost, "/user", preparePayload(t, c.payload))
			w := httptest.NewRecorder()

			srv := server.NewServer(users, NewChat(NewHub(), nil, nil, nil))
			srv.CreateUser(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)
			assert.Equal(t, c.message, w.Body.String())

			usersRepo.AssertExpectations(t)
		})
	}
}

func Test_users_HandleUserLogin(t *testing.T) {
	cases := map[string]struct {
		payload    userPayload
		statusCode int
		message    string
		setupMock  func(usersRepo *mocks.UsersRepo, tokensRepo *mocks.TokensRepo)
	}{
		"user can't login if it's not exists": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "User with userName john not found\n",
			setupMock: func(usersRepo *mocks.UsersRepo, tokensRepo *mocks.TokensRepo) {
				usersRepo.On("GetByUserName", "john").Return(models.User{}, db.ErrUserNotFound)
			},
		},
		"user can't login if db cannot perform query": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusInternalServerError,
			message:    "Cannot query user with userName john\n",
			setupMock: func(usersRepo *mocks.UsersRepo, tokensRepo *mocks.TokensRepo) {
				usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("cannot query user"))
			},
		},
		"user can't login if password is not valid": {
			payload: userPayload{
				UserName: "john",
				Password: "not-a-valid-password",
			},
			statusCode: http.StatusBadRequest,
			message:    "Invalid password\n",
			setupMock: func(usersRepo *mocks.UsersRepo, tokensRepo *mocks.TokensRepo) {
				u := models.User{ID: uuid.New(), UserName: "john", PasswordHash: "$2a$04$hOGVri5G8ZLnrXtO/EP6keZkdzveoVGfh9krMXxeI/OP2QcSDJWOW"}
				usersRepo.On("GetByUserName", "john").Return(u, nil)
			},
		},
		"user can't login if token was not created": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusInternalServerError,
			message:    "Cannot create token for user\n",
			setupMock: func(usersRepo *mocks.UsersRepo, tokensRepo *mocks.TokensRepo) {
				u := models.User{ID: uuid.New(), UserName: "john", PasswordHash: "$2a$04$hOGVri5G8ZLnrXtO/EP6keZkdzveoVGfh9krMXxeI/OP2QcSDJWOW"}

				usersRepo.On("GetByUserName", "john").Return(u, nil)
				tokensRepo.On("Create", mock.AnythingOfType("Token")).Return(models.Token{}, errors.New("cannot create token"))
			},
		},
		"user can login": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusOK,
			message:    "{\"url\":\"ws://example.com/chat/ws.rtm.start?token=d0b61534-256c-4ef1-b98e-53c5424ce7cd\"}\n",
			setupMock: func(usersRepo *mocks.UsersRepo, tokensRepo *mocks.TokensRepo) {
				u := models.User{ID: uuid.New(), UserName: "john", PasswordHash: "$2a$04$hOGVri5G8ZLnrXtO/EP6keZkdzveoVGfh9krMXxeI/OP2QcSDJWOW"}
				t := models.Token{ID: uuid.MustParse("d0b61534-256c-4ef1-b98e-53c5424ce7cd"), UserID: u.ID, ExpiresAt: time.Now().Add(time.Minute)}

				usersRepo.On("GetByUserName", "john").Return(u, nil)
				tokensRepo.On("Create", mock.AnythingOfType("Token")).Return(t, nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			usersRepo := new(mocks.UsersRepo)
			tokensRepo := new(mocks.TokensRepo)
			users := NewUsers(usersRepo, nil, tokensRepo)
			c.setupMock(usersRepo, tokensRepo)

			req := httptest.NewRequest(http.MethodPost, "/user/login", preparePayload(t, c.payload))
			w := httptest.NewRecorder()

			srv := server.NewServer(users, NewChat(NewHub(), nil, nil, nil))
			srv.LoginUser(w, req)

			assert.Equal(t, c.statusCode, w.Result().StatusCode)
			assert.Equal(t, c.message, w.Body.String())

			usersRepo.AssertExpectations(t)
			tokensRepo.AssertExpectations(t)
		})
	}
}

func Test_users_HandleUserActive(t *testing.T) {
	activeUsersRepo := new(mocks.ActiveUsersRepo)
	users := NewUsers(nil, activeUsersRepo, nil)

	activeUsersRepo.On("Count").Return(10)

	req := httptest.NewRequest(http.MethodGet, "/user/active", nil)
	w := httptest.NewRecorder()

	srv := server.NewServer(users, NewChat(NewHub(), nil, nil, nil))
	srv.GetActiveUsers(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, "{\"count\":10}\n", w.Body.String())

	activeUsersRepo.AssertExpectations(t)
}
