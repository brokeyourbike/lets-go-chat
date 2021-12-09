package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	payload := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{
		UserName: "john",
		Password: "12345678",
	}

	user := models.User{ID: uuid.New(), UserName: payload.UserName}

	s.usersRepo.On("GetByUserName", payload.UserName).Return(models.User{}, errors.New("user not found"))
	s.usersRepo.On("Create", mock.AnythingOfType("User")).Return(user, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", s.preparePayload(payload))
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusOK, w.Result().StatusCode)
	require.Contains(s.T(), w.Body.String(), user.ID.String())
}

func (s *UsersSuite) Test_users_HandleUserCreate_InvalidJson() {
	table := map[string]struct {
		payload    string
		statusCode int
		message    string
	}{
		"json shold be valid": {
			payload:    `{not a valid json}`,
			statusCode: http.StatusBadRequest,
			message:    "Invalid json\n",
		},
		"username is required": {
			payload:    `{"password":"12345678"}`,
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
		},
		"password is required": {
			payload:    `{"username":"john"}`,
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
		},
		"username should be min 4 symbols": {
			payload:    `{"username":"jo", "password":"12345678"}`,
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
		},
		"password should be min 8 symbols": {
			payload:    `{"username":"john", "password":"1234567"}`,
			statusCode: http.StatusBadRequest,
			message:    "Invalid username or password\n",
		},
	}

	for name, tt := range table {
		s.T().Run(name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/user", strings.NewReader(tt.payload))
			w := httptest.NewRecorder()

			srv := server.NewServer(chi.NewRouter())
			srv.Routes(s.users)
			srv.ServeHTTP(w, req)

			require.Equal(s.T(), tt.statusCode, w.Result().StatusCode)
			require.Equal(s.T(), tt.message, w.Body.String())
		})
	}
}

func (s *UsersSuite) Test_users_HandleUserCreate_UserExist() {
	payload := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{
		UserName: "john",
		Password: "12345678",
	}

	s.usersRepo.On("GetByUserName", payload.UserName).Return(models.User{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/user", s.preparePayload(payload))
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusBadRequest, w.Result().StatusCode)
	require.Contains(s.T(), w.Body.String(), payload.UserName)
}

func (s *UsersSuite) Test_users_HandleUserCreate_CannotCreateUser() {
	payload := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{
		UserName: "john",
		Password: "12345678",
	}

	s.usersRepo.On("GetByUserName", payload.UserName).Return(models.User{}, errors.New("user not found"))
	s.usersRepo.On("Create", mock.AnythingOfType("User")).Return(models.User{}, errors.New("cannot create user"))

	req := httptest.NewRequest(http.MethodPost, "/v1/user", s.preparePayload(payload))
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusInternalServerError, w.Result().StatusCode)
}

func (s *UsersSuite) Test_users_HandleUserLogin() {
	payload := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{
		UserName: "john",
		Password: "12345678",
	}

	user := models.User{ID: uuid.New(), UserName: "john", PasswordHash: "$2a$04$hOGVri5G8ZLnrXtO/EP6keZkdzveoVGfh9krMXxeI/OP2QcSDJWOW"}
	token := models.Token{ID: uuid.New(), UserID: user.ID, ExpiresAt: time.Now().Add(time.Minute)}

	s.usersRepo.On("GetByUserName", payload.UserName).Return(user, nil)
	s.tokensRepo.On("Create", mock.AnythingOfType("Token")).Return(token, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", s.preparePayload(payload))
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusOK, w.Result().StatusCode)
	require.Equal(s.T(), token.ExpiresAt.UTC().String(), w.Result().Header.Get("X-Expires-After"))
}

func (s *UsersSuite) Test_users_HandleUserLogin_InvalidJson() {
	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", strings.NewReader("{not a valid json}"))
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusBadRequest, w.Result().StatusCode)
}

func (s *UsersSuite) Test_users_HandleUserLogin_UserNotFound() {
	payload := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{
		UserName: "john",
		Password: "12345678",
	}

	s.usersRepo.On("GetByUserName", payload.UserName).Return(models.User{}, db.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", s.preparePayload(payload))
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusBadRequest, w.Result().StatusCode)
}

func (s *UsersSuite) Test_users_HandleUserLogin_CannotQueryUser() {
	payload := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{
		UserName: "john",
		Password: "12345678",
	}

	s.usersRepo.On("GetByUserName", payload.UserName).Return(models.User{}, errors.New("cannot query user"))

	req := httptest.NewRequest(http.MethodPost, "/v1/user/login", s.preparePayload(payload))
	w := httptest.NewRecorder()

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(s.users)
	srv.ServeHTTP(w, req)

	require.Equal(s.T(), http.StatusInternalServerError, w.Result().StatusCode)
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
