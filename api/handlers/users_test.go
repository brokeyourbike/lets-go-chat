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
	"github.com/stretchr/testify/suite"
)

type userPayload struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type UsersSuite struct {
	suite.Suite
	usersRepo       *mocks.UsersRepo
	activeUsersRepo *mocks.ActiveUsersRepo
	tokensRepo      *mocks.TokensRepo
	users           *Users
}

func (s *UsersSuite) SetupTest() {
	s.usersRepo = new(mocks.UsersRepo)
	s.activeUsersRepo = new(mocks.ActiveUsersRepo)
	s.tokensRepo = new(mocks.TokensRepo)

	s.users = NewUsers(s.usersRepo, s.activeUsersRepo, s.tokensRepo)
}

func (s *UsersSuite) AfterTest(_, _ string) {
	require.True(s.T(), s.usersRepo.AssertExpectations(s.T()))
	require.True(s.T(), s.activeUsersRepo.AssertExpectations(s.T()))
	require.True(s.T(), s.tokensRepo.AssertExpectations(s.T()))
}

func TestUsers(t *testing.T) {
	suite.Run(t, new(UsersSuite))
}

func (s *UsersSuite) preparePayload(p interface{}) *bytes.Buffer {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	require.NoError(s.T(), err)
	return &buf
}

func (s *UsersSuite) Test_users_HandleUserCreate() {
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
				s.usersRepo.On("GetByUserName", "john").Return(models.User{}, nil)
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
				s.usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("user not found"))
				s.usersRepo.On("Create", mock.AnythingOfType("User")).Return(models.User{}, errors.New("cannot create user"))
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

				s.usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("user not found"))
				s.usersRepo.On("Create", mock.AnythingOfType("User")).Return(u, nil)
			},
		},
	}

	for name, c := range cases {
		s.T().Run(name, func(t *testing.T) {
			t.Parallel()
			c.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/v1/user", s.preparePayload(c.payload))
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewRouter())
			srv.Routes(s.users)
			srv.ServeHTTP(w, req)

			require.Equal(s.T(), c.statusCode, w.Result().StatusCode)
			require.Equal(s.T(), c.message, w.Body.String())
		})
	}
}

func (s *UsersSuite) Test_users_HandleUserLogin() {
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
				s.usersRepo.On("GetByUserName", "john").Return(models.User{}, db.ErrUserNotFound)
			},
		},
		"user can't login if db cannot perform query": {
			payload: userPayload{
				UserName: "john",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
			message:    "User with userName john not found\n",
			setupMock: func() {
				s.usersRepo.On("GetByUserName", "john").Return(models.User{}, errors.New("cannot query user"))
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

				s.usersRepo.On("GetByUserName", "john").Return(u, nil)
				s.tokensRepo.On("Create", mock.AnythingOfType("Token")).Return(t, nil)
			},
		},
	}

	for name, c := range cases {
		s.T().Run(name, func(t *testing.T) {
			t.Parallel()
			c.setupMock()

			req := httptest.NewRequest(http.MethodPost, "/v1/user/login", s.preparePayload(c.payload))
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewRouter())
			srv.Routes(s.users)
			srv.ServeHTTP(w, req)

			require.Equal(s.T(), c.statusCode, w.Result().StatusCode)
			require.Equal(s.T(), c.message, w.Body.String())
		})
	}
}

func (s *UsersSuite) Test_users_HandleUserActive() {
	s.activeUsersRepo.On("Count").Return(10)

	req := httptest.NewRequest(http.MethodGet, "/v1/user/active", nil)
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusOK, w.Result().StatusCode)
	require.Equal(s.T(), "{\"count\":10}\n", w.Body.String())
}

func (s *UsersSuite) Test_users_HandleChat_InvalidTokenFormat() {
	tests := []struct {
		name   string
		token  string
		status int
	}{
		{
			name:   "token can not be empty",
			token:  "",
			status: http.StatusBadRequest,
		},
		{
			name:   "token must be valid uuid",
			token:  "not-uuid",
			status: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/chat/ws.rtm.start?token=%s", tt.token), nil)
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewRouter())
			srv.Routes(s.users)
			srv.ServeHTTP(w, req)

			require.Equal(s.T(), tt.status, w.Result().StatusCode)
		})
	}
}

func (s *UsersSuite) Test_users_HandleChat_CannotValidateToken() {
	tokenID := uuid.New()

	s.tokensRepo.On("Get", tokenID).Return(models.Token{}, errors.New("cannot validate token"))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/chat/ws.rtm.start?token=%s", tokenID.String()), nil)
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusInternalServerError, w.Result().StatusCode)
	require.Equal(s.T(), "Token cannot be validated\n", w.Body.String())
}

func (s *UsersSuite) Test_users_HandleChat_TokenExpired() {
	token := models.Token{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: time.Now().Add(-time.Minute)}

	s.tokensRepo.On("Get", token.ID).Return(token, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/chat/ws.rtm.start?token=%s", token.ID.String()), nil)
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusBadRequest, w.Result().StatusCode)
	require.Equal(s.T(), "Token expired\n", w.Body.String())
}

func (s *UsersSuite) Test_users_HandleChat_RequestNotUpgradable() {
	token := models.Token{ID: uuid.New(), UserID: uuid.New(), ExpiresAt: time.Now().Add(time.Minute)}

	s.tokensRepo.On("Get", token.ID).Return(token, nil)
	s.tokensRepo.On("InvalidateByUserId", token.UserID).Return(nil)
	s.activeUsersRepo.On("Add", token.UserID).Return(nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/chat/ws.rtm.start?token=%s", token.ID.String()), nil)
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusBadRequest, w.Result().StatusCode)
	require.Equal(s.T(), "Bad Request\nCannot upgrade request to websocket protocol\n", w.Body.String())
}
