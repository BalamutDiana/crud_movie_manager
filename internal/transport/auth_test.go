package transport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BalamutDiana/crud_movie_manager/internal/domain"
	repo "github.com/BalamutDiana/crud_movie_manager/internal/repository"
	mock_service "github.com/BalamutDiana/crud_movie_manager/internal/transport/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, ctx context.Context, user domain.SignUpInput)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          domain.SignUpInput
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","email":"test@gmail.com","password":"qwerty"}`,
			inputUser: domain.SignUpInput{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUser, ctx context.Context, user domain.SignUpInput) {
				s.EXPECT().SignUp(gomock.Any(), user).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:               "Empty Fields",
			inputBody:          `{"email":"test@gmail.com","password":"qwerty"}`,
			mockBehavior:       func(s *mock_service.MockUser, ctx context.Context, user domain.SignUpInput) {},
			expectedStatusCode: 400,
		},
		{
			name:               "Email not valid",
			inputBody:          `{"name":"Test","email":"test","password":"qwerty"}`,
			mockBehavior:       func(s *mock_service.MockUser, ctx context.Context, user domain.SignUpInput) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name":"Test","email":"test@gmail.com","password":"qwerty"}`,
			inputUser: domain.SignUpInput{
				Name:     "Test",
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUser, ctx context.Context, user domain.SignUpInput) {
				s.EXPECT().SignUp(gomock.Any(), user).Return(errors.New("service failure"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockUser(c)
			testCase.mockBehavior(auth, context.Background(), testCase.inputUser)

			movieRepo := &repo.Movies{}
			h := NewHandler(movieRepo, auth)

			// Test Server
			r := mux.NewRouter()
			a := r.PathPrefix("/auth").Subrouter()
			{
				a.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
			}

			//Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, ctx context.Context, user domain.SignInInput)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            domain.SignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"test@gmail.com","password":"qwerty"}`,
			inputUser: domain.SignInInput{
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUser, ctx context.Context, user domain.SignInInput) {
				s.EXPECT().SignIn(gomock.Any(), user).Return("a", "r", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"a"}`,
		},
		{
			name:                 "Empty fields",
			inputBody:            `{"password":"qwerty"}`,
			mockBehavior:         func(s *mock_service.MockUser, ctx context.Context, user domain.SignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			name:                 "Email not valid",
			inputBody:            `{"email":"test","password":"qwerty"}`,
			mockBehavior:         func(s *mock_service.MockUser, ctx context.Context, user domain.SignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			name:      "Service Failure",
			inputBody: `{"email":"test@gmail.com","password":"qwerty"}`,
			inputUser: domain.SignInInput{
				Email:    "test@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockUser, ctx context.Context, user domain.SignInInput) {
				s.EXPECT().SignIn(gomock.Any(), user).Return("", "", errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: ``,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockUser(c)
			testCase.mockBehavior(auth, context.Background(), testCase.inputUser)

			movieRepo := &repo.Movies{}
			h := NewHandler(movieRepo, auth)

			// Test Server
			r := mux.NewRouter()
			a := r.PathPrefix("/auth").Subrouter()
			{
				a.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
			}

			//Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/auth/sign-in", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_refresh(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, ctx context.Context, cookie string)

	testTable := []struct {
		name                 string
		cookie               string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "OK",
			cookie: "r",
			mockBehavior: func(s *mock_service.MockUser, ctx context.Context, cookie string) {
				s.EXPECT().RefreshTokens(gomock.Any(), cookie).Return("a", "r", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"a"}`,
		},
		{
			name:   "Service Failure",
			cookie: "r",
			mockBehavior: func(s *mock_service.MockUser, ctx context.Context, cookie string) {
				s.EXPECT().RefreshTokens(gomock.Any(), cookie).Return("", "", errors.New("service failure"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: ``,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockUser(c)
			testCase.mockBehavior(auth, context.Background(), testCase.cookie)

			movieRepo := &repo.Movies{}
			h := NewHandler(movieRepo, auth)

			// Test Server
			r := mux.NewRouter()
			a := r.PathPrefix("/auth").Subrouter()
			{
				a.HandleFunc("/refresh", h.refresh).Methods(http.MethodGet)
			}

			//Test Request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/auth/refresh", nil)
			co := &http.Cookie{
				Name:   "refresh-token",
				Value:  testCase.cookie,
				MaxAge: 0,
			}
			req.AddCookie(co)

			r.ServeHTTP(w, req)

			//Assert
			fmt.Println(req.Cookies())
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
