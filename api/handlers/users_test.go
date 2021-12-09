package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brokeyourbike/lets-go-chat/api/server"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/brokeyourbike/lets-go-chat/mocks"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	usersRepo := new(mocks.UsersRepo)
	users := NewUsers(usersRepo, nil, nil)

	cases := map[string]struct {
		payload    userPayload
		statusCode int
		message    string
		setupMock  func()
	}{
		"username is required": {
			payload: userPayload{
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
			setupMock:  func() {},
		},
		"password is required": {
			payload: userPayload{
				UserName: "john",
			},
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
			setupMock:  func() {},
		},
		"username should be min 4 symbols": {
			payload: userPayload{
				UserName: "joh",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
			setupMock:  func() {},
		},
		"password should be min 8 symbols": {
			payload: userPayload{
				UserName: "john",
				Password: "1234567",
			},
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
			setupMock:  func() {},
		},
		"user will not be created if it's exists": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "User with userName john already exists\n",
			setupMock: func() {
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
			setupMock: func() {
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
			setupMock: func() {
				u := models.User{ID: uuid.MustParse("50b470c4-223f-443b-b525-26be5e005063"), UserName: "john"}

				usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("user not found"))
				usersRepo.On("Create", mock.AnythingOfType("User")).Return(u, nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/v1/user", preparePayload(t, c.payload))
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewRouter())
			srv.Routes(users)
			srv.ServeHTTP(w, req)

			require.Equal(t, c.statusCode, w.Result().StatusCode)
			require.Equal(t, c.message, w.Body.String())

			usersRepo.AssertExpectations(t)
		})
	}
}

func Test_users_HandleUserLogin(t *testing.T) {
	usersRepo := new(mocks.UsersRepo)
	tokensRepo := new(mocks.TokensRepo)
	users := NewUsers(usersRepo, nil, tokensRepo)

	cases := map[string]struct {
		payload    userPayload
		statusCode int
		message    string
		setupMock  func()
	}{
		"user can't login if it's not exists": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "User with userName john not found\n",
			setupMock: func() {
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
			setupMock: func() {
				usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("cannot query user"))
			},
		},
		"user can login": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusOK,
			message:    "{\"url\":\"ws://example.com/v1/chat/ws.rtm.start?token=d0b61534-256c-4ef1-b98e-53c5424ce7cd\"}\n",
			setupMock: func() {
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
			c.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/v1/user/login", preparePayload(t, c.payload))
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewRouter())
			srv.Routes(users)
			srv.ServeHTTP(w, req)

			require.Equal(t, c.statusCode, w.Result().StatusCode)
			require.Equal(t, c.message, w.Body.String())

			usersRepo.AssertExpectations(t)
			tokensRepo.AssertExpectations(t)
		})
	}
}

func Test_users_HandleUserActive(t *testing.T) {
	activeUsersRepo := new(mocks.ActiveUsersRepo)
	users := NewUsers(nil, activeUsersRepo, nil)

	activeUsersRepo.On("Count").Return(10)

	req := httptest.NewRequest(http.MethodGet, "/v1/user/active", nil)
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(users)
	srv.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.Equal(t, "{\"count\":10}\n", w.Body.String())

	activeUsersRepo.AssertExpectations(t)
}

func Test_users_HandleChat(t *testing.T) {
	activeUsersRepo := new(mocks.ActiveUsersRepo)
	tokensRepo := new(mocks.TokensRepo)
	users := NewUsers(nil, activeUsersRepo, tokensRepo)

	cases := map[string]struct {
		token      string
		statusCode int
		message    string
		setupMock  func()
	}{
		"token can not be empty": {
			token:      "",
			statusCode: http.StatusBadRequest,
			message:    "Token format invalid\n",
			setupMock:  func() {},
		},
		"token must be valid uuid": {
			token:      "not-uuid",
			statusCode: http.StatusBadRequest,
			message:    "Token format invalid\n",
			setupMock:  func() {},
		},
		"token must exist": {
			token:      "c0834646-95ce-4d71-9cc3-e54ae187d1b9",
			statusCode: http.StatusBadRequest,
			message:    "Token invalid\n",
			setupMock: func() {
				id := uuid.MustParse("c0834646-95ce-4d71-9cc3-e54ae187d1b9")
				tokensRepo.On("Get", id).Return(models.Token{}, db.ErrTokenNotFound)
			},
		},
		"token should be returned from the query": {
			token:      "c0834646-95ce-4d71-9cc3-e54ae187d1b9",
			statusCode: http.StatusInternalServerError,
			message:    "Token cannot be validated\n",
			setupMock: func() {
				id := uuid.MustParse("c0834646-95ce-4d71-9cc3-e54ae187d1b9")
				tokensRepo.On("Get", id).Return(models.Token{}, errors.New("cannot quary token"))
			},
		},
		"token must not be expired": {
			token:      "c0834646-95ce-4d71-9cc3-e54ae187d1b9",
			statusCode: http.StatusBadRequest,
			message:    "Token expired\n",
			setupMock: func() {
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
			setupMock: func() {
				t := models.Token{
					ID:        uuid.MustParse("c0834646-95ce-4d71-9cc3-e54ae187d1b9"),
					UserID:    uuid.New(),
					ExpiresAt: time.Now().Add(time.Minute),
				}
				tokensRepo.On("Get", t.ID).Return(t, nil)
				tokensRepo.On("InvalidateByUserId", t.UserID).Return(nil)
				activeUsersRepo.On("Add", t.UserID).Return(nil)
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c.setupMock()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/chat/ws.rtm.start?token=%s", c.token), nil)
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewRouter())
			srv.Routes(users)
			srv.ServeHTTP(w, req)

			require.Equal(t, c.statusCode, w.Result().StatusCode)
			require.Equal(t, c.message, w.Body.String())

			tokensRepo.AssertExpectations(t)
			activeUsersRepo.AssertExpectations(t)
		})
	}
}
